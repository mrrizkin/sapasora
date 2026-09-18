package contact

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

// CSVHeader is the stable, provider-neutral CSV schema used by the contact
// import/export foundation. Address values are intentionally present only in
// the CSV payload; diagnostics never echo them.
var CSVHeader = []string{
	"contact_public_id",
	"display_name",
	"status",
	"source",
	"address_kind",
	"address_namespace",
	"address_value",
}

const (
	DefaultCSVMaxRows       = 10_000
	DefaultCSVMaxFieldBytes = 4 * 1024
)

var (
	ErrCSVSchema    = errors.New("contact csv schema is invalid")
	ErrCSVLimit     = errors.New("contact csv limit exceeded")
	ErrCSVOperation = errors.New("contact csv operation failed")
)

// CSVImportOptions bounds a streaming import. Zero limits select safe
// defaults. DryRun validates and reports the same rows without mutations.
type CSVImportOptions struct {
	DryRun        bool
	MaxRows       int
	MaxFieldBytes int
}

// CSVExportOptions bounds the records and field sizes emitted. Zero selects
// safe defaults. Export writes directly to the supplied writer and never
// buffers the complete CSV document.
type CSVExportOptions struct {
	MaxRows       int
	MaxFieldBytes int
}

// CSVImportIssue is a safe row-level diagnostic. It contains only a row,
// schema field, and stable code; it never contains a raw or normalized value.
type CSVImportIssue struct {
	Row   int
	Field string
	Code  string
}

func (i CSVImportIssue) Error() string {
	return fmt.Sprintf("contact csv row %d field %s: %s", i.Row, i.Field, i.Code)
}

// CSVImportResult reports bounded import work without returning imported PII.
type CSVImportResult struct {
	DryRun           bool
	RowsRead         int
	RowsAccepted     int
	RowsImported     int
	RowsWouldImport  int
	RowsSkipped      int
	DuplicateRows    int
	ContactsCreated  int
	AddressesCreated int
	Issues           []CSVImportIssue
}

// CSVError is returned for a fatal schema, limit, or repository operation
// failure. Its message contains only safe structural references.
type CSVError struct {
	Kind  error
	Row   int
	Field string
	Code  string
}

func (e *CSVError) Error() string {
	if e == nil {
		return "contact csv error"
	}
	if e.Row > 0 {
		return fmt.Sprintf("contact csv row %d field %s: %s", e.Row, e.Field, e.Code)
	}
	return fmt.Sprintf("contact csv: %s", e.Code)
}

func (e *CSVError) Unwrap() error { return e.Kind }

