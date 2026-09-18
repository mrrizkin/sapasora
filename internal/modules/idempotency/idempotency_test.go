package idempotency

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func testRequest(payload string) Request {
	return Request{
		Key: "request-123",
		Scope: Scope{
			WorkspaceID: "workspace-1",
			ActorID:     "actor-1",
			Endpoint:    "messages.create",
		},
		Payload: []byte(payload),
	}
}

func TestParseKeyIsBoundedAndDoesNotEchoInput(t *testing.T) {
	valid, err := ParseKey("order-123")
	require.NoError(t, err)
	require.Equal(t, "order-123", valid)

	for _, key := range []string{"", "has space", "has\nnewline", string(make([]byte, MaxKeyLength+1))} {
		_, err := ParseKey(key)
		require.ErrorIs(t, err, ErrInvalidKey)
		if key != "" {
			require.NotContains(t, err.Error(), key)
		}
	}
}

func TestFingerprintNormalizesJSONWithoutRetainingPayload(t *testing.T) {
	first, err := Fingerprint([]byte(`{"b":2,"a":1}`))
	require.NoError(t, err)
	second, err := Fingerprint([]byte(" { \"a\": 1, \"b\": 2 } "))
	require.NoError(t, err)
	require.Equal(t, first, second)

	different, err := Fingerprint([]byte(`{"a":1,"b":3}`))
	require.NoError(t, err)
	require.NotEqual(t, first, different)
}

func TestClaimScopesKeyByWorkspaceActorAndEndpoint(t *testing.T) {
	store := NewMemoryStore()
	request := testRequest(`{"message":"hello"}`)

	claim, err := store.Claim(request)
	require.NoError(t, err)
	require.True(t, claim.IsOwner())

	for name, mutate := range map[string]func(*Request){
		"workspace": func(r *Request) { r.Scope.WorkspaceID = "workspace-2" },
		"actor":     func(r *Request) { r.Scope.ActorID = "actor-2" },
		"endpoint":  func(r *Request) { r.Scope.Endpoint = "messages.retry" },
	} {
		t.Run(name, func(t *testing.T) {
			other := request
			mutate(&other)
			otherClaim, err := store.Claim(other)
			require.NoError(t, err)
			require.True(t, otherClaim.IsOwner())
		})
	}
}

func TestDuplicateReplaysOriginalResultAndConflictRedactsDetails(t *testing.T) {
	store := NewMemoryStore()
	request := testRequest(`{"amount":100}`)
	claim, err := store.Claim(request)
	require.NoError(t, err)

	result := Result{
		Status:  201,
		Headers: map[string]string{"X-Request-ID": "request-1"},
		Body:    []byte(`{"id":"message-1","status":"queued"}`),
	}
	require.NoError(t, store.Complete(claim, result))

	// Mutating the caller's result must not mutate the replayed result.
	result.Body[2] = 'X'
	result.Headers["X-Request-ID"] = "changed"

	replay, err := store.Claim(request)
	require.NoError(t, err)
	require.Equal(t, ClaimReplay, replay.Status)
	original, ok := replay.Replay()
	require.True(t, ok)
	require.Equal(t, 201, original.Status)
	require.Equal(t, []byte(`{"id":"message-1","status":"queued"}`), original.Body)
	require.Equal(t, "request-1", original.Headers["X-Request-ID"])

	conflicting := request
	conflicting.Payload = []byte(`{"amount":999,"secret":"do-not-echo"}`)
	_, err = store.Claim(conflicting)
	require.ErrorIs(t, err, ErrConflict)
	require.NotContains(t, err.Error(), "999")
	require.NotContains(t, err.Error(), "do-not-echo")
	require.NotContains(t, err.Error(), request.Key)
}

func TestTTLExpiresPendingAndCompletedRecords(t *testing.T) {
	current := time.Unix(100, 0)
	store := NewMemoryStore(
		WithClock(func() time.Time { return current }),
		WithDefaultTTL(time.Second),
		WithMaxTTL(time.Minute),
	)

	pending, err := store.Claim(testRequest("pending"))
	require.NoError(t, err)
	current = current.Add(2 * time.Second)
	reclaimed, err := store.Claim(testRequest("pending"))
	require.NoError(t, err)
	require.True(t, reclaimed.IsOwner())
	require.NotEqual(t, pending.token, reclaimed.token)

	completedRequest := testRequest("completed")
	completedRequest.Key = "completed-request"
	completed, err := store.Claim(completedRequest)
	require.NoError(t, err)
	require.NoError(t, store.Complete(completed, Result{Status: 200, Body: []byte("ok")}))
	replay, err := store.Claim(completedRequest)
	require.NoError(t, err)
	require.Equal(t, ClaimReplay, replay.Status)

	current = current.Add(2 * time.Second)
	expired, err := store.Claim(completedRequest)
	require.NoError(t, err)
	require.True(t, expired.IsOwner())
}

func TestConcurrentDuplicateClaimsHaveSingleOwner(t *testing.T) {
	store := NewMemoryStore()
	const callers = 32

	start := make(chan struct{})
	claims := make(chan Claim, callers)
	errs := make(chan error, callers)
	var started atomic.Int32
	var group sync.WaitGroup
	group.Add(callers)
	for i := 0; i < callers; i++ {
		go func() {
			defer group.Done()
			started.Add(1)
			<-start
			claim, err := store.Claim(testRequest(`{"same":"payload"}`))
			claims <- claim
			errs <- err
		}()
	}
	for started.Load() != callers {
		time.Sleep(time.Millisecond)
	}
	close(start)
	group.Wait()
	close(claims)
	close(errs)

	var owner Claim
	owners := 0
	inProgress := 0
	for claim := range claims {
		switch claim.Status {
		case ClaimAcquired:
			owners++
			owner = claim
		case ClaimInProgress:
			inProgress++
		}
	}
	for err := range errs {
		if err != nil {
			require.ErrorIs(t, err, ErrInProgress)
		}
	}
	require.Equal(t, 1, owners)
	require.Equal(t, callers-1, inProgress)

	require.NoError(t, store.Complete(owner, Result{Status: 202, Body: []byte("accepted")}))
	replay, err := store.Claim(testRequest(`{"same":"payload"}`))
	require.NoError(t, err)
	require.Equal(t, ClaimReplay, replay.Status)
	require.Equal(t, []byte("accepted"), replay.StoredResult.Body)
}

func TestInvalidInputsAreBoundedAndRedacted(t *testing.T) {
	store := NewMemoryStore(WithMaxPayloadBytes(4), WithMaxResultBytes(4))

	request := testRequest("secret-payload")
	_, err := store.Claim(request)
	require.ErrorIs(t, err, ErrPayloadTooLarge)
	require.NotContains(t, err.Error(), "secret-payload")

	request = testRequest("ok")
	claim, err := store.Claim(request)
	require.NoError(t, err)
	require.ErrorIs(t, store.Complete(claim, Result{Body: []byte("secret-result")}), ErrResultTooLarge)

	_, err = store.Claim(Request{Key: request.Key, Scope: request.Scope, Payload: []byte("ok"), TTL: -time.Second})
	require.ErrorIs(t, err, ErrInvalidTTL)
	require.True(t, errors.Is(err, ErrInvalidTTL))
}
