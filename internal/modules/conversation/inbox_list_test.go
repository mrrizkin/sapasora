package conversation

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func inboxConversation(scope, publicID string, activity time.Time) Conversation {
	return Conversation{
		TenantID:        scope,
		PublicID:        publicID,
		ChannelPublicID: "channel-1",
		ContactPublicID: "contact-" + publicID,
		Status:          ConversationStatusOpen,
		Priority:        ConversationPriorityNormal,
		Tags:            []string{"support"},
		LastActivityAt:  activity,
		CreatedAt:       activity,
		UnreadCount:     1,
	}
}

func createInboxConversations(t *testing.T, repository *InMemoryRepository, values ...Conversation) {
	t.Helper()
	for index := range values {
		if err := repository.CreateConversation(context.Background(), &values[index]); err != nil {
			t.Fatalf("CreateConversation(%s): %v", values[index].PublicID, err)
		}
	}
}

func TestInboxListUsesStableBoundedCursorOrdering(t *testing.T) {
	repository := NewInMemoryRepository()
	at := time.Date(2026, time.January, 2, 15, 4, 5, 0, time.UTC)
	values := []Conversation{
		inboxConversation("tenant-a", "conversation-1", at),
		inboxConversation("tenant-a", "conversation-2", at.Add(-time.Hour)),
		inboxConversation("tenant-a", "conversation-3", at.Add(-2*time.Hour)),
		inboxConversation("tenant-a", "conversation-4", at.Add(-3*time.Hour)),
		inboxConversation("tenant-a", "conversation-5", at.Add(-4*time.Hour)),
	}
	createInboxConversations(t, repository, values...)
	service := NewInboxListService(repository)
	request := InboxListRequest{TenantID: "tenant-a", Limit: 2, AsOf: at}

	first := service.List(context.Background(), request)
	if first.State != InboxListStateReady || !first.HasMore || len(first.Items) != 2 {
		t.Fatalf("first result = %+v", first)
	}
	if got := first.Items[0].PublicID + "," + first.Items[1].PublicID; got != "conversation-1,conversation-2" {
		t.Fatalf("first order = %s", got)
	}
	if first.NextCursor == "" || len(first.NextCursor) > maxInboxCursorSize {
		t.Fatalf("next cursor = %q", first.NextCursor)
	}

	request.Cursor = first.NextCursor
	second := service.List(context.Background(), request)
	if second.State != InboxListStateReady || !second.HasMore || len(second.Items) != 2 {
		t.Fatalf("second result = %+v", second)
	}
	if got := second.Items[0].PublicID + "," + second.Items[1].PublicID; got != "conversation-3,conversation-4" {
		t.Fatalf("second order = %s", got)
	}

	request.Cursor = second.NextCursor
	third := service.List(context.Background(), request)
	if third.State != InboxListStateReady || third.HasMore || len(third.Items) != 1 || third.Items[0].PublicID != "conversation-5" {
		t.Fatalf("third result = %+v", third)
	}
	request.Cursor = encodeInboxCursor(third.Items[0], "tenant-a", request.Filter, time.Time{})
	fourth := service.List(context.Background(), request)
	if fourth.State != InboxListStateEmpty || fourth.HasMore || len(fourth.Items) != 0 {
		t.Fatalf("empty result = %+v", fourth)
	}
}

func TestInboxListTieBreakAndTenantIsolation(t *testing.T) {
	repository := NewInMemoryRepository()
	at := time.Date(2026, time.January, 2, 15, 4, 5, 0, time.UTC)
	createInboxConversations(t, repository,
		inboxConversation("tenant-a", "a-low", at),
		inboxConversation("tenant-a", "z-high", at),
		inboxConversation("tenant-b", "tenant-b-only", at.Add(time.Hour)),
	)
	service := NewInboxListService(repository)
	result := service.List(context.Background(), InboxListRequest{TenantID: "tenant-a", Limit: MaxInboxPageSize, AsOf: at})
	if result.State != InboxListStateReady || len(result.Items) != 2 {
		t.Fatalf("tenant-scoped result = %+v", result)
	}
	if result.Items[0].PublicID != "z-high" || result.Items[1].PublicID != "a-low" {
		t.Fatalf("tie order = %q, %q", result.Items[0].PublicID, result.Items[1].PublicID)
	}
	for _, item := range result.Items {
		if item.ScopeID() != "tenant-a" {
			t.Fatalf("cross-tenant item = %+v", item)
		}
	}
	workspaceResult := service.List(context.Background(), InboxListRequest{WorkspaceID: "tenant-a", Limit: MaxInboxPageSize, AsOf: at})
	if workspaceResult.State != InboxListStateReady || len(workspaceResult.Items) != 2 {
		t.Fatalf("workspace-scoped result = %+v", workspaceResult)
	}
	mismatched := service.List(context.Background(), InboxListRequest{TenantID: "tenant-a", WorkspaceID: "tenant-b"})
	if mismatched.State != InboxListStateError || mismatched.Error == nil || mismatched.Error.Code != "invalid_request" {
		t.Fatalf("mismatched scope result = %+v", mismatched)
	}
}