// ImportCSV validates and imports the CSV stream for exactly tenantID. It is
// deliberately transport-neutral: HTTP upload limits, authorization, and job
// orchestration belong to the integration boundary.
func ImportCSV(ctx context.Context, service ContactService, tenantID string, source io.Reader, options CSVImportOptions) (CSVImportResult, error) {
	result := CSVImportResult{DryRun: options.DryRun}
	if err := validateCSVInputs(ctx, service, tenantID, source); err != nil {
		return result, err
	}
	maxRows, maxFieldBytes, err := csvLimits(options.MaxRows, options.MaxFieldBytes)
	if err != nil {
		return result, err
	}

	reader := csv.NewReader(source)
	reader.FieldsPerRecord = -1
	reader.ReuseRecord = true
	header, err := reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return result, &CSVError{Kind: ErrCSVSchema, Field: "header", Code: "missing"}
		}
		return result, &CSVError{Kind: ErrCSVSchema, Field: "header", Code: "malformed"}
	}
	if issue := validateCSVFields(header, 1, maxFieldBytes); issue != nil {
		return result, &CSVError{Kind: ErrCSVLimit, Row: issue.Row, Field: issue.Field, Code: issue.Code}
	}
	if !sameCSVHeader(header) {
		return result, &CSVError{Kind: ErrCSVSchema, Row: 1, Field: "header", Code: "unexpected"}
	}

	seenIdentities := make(map[string]struct{})
	localContacts := make(map[string]*Contact)
	plannedContacts := make(map[string]struct{})
	for {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		record, readErr := reader.Read()
		if errors.Is(readErr, io.EOF) {
			break
		}
		rowNumber := result.RowsRead + 2
		if result.RowsRead >= maxRows {
			return result, &CSVError{Kind: ErrCSVLimit, Row: rowNumber, Field: "record", Code: "max_rows"}
		}
		result.RowsRead++

		if readErr != nil && !errors.Is(readErr, csv.ErrFieldCount) {
			addCSVIssue(&result, rowNumber, "record", "malformed")
			continue
		}
		if len(record) != len(CSVHeader) {
			addCSVIssue(&result, rowNumber, "record", "field_count")
			continue
		}
		if issue := validateCSVFields(record, rowNumber, maxFieldBytes); issue != nil {
			result.RowsSkipped++
			result.Issues = append(result.Issues, *issue)
			continue
		}
		row, issue := parseCSVImportRow(record, rowNumber)
		if issue != nil {
			result.RowsSkipped++
			result.Issues = append(result.Issues, *issue)
			continue
		}

		var identity AddressIdentity
		hasAddress := row.kind != ""
		if hasAddress {
			identity, err = NormalizeAddress(AddressKind(row.kind), row.namespace, row.value)
			if err != nil {
				result.RowsSkipped++
				result.Issues = append(result.Issues, CSVImportIssue{Row: rowNumber, Field: "address_value", Code: "invalid_address"})
				continue
			}
			if _, exists := seenIdentities[identity.Key()]; exists {
				addDuplicateIssue(&result, rowNumber)
				continue
			}
			existing, lookupErr := serviceIdentityExists(ctx, service, tenantID, identity)
			if lookupErr != nil {
				return result, lookupErr
			}
			if existing {
				addDuplicateIssue(&result, rowNumber)
				seenIdentities[identity.Key()] = struct{}{}
				continue
			}
		}

		contact, contactCreated, err := importContact(ctx, service, tenantID, row, options.DryRun, localContacts, plannedContacts)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return result, err
			}
			result.RowsSkipped++
			result.Issues = append(result.Issues, CSVImportIssue{Row: rowNumber, Field: "contact_public_id", Code: "contact_operation_failed"})
			continue
		}
		if contactCreated {
			result.ContactsCreated++
		}
		if hasAddress && !options.DryRun {
			address, addressErr := NewContactAddress(tenantID, contact.ID, identity.Kind, identity.Namespace, row.value)
			if addressErr != nil {
				result.RowsSkipped++
				result.Issues = append(result.Issues, CSVImportIssue{Row: rowNumber, Field: "address_value", Code: "invalid_address"})
				continue
			}
			address.Source = ContactSourceImport
			if err := service.CreateContactAddress(ctx, tenantID, &address); err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return result, err
				}
				if errors.Is(err, ErrConflict) {
					addDuplicateIssue(&result, rowNumber)
					seenIdentities[identity.Key()] = struct{}{}
					continue
				}
				result.RowsSkipped++
				result.Issues = append(result.Issues, CSVImportIssue{Row: rowNumber, Field: "address_value", Code: "address_operation_failed"})
				continue
			}
			result.AddressesCreated++
		}
		if hasAddress {
			seenIdentities[identity.Key()] = struct{}{}
		}
		result.RowsAccepted++
		if options.DryRun {
			result.RowsWouldImport++
		} else {
			result.RowsImported++
		}
	}
	return result, nil
}

