package message

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
)

func testMessage(scope, publicID string) Message {
	return Message{
		WorkspaceID:          scope,
		PublicID:             publicID,
		Direction:            MessageDirectionOutbound,
		Provider:             "provider-neutral",
		ProviderMessageID:    "provider-message-1",
		IdempotencyKey:       "idem-" + publicID,
		RequestID:            "request-" + publicID,
		ChannelPublicID:      "channel-1",
		ContactPublicID:      "contact-1",
		ConversationPublicID: "conversation-1",
		NormalizedContent:    "  private message body  ",
		ContentMetadata:      NormalizedContentMetadata{Kind: ContentKindText},
		ProviderMetadata:     map[string]string{"provider_status": "accepted", "access_token": "private"},
	}
}

func TestInMemoryRepositoryScopesMessagesAndCanonicalizesContent(t *testing.T) {
	repository := NewInMemoryRepository()
	ctx := context.Background()
	first := testMessage("workspace-a", "message-1")
	if err := repository.CreateMessage(ctx, &first); err != nil {
		t.Fatalf("CreateMessage() error = %v", err)
	}
	if first.ID == 0 || first.TenantID != "workspace-a" || first.WorkspaceID != "workspace-a" {
		t.Fatalf("canonical scope/persistence fields = %+v", first)
	}
	if first.NormalizedContent != "private message body" || first.ContentMetadata.CharacterCount != len([]rune("private message body")) {
		t.Fatalf("normalized content metadata = %+v", first)
	}
	if first.RedactedProviderMetadata["access_token"] != "[REDACTED]" {
		t.Fatalf("provider metadata was not redacted: %+v", first.RedactedProviderMetadata)
	}
	if _, err := repository.GetMessageByPublicID(ctx, "workspace-b", first.PublicID); !errors.Is(err, ErrMessageNotFound) {
		t.Fatalf("cross-workspace lookup error = %v", err)
	}
	second := testMessage("workspace-b", "message-1")
	if err := repository.CreateMessage(ctx, &second); err != nil {
		t.Fatalf("same public ID in another workspace error = %v", err)
	}
	if err := repository.CreateMessage(ctx, &Message{
		TenantID: "workspace-a", PublicID: "message-2", Direction: MessageDirectionOutbound,
		IdempotencyKey: first.IdempotencyKey,
	}); !errors.Is(err, ErrMessageConflict) {
		t.Fatalf("duplicate idempotency key error = %v", err)
	}
	first.IdempotencyKey = ""
	if err := repository.UpdateMessage(ctx, &first); err != nil {
		t.Fatalf("clearing idempotency key error = %v", err)
	}
	releasedKey := Message{TenantID: "workspace-a", PublicID: "message-2", Direction: MessageDirectionOutbound, IdempotencyKey: "idem-message-1"}
	if err := repository.CreateMessage(ctx, &releasedKey); err != nil {
		t.Fatalf("reusing cleared idempotency key error = %v", err)
	}
}

func TestMessageDiagnosticsExcludeContentAndOpaqueSecrets(t *testing.T) {
	value := testMessage("workspace-a", "message-1")
	if err := NewInMemoryRepository().CreateMessage(context.Background(), &value); err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	output := value.String() + string(encoded)
	for _, secret := range []string{"private message body", "provider-message-1", "idem-message-1", "request-message-1", "private"} {
		if strings.Contains(output, secret) {
			t.Fatalf("diagnostics leaked %q: %s", secret, output)
		}
	}
	if !strings.Contains(string(encoded), `"content_metadata"`) {
		t.Fatalf("safe diagnostics omitted content metadata: %s", encoded)
	}
}

