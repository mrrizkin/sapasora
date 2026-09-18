package conversation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func testConversation(scope, publicID string) Conversation {
	return Conversation{
		TenantID:        scope,
		PublicID:        publicID,
		ChannelPublicID: "channel-1",
		ContactPublicID: "contact-1",
		Status:          ConversationStatusOpen,
		Priority:        ConversationPriorityNormal,
		Tags:            []string{"VIP", "support"},
		UnreadCount:     2,
	}
}

func TestInMemoryRepositoryScopesConversationStateToTenant(t *testing.T) {
	repository := NewInMemoryRepository()
	ctx := context.Background()
	conversationA := testConversation("tenant-a", "conversation-1")
	conversationB := testConversation("tenant-b", "conversation-1")
	conversationB.Status = ConversationStatusPending

	if err := repository.CreateConversation(ctx, &conversationA); err != nil {
		t.Fatalf("CreateConversation(a) error = %v", err)
	}
	if err := repository.CreateConversation(ctx, &conversationB); err != nil {
		t.Fatalf("CreateConversation(b) error = %v", err)
	}
	if conversationA.ID == conversationB.ID || conversationA.ID == 0 {
		t.Fatalf("persistence IDs = %d and %d", conversationA.ID, conversationB.ID)
	}

	got, err := repository.GetConversationByPublicID(ctx, "tenant-a", "conversation-1")
	if err != nil || got.Status != ConversationStatusOpen {
		t.Fatalf("tenant-a lookup = %+v, %v", got, err)
	}
	if _, err := repository.GetConversationByPublicID(ctx, "tenant-b", "missing"); !errors.Is(err, ErrConversationNotFound) {
		t.Fatalf("missing lookup error = %v", err)
	}
	if _, err := repository.GetConversationByID(ctx, "tenant-b", conversationA.ID); !errors.Is(err, ErrConversationNotFound) {
		t.Fatalf("cross-tenant internal lookup error = %v", err)
	}

	pending := ConversationStatusPending
	values, err := repository.ListConversations(ctx, "tenant-b", ConversationFilter{Status: &pending})
	if err != nil || len(values) != 1 || values[0].PublicID != conversationB.PublicID {
		t.Fatalf("tenant-b state-filtered list = %+v, %v", values, err)
	}
	values, err = repository.ListConversations(ctx, "tenant-a", ConversationFilter{Status: &pending})
	if err != nil || len(values) != 0 {
		t.Fatalf("tenant-a leaked tenant-b state = %+v, %v", values, err)
	}

	conversationA.Status = ConversationStatusResolved
	if err := repository.UpdateConversation(ctx, &conversationA); err != nil {
		t.Fatalf("UpdateConversation() error = %v", err)
	}
	if got, err := repository.GetConversationByPublicID(ctx, "tenant-a", conversationA.PublicID); err != nil || got.Status != ConversationStatusResolved {
		t.Fatalf("updated state = %+v, %v", got, err)
	}
}

func TestInMemoryRepositoryAssociationsAreTenantAndConversationScoped(t *testing.T) {
	repository := NewInMemoryRepository()
	ctx := context.Background()
	conversationA := testConversation("tenant-a", "conversation-a")
	conversationB := testConversation("tenant-b", "conversation-b")
	if err := repository.CreateConversation(ctx, &conversationA); err != nil {
		t.Fatal(err)
	}
	if err := repository.CreateConversation(ctx, &conversationB); err != nil {
		t.Fatal(err)
	}

	participant := Participant{
		TenantID:             "tenant-a",
		ConversationPublicID: conversationA.PublicID,
		ContactPublicID:      "contact-a",
		Role:                 ParticipantRoleCustomer,
		DisplayName:          "Private Customer",
		Address:              "+15551234567",
		Metadata:             map[string]string{"private": "secret"},
	}
	if err := repository.CreateParticipant(ctx, &participant); err != nil {
		t.Fatalf("CreateParticipant() error = %v", err)
	}
	assignment := ConversationAssignment{
		TenantID:             "tenant-a",
		ConversationPublicID: conversationA.PublicID,
		AssigneePublicID:     "agent-a",
	}
	if err := repository.CreateAssignment(ctx, &assignment); err != nil {
		t.Fatalf("CreateAssignment() error = %v", err)
	}
	if _, err := repository.ListParticipants(ctx, "tenant-b", conversationA.PublicID); !errors.Is(err, ErrConversationNotFound) {
		t.Fatalf("cross-tenant participant list error = %v", err)
	}
	if _, err := repository.GetAssignmentByPublicID(ctx, "tenant-b", assignment.PublicID); !errors.Is(err, ErrConversationAssignmentNotFound) {
		t.Fatalf("cross-tenant assignment lookup error = %v", err)
	}
	participants, err := repository.ListParticipants(ctx, "tenant-a", conversationA.PublicID)
	if err != nil || len(participants) != 1 || participants[0].ID != participant.ID {
		t.Fatalf("participants = %+v, %v", participants, err)
	}
	assignments, err := repository.ListAssignments(ctx, "tenant-a", conversationA.PublicID)
	if err != nil || len(assignments) != 1 || assignments[0].ID != assignment.ID {
		t.Fatalf("assignments = %+v, %v", assignments, err)
	}

	if err := repository.DeleteConversation(ctx, "tenant-a", conversationA.PublicID); err != nil {
		t.Fatalf("DeleteConversation() error = %v", err)
	}
	if _, err := repository.ListAssignments(ctx, "tenant-a", conversationA.PublicID); !errors.Is(err, ErrConversationNotFound) {
		t.Fatalf("deleted conversation association list error = %v", err)
	}
	if _, err := repository.GetConversationByPublicID(ctx, "tenant-b", conversationB.PublicID); err != nil {
		t.Fatalf("tenant-b conversation after tenant-a delete = %v", err)
	}
}