// ExportCSV emits tenant-scoped contacts and addresses in CSVHeader order.
// Contacts without addresses receive one row with empty address fields.
func ExportCSV(ctx context.Context, service ContactService, tenantID string, destination io.Writer, options CSVExportOptions) error {
	if err := validateCSVInputs(ctx, service, tenantID, destination); err != nil {
		return err
	}
	maxRows, maxFieldBytes, err := csvLimits(options.MaxRows, options.MaxFieldBytes)
	if err != nil {
		return err
	}
	contacts, err := service.ListContacts(ctx, tenantID, ContactFilter{})
	if err != nil {
		return csvOperationError(err)
	}
	sort.SliceStable(contacts, func(i, j int) bool {
		if contacts[i].ID != contacts[j].ID {
			return contacts[i].ID < contacts[j].ID
		}
		return contacts[i].PublicID < contacts[j].PublicID
	})

	writer := csv.NewWriter(destination)
	if err := writer.Write(CSVHeader); err != nil {
		return csvOperationError(err)
	}
	rows := 0
	for _, contact := range contacts {
		if err := ctx.Err(); err != nil {
			return err
		}
		addresses, listErr := service.ListContactAddressesByContactID(ctx, tenantID, contact.ID)
		if listErr != nil {
			return csvOperationError(listErr)
		}
		sort.SliceStable(addresses, func(i, j int) bool {
			if addresses[i].ID != addresses[j].ID {
				return addresses[i].ID < addresses[j].ID
			}
			return addresses[i].PublicID < addresses[j].PublicID
		})
		if len(addresses) == 0 {
			if rows >= maxRows {
				return &CSVError{Kind: ErrCSVLimit, Field: "record", Code: "max_rows"}
			}
			if err := writeCSVContactRow(writer, contact, nil, rows+2, maxFieldBytes); err != nil {
				return csvOperationError(err)
			}
			rows++
		} else {
			for _, address := range addresses {
				if err := ctx.Err(); err != nil {
					return err
				}
				if rows >= maxRows {
					return &CSVError{Kind: ErrCSVLimit, Field: "record", Code: "max_rows"}
				}
				if err := writeCSVContactRow(writer, contact, address, rows+2, maxFieldBytes); err != nil {
					return csvOperationError(err)
				}
				rows++
			}
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return csvOperationError(err)
	}
	return nil
}

// NewCSVImporter/NewCSVExporter provide an object form for application
// composition while keeping the functional API convenient for tests.
type CSVImporter struct{ service ContactService }
type CSVExporter struct{ service ContactService }

func NewCSVImporter(service ContactService) *CSVImporter { return &CSVImporter{service: service} }
func NewCSVExporter(service ContactService) *CSVExporter { return &CSVExporter{service: service} }

func (i *CSVImporter) Import(ctx context.Context, tenantID string, source io.Reader, options CSVImportOptions) (CSVImportResult, error) {
	if i == nil {
		return CSVImportResult{}, &CSVError{Kind: ErrCSVOperation, Code: "importer_unavailable"}
	}
	return ImportCSV(ctx, i.service, tenantID, source, options)
}

func (e *CSVExporter) Export(ctx context.Context, tenantID string, destination io.Writer, options CSVExportOptions) error {
	if e == nil {
		return &CSVError{Kind: ErrCSVOperation, Code: "exporter_unavailable"}
	}
	return ExportCSV(ctx, e.service, tenantID, destination, options)
}

type csvImportRow struct {
	publicID    string
	displayName string
	status      ContactStatus
	source      ContactSource
	kind        string
	namespace   string
	value       string
}

func validateCSVInputs(ctx context.Context, service ContactService, tenantID string, stream any) error {
	if ctx == nil {
		return &CSVError{Kind: ErrCSVOperation, Code: "context_unavailable"}
	}
	if service == nil || stream == nil {
		return &CSVError{Kind: ErrCSVOperation, Code: "boundary_unavailable"}
	}
	if strings.TrimSpace(tenantID) == "" {
		return &CSVError{Kind: ErrCSVOperation, Code: "tenant_required"}
	}
	return nil
}

func csvLimits(maxRows, maxFieldBytes int) (int, int, error) {
	if maxRows == 0 {
		maxRows = DefaultCSVMaxRows
	}
	if maxFieldBytes == 0 {
		maxFieldBytes = DefaultCSVMaxFieldBytes
	}
	if maxRows < 1 || maxFieldBytes < 1 {
		return 0, 0, &CSVError{Kind: ErrCSVLimit, Field: "options", Code: "invalid_limit"}
	}
	return maxRows, maxFieldBytes, nil
}

func sameCSVHeader(header []string) bool {
	if len(header) != len(CSVHeader) {
		return false
	}
	for index, value := range header {
		if index == 0 {
			value = strings.TrimPrefix(value, "\ufeff")
		}
		if value != CSVHeader[index] {
			return false
		}
	}
	return true
}

func validateCSVFields(record []string, row, maxFieldBytes int) *CSVImportIssue {
	for index, value := range record {
		if len(value) > maxFieldBytes {
			field := "record"
			if index < len(CSVHeader) {
				field = CSVHeader[index]
			}
			return &CSVImportIssue{Row: row, Field: field, Code: "max_field_bytes"}
		}
	}
	return nil
}

func parseCSVImportRow(record []string, row int) (csvImportRow, *CSVImportIssue) {
	value := csvImportRow{
		publicID:    strings.TrimSpace(record[0]),
		displayName: record[1],
		kind:        strings.ToLower(strings.TrimSpace(record[4])),
		namespace:   strings.TrimSpace(record[5]),
		value:       record[6],
	}
	value.status = ContactStatus(strings.TrimSpace(record[2]))
	if value.status == "" {
		value.status = ContactStatusActive
	}
	if !value.status.Valid() {
		return csvImportRow{}, &CSVImportIssue{Row: row, Field: "status", Code: "invalid_status"}
	}
	value.source = ContactSource(strings.TrimSpace(record[3]))
	if value.source == "" {
		value.source = ContactSourceImport
	}
	if !value.source.Valid() {
		return csvImportRow{}, &CSVImportIssue{Row: row, Field: "source", Code: "invalid_source"}
	}
	if value.kind == "" && (strings.TrimSpace(value.namespace) != "" || strings.TrimSpace(value.value) != "") {
		return csvImportRow{}, &CSVImportIssue{Row: row, Field: "address_kind", Code: "address_fields_incomplete"}
	}
	if value.kind != "" {
		kind := AddressKind(value.kind)
		if !kind.Valid() {
			return csvImportRow{}, &CSVImportIssue{Row: row, Field: "address_kind", Code: "invalid_kind"}
		}
		if strings.TrimSpace(value.value) == "" {
			return csvImportRow{}, &CSVImportIssue{Row: row, Field: "address_value", Code: "required"}
		}
	}
	return value, nil
}

func serviceIdentityExists(ctx context.Context, service ContactService, tenantID string, identity AddressIdentity) (bool, error) {
	// IsAddressSendable is an existing tenant-scoped lookup that does not expose
	// identity values. Its boolean is irrelevant here; a nil error means the
	// identity exists even when consent makes it unsendable.
	_, err := service.IsAddressSendable(ctx, tenantID, identity)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, ErrNotFound) {
		return false, nil
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false, err
	}
	return false, csvOperationError(err)
}