func TestInMemoryRepositoryStoresAssociationsAndAppendOnlyRecords(t *testing.T) {
	repository := NewInMemoryRepository()
	ctx := context.Background()
	value := testMessage("tenant-a", "message-1")
	if err := repository.CreateMessage(ctx, &value); err != nil {
		t.Fatal(err)
	}
	attachment := MessageAttachment{
		TenantID: "tenant-a", MessagePublicID: value.PublicID, Kind: AttachmentKindDocument,
		FileName: "private.pdf", MIMEType: "application/pdf", SizeBytes: 42,
		StorageRef: "private-storage-ref", Metadata: map[string]string{"checksum": "safe"},
	}
	if err := repository.CreateAttachment(ctx, &attachment); err != nil {
		t.Fatalf("CreateAttachment() error = %v", err)
	}
	if attachment.MessageID != value.ID || attachment.ID == 0 {
		t.Fatalf("attachment associations = %+v", attachment)
	}
	attachmentList, err := repository.ListAttachments(ctx, "tenant-a", value.PublicID)
	if err != nil || len(attachmentList) != 1 {
		t.Fatalf("ListAttachments() = %+v, %v", attachmentList, err)
	}
	attachmentList[0].Metadata["checksum"] = "caller-mutated"
	attachmentList, err = repository.ListAttachments(ctx, "tenant-a", value.PublicID)
	if err != nil || attachmentList[0].Metadata["checksum"] != "safe" {
		t.Fatalf("attachment defensive copy = %+v, %v", attachmentList, err)
	}

	delivery := DeliveryRecord{
		TenantID: "tenant-a", MessagePublicID: value.PublicID, ChannelPublicID: "channel-1",
		Provider: "provider-neutral", ProviderMessageID: "provider-message-1", Status: MessageStatusSent,
		RedactedProviderMetadata: map[string]string{"provider_token": "secret"},
	}
	if err := repository.CreateDeliveryRecord(ctx, &delivery); err != nil {
		t.Fatalf("CreateDeliveryRecord() error = %v", err)
	}
	if delivery.MessageID != value.ID || delivery.RedactedProviderMetadata["provider_token"] != "[REDACTED]" {
		t.Fatalf("delivery record = %+v", delivery)
	}

	for _, eventType := range []MessageEventType{MessageEventCreated, MessageEventSent, MessageEventDelivered} {
		event := MessageEvent{TenantID: "tenant-a", MessagePublicID: value.PublicID, Type: eventType, Status: MessageStatus(eventType)}
		if eventType == MessageEventCreated {
			event.Status = MessageStatusAccepted
		}
		if err := repository.AppendMessageEvent(ctx, &event); err != nil {
			t.Fatalf("AppendMessageEvent(%s) error = %v", eventType, err)
		}
	}
	events, err := repository.ListMessageEvents(ctx, "tenant-a", value.PublicID)
	if err != nil || len(events) != 3 {
		t.Fatalf("ListMessageEvents() = %+v, %v", events, err)
	}
	for index, event := range events {
		if event.Sequence != uint64(index+1) {
			t.Fatalf("event sequence = %d at %d", event.Sequence, index)
		}
	}
	if _, err := repository.ListMessageEvents(ctx, "tenant-b", value.PublicID); !errors.Is(err, ErrMessageNotFound) {
		t.Fatalf("cross-tenant event lookup error = %v", err)
	}
}

func TestMessageValidationRejectsUnsafeOrIncompleteOwnership(t *testing.T) {
	if err := (Message{TenantID: "tenant-a", WorkspaceID: "tenant-b", PublicID: "message", Direction: MessageDirectionInbound, Status: MessageStatusAccepted}).Valid(); err == nil {
		t.Fatal("mismatched tenant/workspace scope accepted")
	}
	if err := (Message{TenantID: "tenant-a", PublicID: "message", Direction: MessageDirection("sideways"), Status: MessageStatusAccepted}).Valid(); err == nil {
		t.Fatal("invalid direction accepted")
	}
	if err := (MessageAttachment{TenantID: "tenant-a", PublicID: "attachment", MessagePublicID: "message", Kind: AttachmentKindImage, SizeBytes: -1}).Valid(); err == nil {
		t.Fatal("negative attachment size accepted")
	}
}

func TestInMemoryRepositoryConcurrentCreatesAndEvents(t *testing.T) {
	repository := NewInMemoryRepository()
	ctx := context.Background()
	const workers = 32
	var wait sync.WaitGroup
	errorsCh := make(chan error, workers)
	for i := 0; i < workers; i++ {
		wait.Add(1)
		go func(index int) {
			defer wait.Done()
			value := Message{TenantID: "tenant-a", PublicID: fmt.Sprintf("message-%d", index), Direction: MessageDirectionInbound, Content: "inbound"}
			if err := repository.CreateMessage(ctx, &value); err != nil {
				errorsCh <- err
				return
			}
			event := MessageEvent{TenantID: "tenant-a", MessagePublicID: value.PublicID, Type: MessageEventCreated}
			if err := repository.AppendMessageEvent(ctx, &event); err != nil {
				errorsCh <- err
			}
		}(i)
	}
	wait.Wait()
	close(errorsCh)
	for err := range errorsCh {
		t.Fatal(err)
	}
	messages, err := repository.ListMessages(ctx, "tenant-a", MessageFilter{})
	if err != nil || len(messages) != workers {
		t.Fatalf("concurrent messages = %d, %v", len(messages), err)
	}
	for _, value := range messages {
		events, err := repository.ListMessageEvents(ctx, "tenant-a", value.PublicID)
		if err != nil || len(events) != 1 {
			t.Fatalf("events for %s = %+v, %v", value.PublicID, events, err)
		}
	}
}