func TestInboxListFiltersSearchAndSLA(t *testing.T) {
	repository := NewInMemoryRepository()
	at := time.Date(2026, time.January, 2, 15, 4, 5, 0, time.UTC)
	breached := inboxConversation("tenant-a", "breached", at)
	breached.ChannelPublicID = "channel-whatsapp"
	breached.Status = ConversationStatusPending
	breached.Priority = ConversationPriorityUrgent
	breached.AssigneePublicID = "agent-7"
	breached.TeamPublicID = "team-red"
	breached.Tags = []string{"vip"}
	due := at.Add(-time.Minute)
	breached.SLAResolutionDueAt = &due
	unbreached := inboxConversation("tenant-a", "unbreached", at.Add(-time.Hour))
	unbreached.ChannelPublicID = "channel-email"
	unbreached.UnreadCount = 0
	createInboxConversations(t, repository, breached, unbreached)

	service := NewInboxListService(repository)
	status := ConversationStatusPending
	priority := ConversationPriorityUrgent
	result := service.List(context.Background(), InboxListRequest{
		TenantID: "tenant-a", AsOf: at, Limit: 10,
		Filter: InboxListFilter{
			Search: "WHATSAPP", ChannelPublicID: "channel-whatsapp", Status: &status,
			AssigneePublicID: "agent-7", TeamPublicID: "team-red", Priority: &priority,
			Tag: "VIP", UnreadOnly: true, SLABreachedOnly: true,
		},
	})
	if result.State != InboxListStateReady || len(result.Items) != 1 || result.Items[0].PublicID != "breached" {
		t.Fatalf("filtered result = %+v", result)
	}
}

func TestParseInboxListQueryIsStrictAndBounded(t *testing.T) {
	query, err := ParseInboxListQuery(url.Values{
		"search": {"channel"}, "channel_id": {"channel-1"}, "status": {"OPEN"},
		"priority": {"high"}, "tag": {"VIP"}, "unread": {"1"}, "sla": {"breached"}, "limit": {"10"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if query.Limit != 10 || query.Status == nil || *query.Status != ConversationStatusOpen || !query.UnreadOnly || !query.SLABreachedOnly || query.Tag != "vip" {
		t.Fatalf("parsed query = %+v", query)
	}
	for _, raw := range []string{
		"unknown=value", "limit=0", "limit=101", "limit=", "unread=yes", "status=bogus",
		"channel=a&channel_id=b", "search=" + strings.Repeat("x", maxInboxQueryValue+1),
	} {
		if _, err := ParseInboxQuery(raw); !errors.Is(err, ErrInvalidInboxQuery) {
			t.Fatalf("ParseInboxQuery(%q) error = %v", raw, err)
		}
	}
	if got := NewInboxListLoadingResult(); got.State != InboxListStateLoading || got.Items == nil {
		t.Fatalf("loading result = %+v", got)
	}
}

func TestInboxListErrorAndDiagnosticsAreSafe(t *testing.T) {
	repositoryError := errors.New("sql failed for Private Person +15551234567 secret")
	service := NewInboxListService(errorInboxRepository{Repository: NewInMemoryRepository(), err: repositoryError})
	result := service.List(context.Background(), InboxListRequest{TenantID: "tenant-a"})
	if result.State != InboxListStateError || result.Error == nil || result.Error.Code != "repository_error" {
		t.Fatalf("error result = %+v", result)
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	for _, secret := range []string{"Private Person", "+15551234567", "secret", "sql failed"} {
		if strings.Contains(string(encoded), secret) {
			t.Fatalf("error diagnostics leaked %q: %s", secret, encoded)
		}
	}
}

type errorInboxRepository struct {
	Repository
	err error
}

func (r errorInboxRepository) ListConversations(context.Context, string, ConversationFilter) ([]*Conversation, error) {
	return nil, r.err
}

func TestInboxListConcurrentReadsAreRaceSafe(t *testing.T) {
	repository := NewInMemoryRepository()
	at := time.Date(2026, time.January, 2, 15, 4, 5, 0, time.UTC)
	values := make([]Conversation, 20)
	for index := range values {
		values[index] = inboxConversation("tenant-a", "conversation-"+strconv.Itoa(index), at.Add(-time.Duration(index)*time.Minute))
	}
	createInboxConversations(t, repository, values...)
	service := NewInboxListService(repository)

	const workers = 32
	var wait sync.WaitGroup
	wait.Add(workers)
	errorsCh := make(chan string, workers)
	for index := 0; index < workers; index++ {
		go func() {
			defer wait.Done()
			result := service.List(context.Background(), InboxListRequest{TenantID: "tenant-a", Limit: 5, AsOf: at})
			if result.State != InboxListStateReady || len(result.Items) != 5 {
				errorsCh <- "unexpected concurrent list result"
			}
		}()
	}
	wait.Wait()
	close(errorsCh)
	for message := range errorsCh {
		t.Fatal(message)
	}
}
