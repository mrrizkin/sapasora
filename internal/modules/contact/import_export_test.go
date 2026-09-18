package contact

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestImportCSVValidRowsNormalizeAndGroupContacts(t *testing.T) {
	service := NewContactService(NewInMemoryContactRepository())
	input := strings.Join([]string{
		strings.Join(CSVHeader, ","),
		"contact-alice,Alice,active,import,email,, Alice@Example.COM ",
		"contact-alice,Alice,active,import,phone,,+1 (555) 123-4567",
	}, "\n")

	result, err := ImportCSV(context.Background(), service, "tenant-a", strings.NewReader(input), CSVImportOptions{})
	if err != nil {
		t.Fatalf("ImportCSV() error = %v", err)
	}
	if result.RowsImported != 2 || result.ContactsCreated != 1 || result.AddressesCreated != 2 || len(result.Issues) != 0 {
		t.Fatalf("ImportCSV() result = %+v", result)
	}
	contact, err := service.GetContact(context.Background(), "tenant-a", "contact-alice")
	if err != nil || contact.DisplayName != "Alice" {
		t.Fatalf("imported contact = %+v, %v", contact, err)
	}
	addresses, err := service.ListContactAddressesByContactID(context.Background(), "tenant-a", contact.ID)
	if err != nil || len(addresses) != 2 {
		t.Fatalf("imported addresses = %+v, %v", addresses, err)
	}
	if addresses[0].NormalizedValue != "alice@example.com" || addresses[1].NormalizedValue != "+15551234567" {
		t.Fatalf("normalized addresses = %+v", addresses)
	}
}

func TestImportCSVReportsMalformedRowsWithoutPII(t *testing.T) {
	service := NewContactService(NewInMemoryContactRepository())
	secret := "private person@example.com"
	input := strings.Join([]string{
		strings.Join(CSVHeader, ","),
		"good,Good,active,import,email,,good@example.com",
		"bad,Bad,active,import,email,," + secret,
		"bad-status,Bad,not-a-status,import,email,,status@example.com",
		"wrong,row",
	}, "\n")

	result, err := ImportCSV(context.Background(), service, "tenant-a", strings.NewReader(input), CSVImportOptions{})
	if err != nil {
		t.Fatalf("ImportCSV() error = %v", err)
	}
	if result.RowsImported != 1 || result.RowsSkipped != 3 || len(result.Issues) != 3 {
		t.Fatalf("ImportCSV() result = %+v", result)
	}
	for _, issue := range result.Issues {
		if strings.Contains(issue.Error(), secret) {
			t.Fatalf("issue leaked PII: %v", issue)
		}
	}
	if _, err := service.GetContact(context.Background(), "tenant-a", "bad"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("malformed contact lookup error = %v", err)
	}
}

func TestImportCSVReportsDuplicateIdentities(t *testing.T) {
	service := NewContactService(NewInMemoryContactRepository())
	input := strings.Join([]string{
		strings.Join(CSVHeader, ","),
		"first,First,active,import,email,,Alice@Example.com",
		"second,Second,active,import,email,, alice@example.com ",
	}, "\n")

	result, err := ImportCSV(context.Background(), service, "tenant-a", strings.NewReader(input), CSVImportOptions{})
	if err != nil {
		t.Fatalf("ImportCSV() error = %v", err)
	}
	if result.RowsImported != 1 || result.DuplicateRows != 1 || len(result.Issues) != 1 || result.Issues[0].Code != "duplicate_identity" {
		t.Fatalf("duplicate result = %+v", result)
	}
}

func TestImportCSVEnforcesMaxRows(t *testing.T) {
	service := NewContactService(NewInMemoryContactRepository())
	input := strings.Join([]string{
		strings.Join(CSVHeader, ","),
		"first,First,active,import,email,,first@example.com",
		"second,Second,active,import,email,,second@example.com",
	}, "\n")

	result, err := ImportCSV(context.Background(), service, "tenant-a", strings.NewReader(input), CSVImportOptions{MaxRows: 1})
	if !errors.Is(err, ErrCSVLimit) || result.RowsRead != 1 || result.RowsImported != 1 {
		t.Fatalf("max rows result = %+v, error = %v", result, err)
	}
}