func TestConversationValidationAndDiagnosticsArePIISafe(t *testing.T) {
	conversation := testConversation("workspace-a", "conversation-1")
	conversation.WorkspaceID = "workspace-a"
	conversation.Tags = []string{"VIP"}
	if err := conversation.Valid(); err != nil {
		t.Fatalf("Conversation.Valid() error = %v", err)
	}
	if err := (Conversation{TenantID: "tenant-a", WorkspaceID: "workspace-b", PublicID: "id", Status: ConversationStatusOpen, Priority: ConversationPriorityNormal}).Valid(); err == nil {
		t.Fatal("Conversation.Valid() accepted mismatched tenant/workspace")
	}
	if err := (Conversation{TenantID: "tenant-a", PublicID: "id", Status: ConversationStatusOpen, Priority: ConversationPriorityNormal, Tags: []string{"vip", "VIP"}}).Valid(); err == nil {
		t.Fatal("Conversation.Valid() accepted duplicate tags")
	}

	participant := Participant{
		PublicID:             "participant-1",
		ConversationPublicID: "conversation-1",
		TenantID:             "tenant-a",
		Role:                 ParticipantRoleCustomer,
		DisplayName:          "Private Person",
		Address:              "+15551234567",
		JoinedAt:             time.Now(),
		Metadata:             map[string]string{"note": "secret"},
	}
	stringDiagnostics := participant.String()
	encoded, err := json.Marshal(participant)
	if err != nil {
		t.Fatal(err)
	}
	for _, output := range []string{stringDiagnostics, string(encoded)} {
		for _, secret := range []string{"Private Person", "+15551234567", "secret"} {
			if strings.Contains(output, secret) {
				t.Fatalf("diagnostics leaked %q: %s", secret, output)
			}
		}
	}
}

func TestInMemoryRepositoryDefensiveCopiesAndConcurrentCreates(t *testing.T) {
	repository := NewInMemoryRepository()
	ctx := context.Background()
	conversation := testConversation("tenant-a", "conversation-1")
	if err := repository.CreateConversation(ctx, &conversation); err != nil {
		t.Fatal(err)
	}
	conversation.Tags[0] = "caller-mutated"
	stored, err := repository.GetConversationByPublicID(ctx, "tenant-a", "conversation-1")
	if err != nil || stored.Tags[0] != "vip" {
		t.Fatalf("conversation defensive copy = %+v, %v", stored, err)
	}

	const workers = 32
	var wait sync.WaitGroup
	errorsCh := make(chan error, workers)
	for i := 0; i < workers; i++ {
		wait.Add(1)
		go func(index int) {
			defer wait.Done()
			value := testConversation("tenant-a", fmt.Sprintf("conversation-%d", index+2))
			if err := repository.CreateConversation(ctx, &value); err != nil {
				errorsCh <- err
			}
		}(i)
	}
	wait.Wait()
	close(errorsCh)
	for err := range errorsCh {
		t.Fatalf("concurrent create error = %v", err)
	}
	values, err := repository.ListConversations(ctx, "tenant-a", ConversationFilter{})
	if err != nil || len(values) != workers+1 {
		t.Fatalf("concurrent conversation count = %d, %v", len(values), err)
	}
}