func importContact(ctx context.Context, service ContactService, tenantID string, row csvImportRow, dryRun bool, local map[string]*Contact, planned map[string]struct{}) (*Contact, bool, error) {
	if row.publicID != "" {
		if contact, ok := local[row.publicID]; ok {
			return contact, false, nil
		}
		if _, ok := planned[row.publicID]; ok {
			return &Contact{TenantID: tenantID, PublicID: row.publicID, ID: 1}, false, nil
		}
		contact, err := service.GetContactByPublicID(ctx, tenantID, row.publicID)
		if err == nil {
			local[row.publicID] = contact
			return contact, false, nil
		}
		if !errors.Is(err, ErrNotFound) {
			return nil, false, err
		}
	}
	if dryRun {
		if row.publicID != "" {
			planned[row.publicID] = struct{}{}
		}
		return &Contact{TenantID: tenantID, PublicID: row.publicID, ID: 1}, false, nil
	}
	contact := &Contact{TenantID: tenantID, PublicID: row.publicID, DisplayName: row.displayName, Status: row.status, Source: row.source}
	if err := service.CreateContact(ctx, tenantID, contact); err != nil {
		return nil, false, err
	}
	if contact.PublicID != "" {
		local[contact.PublicID] = contact
	}
	return contact, true, nil
}

func writeCSVContactRow(writer *csv.Writer, contact *Contact, address *ContactAddress, row, maxFieldBytes int) error {
	record := []string{contact.PublicID, contact.DisplayName, contact.Status.String(), contact.Source.String(), "", "", ""}
	if address != nil {
		record[4] = address.Kind.String()
		record[5] = address.Namespace
		record[6] = address.NormalizedValue
	}
	if issue := validateCSVFields(record, row, maxFieldBytes); issue != nil {
		return &CSVError{Kind: ErrCSVLimit, Row: issue.Row, Field: issue.Field, Code: issue.Code}
	}
	return writer.Write(record)
}

func addCSVIssue(result *CSVImportResult, row int, field, code string) {
	result.RowsSkipped++
	result.Issues = append(result.Issues, CSVImportIssue{Row: row, Field: field, Code: code})
}

func addDuplicateIssue(result *CSVImportResult, row int) {
	result.RowsSkipped++
	result.DuplicateRows++
	result.Issues = append(result.Issues, CSVImportIssue{Row: row, Field: "address_value", Code: "duplicate_identity"})
}

func csvOperationError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, ErrCSVLimit) {
		return err
	}
	return &CSVError{Kind: ErrCSVOperation, Code: "contact_operation_failed"}
}