func TestImportCSVDryRunDoesNotMutate(t *testing.T) {
	repository := NewInMemoryContactRepository()
	service := NewContactService(repository)
	input := strings.Join([]string{
		strings.Join(CSVHeader, ","),
		"contact,Contact,active,import,email,,dry-run@example.com",
	}, "\n")

	result, err := ImportCSV(context.Background(), service, "tenant-a", strings.NewReader(input), CSVImportOptions{DryRun: true})
	if err != nil {
		t.Fatalf("ImportCSV() error = %v", err)
	}
	if !result.DryRun || result.RowsImported != 0 || result.RowsWouldImport != 1 || result.ContactsCreated != 0 || result.AddressesCreated != 0 {
		t.Fatalf("dry-run result = %+v", result)
	}
	if contacts, err := service.ListContacts(context.Background(), "tenant-a", ContactFilter{}); err != nil || len(contacts) != 0 {
		t.Fatalf("dry-run mutated contacts = %+v, %v", contacts, err)
	}
}

func TestExportCSVRoundTrip(t *testing.T) {
	sourceService := NewContactService(NewInMemoryContactRepository())
	contact := &Contact{PublicID: "contact-round-trip", DisplayName: "Round Trip", Status: ContactStatusActive, Source: ContactSourceAPI}
	if err := sourceService.CreateContact(context.Background(), "tenant-source", contact); err != nil {
		t.Fatalf("CreateContact() error = %v", err)
	}
	address, err := NewContactAddress("tenant-source", contact.ID, AddressKindEmail, "", "RoundTrip@Example.com")
	if err != nil {
		t.Fatalf("NewContactAddress() error = %v", err)
	}
	if err := sourceService.CreateContactAddress(context.Background(), "tenant-source", &address); err != nil {
		t.Fatalf("CreateContactAddress() error = %v", err)
	}

	var exported strings.Builder
	if err := ExportCSV(context.Background(), sourceService, "tenant-source", &exported, CSVExportOptions{}); err != nil {
		t.Fatalf("ExportCSV() error = %v", err)
	}
	targetService := NewContactService(NewInMemoryContactRepository())
	result, err := ImportCSV(context.Background(), targetService, "tenant-target", strings.NewReader(exported.String()), CSVImportOptions{})
	if err != nil || result.RowsImported != 1 {
		t.Fatalf("round-trip import = %+v, %v; csv=%q", result, err, exported.String())
	}
	copied, err := targetService.GetContact(context.Background(), "tenant-target", "contact-round-trip")
	if err != nil || copied.DisplayName != "Round Trip" || copied.Source != ContactSourceAPI {
		t.Fatalf("round-trip contact = %+v, %v", copied, err)
	}
	addresses, err := targetService.ListContactAddressesByContactID(context.Background(), "tenant-target", copied.ID)
	if err != nil || len(addresses) != 1 || addresses[0].NormalizedValue != "roundtrip@example.com" {
		t.Fatalf("round-trip addresses = %+v, %v", addresses, err)
	}
}

func TestExportCSVEnforcesFieldLimitWithoutPII(t *testing.T) {
	service := NewContactService(NewInMemoryContactRepository())
	contact := &Contact{PublicID: "field-limit", DisplayName: "private display name", Status: ContactStatusActive, Source: ContactSourceManual}
	if err := service.CreateContact(context.Background(), "tenant-a", contact); err != nil {
		t.Fatalf("CreateContact() error = %v", err)
	}

	var exported strings.Builder
	err := ExportCSV(context.Background(), service, "tenant-a", &exported, CSVExportOptions{MaxFieldBytes: 3})
	if !errors.Is(err, ErrCSVLimit) || strings.Contains(err.Error(), "private display name") {
		t.Fatalf("field limit error = %v", err)
	}
}

func TestCSVImportExportTenantIsolation(t *testing.T) {
	service := NewContactService(NewInMemoryContactRepository())
	input := strings.Join([]string{
		strings.Join(CSVHeader, ","),
		"shared,Shared,active,import,email,,same@example.com",
	}, "\n")
	for _, tenant := range []string{"tenant-a", "tenant-b"} {
		result, err := ImportCSV(context.Background(), service, tenant, strings.NewReader(input), CSVImportOptions{})
		if err != nil || result.RowsImported != 1 || result.DuplicateRows != 0 {
			t.Fatalf("tenant %s import = %+v, %v", tenant, result, err)
		}
	}
	if _, err := service.GetContact(context.Background(), "tenant-a", "shared"); err != nil {
		t.Fatalf("tenant-a contact lookup = %v", err)
	}
	if _, err := service.GetContact(context.Background(), "tenant-b", "shared"); err != nil {
		t.Fatalf("tenant-b contact lookup = %v", err)
	}
	var tenantA strings.Builder
	if err := ExportCSV(context.Background(), service, "tenant-a", &tenantA, CSVExportOptions{}); err != nil {
		t.Fatalf("tenant-a export = %v", err)
	}
	if strings.Count(tenantA.String(), "same@example.com") != 1 {
		t.Fatalf("tenant-a export = %q", tenantA.String())
	}
}
