package conversation

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	// DefaultInboxPageSize is deliberately small enough for a transport adapter
	// to return without an unbounded response.
	DefaultInboxPageSize = 25
	// MaxInboxPageSize is the hard upper bound for every inbox page.
	MaxInboxPageSize   = 100
	maxInboxQueryValue = 255
	maxInboxCursorSize = 512
)

// InboxListState is an adapter-friendly outcome for an inbox list operation.
// Loading is exposed for adapters that model an asynchronous request; the
// synchronous use case returns Ready, Empty, or Error.
type InboxListState string

const (
	InboxListStateLoading InboxListState = "loading"
	InboxListStateReady   InboxListState = "ready"
	InboxListStateEmpty   InboxListState = "empty"
	InboxListStateError   InboxListState = "error"

	InboxListLoading     = InboxListStateLoading
	InboxListReady       = InboxListStateReady
	InboxListEmpty       = InboxListStateEmpty
	InboxListErrorResult = InboxListStateError
)

// InboxListError is a safe, stable error projection. The repository error is
// intentionally not retained here because it may contain SQL, provider, or
// PII-bearing details unsuitable for transport or diagnostics.
type InboxListError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// InboxListResult is the complete result envelope consumed by future HTTP/UI
// adapters. Items are always non-nil on successful outcomes.
type InboxListResult struct {
	State      InboxListState  `json:"state"`
	Items      []*Conversation `json:"items"`
	NextCursor string          `json:"next_cursor,omitempty"`
	HasMore    bool            `json:"has_more"`
	Error      *InboxListError `json:"error,omitempty"`
}

// NewInboxListLoadingResult creates the explicit initial state used by an
// asynchronous transport adapter.
func NewInboxListLoadingResult() InboxListResult {
	return InboxListResult{State: InboxListStateLoading, Items: make([]*Conversation, 0)}
}

// InboxListFilter contains transport-neutral inbox filters. Opaque public IDs
// are used instead of names or addresses so searching this aggregate cannot
// accidentally become a PII search.
type InboxListFilter struct {
	Search           string
	ChannelPublicID  string
	Status           *ConversationStatus
	AssigneePublicID string
	TeamPublicID     string
	Priority         *ConversationPriority
	Tag              string
	UnreadOnly       bool
	SLABreachedOnly  bool
	// SLAOverdueOnly is a readable compatibility alias for the same filter.
	SLAOverdueOnly bool
}

// InboxListQuery is the parsed, bounded query portion of an inbox request.
// It is intentionally independent of tenant selection; the authenticated
// integration must provide that scope separately.
type InboxListQuery struct {
	InboxListFilter
	Cursor string
	Limit  int
}

// InboxListRequest binds a parsed query to one tenant/workspace scope.
type InboxListRequest struct {
	TenantID    string
	WorkspaceID string
	Filter      InboxListFilter
	Cursor      string
	Limit       int
	// AsOf makes SLA filtering deterministic for callers and tests. A zero
	// value uses the current UTC time for that invocation only.
	AsOf time.Time
}

// Request creates a tenant-bound request from parsed query parameters.
func (q InboxListQuery) Request(tenantID, workspaceID string) InboxListRequest {
	return InboxListRequest{
		TenantID: tenantID, WorkspaceID: workspaceID,
		Filter: q.InboxListFilter, Cursor: q.Cursor, Limit: q.Limit,
	}
}

// ParseInboxListQuery parses an allowlisted URL query without reflection. It
// rejects unknown keys, duplicate values, explicit empty values, malformed
// enums/booleans, and values larger than their bounded field size.
func ParseInboxListQuery(values url.Values) (InboxListQuery, error) {
	query := InboxListQuery{Limit: DefaultInboxPageSize}
	allowed := map[string]struct{}{
		"search": {}, "channel": {}, "channel_id": {}, "status": {},
		"assignee": {}, "assignee_id": {}, "team": {}, "team_id": {},
		"priority": {}, "tag": {}, "unread": {}, "sla": {}, "sla_breached": {},
		"cursor": {}, "limit": {},
	}
	for key, raw := range values {
		if _, ok := allowed[key]; !ok {
			return InboxListQuery{}, fmt.Errorf("%w: unsupported query parameter %q", ErrInvalidInboxQuery, key)
		}
		if len(raw) != 1 || raw[0] == "" {
			return InboxListQuery{}, fmt.Errorf("%w: query parameter %q must have one non-empty value", ErrInvalidInboxQuery, key)
		}
		if len(raw[0]) > maxInboxQueryValue || hasControlCharacters(raw[0]) {
			return InboxListQuery{}, fmt.Errorf("%w: query parameter %q is invalid", ErrInvalidInboxQuery, key)
		}
	}

	query.Search = strings.TrimSpace(values.Get("search"))
	query.InboxListFilter.ChannelPublicID = firstQueryValue(values, "channel", "channel_id")
	query.InboxListFilter.AssigneePublicID = firstQueryValue(values, "assignee", "assignee_id")
	query.InboxListFilter.TeamPublicID = firstQueryValue(values, "team", "team_id")
	query.InboxListFilter.Tag = strings.ToLower(strings.TrimSpace(values.Get("tag")))
	if hasMultipleInboxAliases(values, "channel", "channel_id") || hasMultipleInboxAliases(values, "assignee", "assignee_id") || hasMultipleInboxAliases(values, "team", "team_id") || hasMultipleInboxAliases(values, "sla", "sla_breached") {
		return InboxListQuery{}, fmt.Errorf("%w: use one name for each filter", ErrInvalidInboxQuery)
	}
	if value := strings.ToLower(strings.TrimSpace(values.Get("status"))); value != "" {
		status := ConversationStatus(value)
		if !status.Valid() {
			return InboxListQuery{}, fmt.Errorf("%w: invalid status", ErrInvalidInboxQuery)
		}
		query.Status = &status
	}
	if value := strings.ToLower(strings.TrimSpace(values.Get("priority"))); value != "" {
		priority := ConversationPriority(value)
		if !priority.Valid() {
			return InboxListQuery{}, fmt.Errorf("%w: invalid priority", ErrInvalidInboxQuery)
		}
		query.Priority = &priority
	}
	var err error
	if query.UnreadOnly, err = parseInboxBool(values, "unread"); err != nil {
		return InboxListQuery{}, err
	}
	if query.SLABreachedOnly, err = parseSLAQuery(values); err != nil {
		return InboxListQuery{}, err
	}
	query.Cursor = strings.TrimSpace(values.Get("cursor"))
	if query.Cursor != "" && len(query.Cursor) > maxInboxCursorSize {
		return InboxListQuery{}, fmt.Errorf("%w: cursor is too long", ErrInvalidInboxQuery)
	}
	if value := values.Get("limit"); value != "" {
		query.Limit, err = strconv.Atoi(value)
		if err != nil || query.Limit < 1 || query.Limit > MaxInboxPageSize {
			return InboxListQuery{}, fmt.Errorf("%w: limit must be between 1 and %d", ErrInvalidInboxQuery, MaxInboxPageSize)
		}
	}
	return query, nil
}

// ParseInboxQuery is a convenience boundary for adapters that have a raw
// query string but do not yet have an HTTP framework dependency.
func ParseInboxQuery(rawQuery string) (InboxListQuery, error) {
	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return InboxListQuery{}, fmt.Errorf("%w: malformed query", ErrInvalidInboxQuery)
	}
	return ParseInboxListQuery(values)
}

func hasMultipleInboxAliases(values url.Values, names ...string) bool {
	count := 0
	for _, name := range names {
		if _, ok := values[name]; ok {
			count++
		}
	}
	return count > 1
}

func firstQueryValue(values url.Values, names ...string) string {
	for _, name := range names {
		if value := strings.TrimSpace(values.Get(name)); value != "" {
			return value
		}
	}
	return ""
}

func parseInboxBool(values url.Values, key string) (bool, error) {
	value := strings.ToLower(strings.TrimSpace(values.Get(key)))
	if value == "" {
		return false, nil
	}
	switch value {
	case "true", "1":
		return true, nil
	case "false", "0":
		return false, nil
	default:
		return false, fmt.Errorf("%w: %s must be true or false", ErrInvalidInboxQuery, key)
	}
}

func parseSLAQuery(values url.Values) (bool, error) {
	if value := strings.ToLower(strings.TrimSpace(values.Get("sla"))); value != "" {
		if value != "breached" && value != "overdue" {
			return false, fmt.Errorf("%w: invalid sla filter", ErrInvalidInboxQuery)
		}
		return true, nil
	}
	return parseInboxBool(values, "sla_breached")
}

func hasControlCharacters(value string) bool {
	for _, character := range value {
		if unicode.IsControl(character) {
			return true
		}
	}
	return false
}

// InboxListUseCase is the tenant-scoped list boundary for transport adapters.
type InboxListUseCase interface {
	List(context.Context, InboxListRequest) InboxListResult
	ListConversations(context.Context, InboxListRequest) InboxListResult
}

// InboxListService applies filtering, deterministic ordering, and bounded
// cursor pagination above the conversation repository.
type InboxListService struct {
	repository Repository
}

var _ InboxListUseCase = (*InboxListService)(nil)

// NewInboxListUseCase constructs the inbox-list use case.
func NewInboxListUseCase(repository Repository) InboxListUseCase {
	return &InboxListService{repository: repository}
}

// NewInboxListService is the concrete-name constructor for application wiring.
func NewInboxListService(repository Repository) *InboxListService {
	return &InboxListService{repository: repository}
}

// List executes one bounded, tenant/workspace-scoped inbox list operation.
func (s *InboxListService) List(ctx context.Context, request InboxListRequest) InboxListResult {
	if s == nil || s.repository == nil {
		return inboxListError("repository_unavailable", "inbox list is unavailable")
	}
	if ctx == nil {
		return inboxListError("invalid_request", "request context is required")
	}
	scope, err := validateScope(request.TenantID, request.WorkspaceID)
	if err != nil {
		return inboxListError("invalid_request", "tenant/workspace scope is required")
	}
	if err := validateInboxFilter(request.Filter); err != nil {
		return inboxListError("invalid_request", "inbox filters are invalid")
	}
	limit := request.Limit
	if limit == 0 {
		limit = DefaultInboxPageSize
	}
	if limit < 1 || limit > MaxInboxPageSize {
		return inboxListError("invalid_request", "inbox page size is out of bounds")
	}
	asOf := request.AsOf
	if asOf.IsZero() {
		asOf = time.Now().UTC()
	}
	cursorAsOf := time.Time{}
	if request.Filter.SLABreachedOnly || request.Filter.SLAOverdueOnly {
		// SLA membership depends on the evaluation instant. Bind that instant
		// to the cursor so pages cannot silently change their filter window.
		cursorAsOf = asOf
	}
	cursor, err := decodeInboxCursor(request.Cursor, scope, request.Filter, cursorAsOf)
	if err != nil {
		return inboxListError("invalid_cursor", "cursor is invalid or no longer matches this list")
	}

	repositoryFilter := conversationFilterForInbox(request.Filter, asOf)
	values, err := s.repository.ListConversations(ctx, scope, repositoryFilter)
	if err != nil {
		// Never expose repository/provider details through this result.
		return inboxListError("repository_error", "inbox list could not be loaded")
	}
	items := make([]*Conversation, 0, len(values))
	for _, value := range values {
		if value == nil || value.DeletedAt != nil || value.ScopeID() != scope || !matchesInboxSearch(value, request.Filter.Search) {
			continue
		}
		items = append(items, cloneConversation(value))
	}
	sort.SliceStable(items, func(i, j int) bool {
		return inboxConversationBefore(items[i], items[j])
	})
	if cursor != nil {
		items = filterAfterInboxCursor(items, cursor)
	}

	hasMore := len(items) > limit
	if hasMore {
		items = items[:limit]
	}
	result := InboxListResult{State: InboxListStateReady, Items: items, HasMore: hasMore}
	if len(items) == 0 {
		result.State = InboxListStateEmpty
		return result
	}
	if hasMore {
		result.NextCursor = encodeInboxCursor(items[len(items)-1], scope, request.Filter, cursorAsOf)
	}
	return result
}

// ListConversations is the aggregate-qualified alias used by callers that
// keep repository and use-case method names aligned.
func (s *InboxListService) ListConversations(ctx context.Context, request InboxListRequest) InboxListResult {
	return s.List(ctx, request)
}

// Execute is a descriptive alias for callers that name application boundaries
// as commands/use cases rather than services.
func (s *InboxListService) Execute(ctx context.Context, request InboxListRequest) InboxListResult {
	return s.List(ctx, request)
}

func validateInboxFilter(filter InboxListFilter) error {
	if len(filter.Search) > maxInboxQueryValue || hasControlCharacters(filter.Search) {
		return ErrInvalidInboxQuery
	}
	if len(filter.ChannelPublicID) > maxInboxQueryValue || len(filter.AssigneePublicID) > maxInboxQueryValue || len(filter.TeamPublicID) > maxInboxQueryValue || len(filter.Tag) > maxInboxQueryValue {
		return ErrInvalidInboxQuery
	}
	for _, value := range []string{filter.ChannelPublicID, filter.AssigneePublicID, filter.TeamPublicID, filter.Tag} {
		if hasControlCharacters(value) {
			return ErrInvalidInboxQuery
		}
	}
	if filter.Status != nil && !filter.Status.Valid() {
		return ErrInvalidInboxQuery
	}
	if filter.Priority != nil && !filter.Priority.Valid() {
		return ErrInvalidInboxQuery
	}
	return nil
}

func conversationFilterForInbox(filter InboxListFilter, asOf time.Time) ConversationFilter {
	conversationFilter := ConversationFilter{
		UnreadOnly: filter.UnreadOnly,
	}
	if filter.Status != nil {
		status := *filter.Status
		conversationFilter.Status = &status
	}
	if filter.Priority != nil {
		priority := *filter.Priority
		conversationFilter.Priority = &priority
	}
	if filter.ChannelPublicID != "" {
		value := filter.ChannelPublicID
		conversationFilter.ChannelPublicID = &value
	}
	if filter.AssigneePublicID != "" {
		value := filter.AssigneePublicID
		conversationFilter.AssigneePublicID = &value
	}
	if filter.TeamPublicID != "" {
		value := filter.TeamPublicID
		conversationFilter.TeamPublicID = &value
	}
	if filter.Tag != "" {
		value := strings.ToLower(strings.TrimSpace(filter.Tag))
		conversationFilter.Tag = &value
	}
	if filter.SLABreachedOnly || filter.SLAOverdueOnly {
		conversationFilter.SLAOverdueAt = &asOf
	}
	return conversationFilter
}

func matchesInboxSearch(value *Conversation, search string) bool {
	search = strings.ToLower(strings.TrimSpace(search))
	if search == "" {
		return true
	}
	candidates := []string{value.PublicID, value.ChannelPublicID, value.ContactPublicID, value.AssigneePublicID, value.TeamPublicID, string(value.Status), string(value.Priority)}
	for _, tag := range value.Tags {
		candidates = append(candidates, tag)
	}
	for _, candidate := range candidates {
		if strings.Contains(strings.ToLower(candidate), search) {
			return true
		}
	}
	return false
}

type inboxCursor struct {
	Version      uint8  `json:"v"`
	ActivityUnix int64  `json:"a"`
	PublicID     string `json:"p"`
	ScopeHash    string `json:"s"`
	FilterHash   string `json:"f"`
	AsOfUnix     int64  `json:"t,omitempty"`
}

func encodeInboxCursor(value *Conversation, scope string, filter InboxListFilter, asOf time.Time) string {
	cursor := inboxCursor{
		Version: 1, ActivityUnix: value.LastActivityAt.UnixNano(), PublicID: value.PublicID,
		ScopeHash: inboxHash(scope), FilterHash: inboxFilterHash(filter), AsOfUnix: inboxAsOfUnix(asOf),
	}
	encoded, _ := json.Marshal(cursor)
	return base64.RawURLEncoding.EncodeToString(encoded)
}

func decodeInboxCursor(encoded, scope string, filter InboxListFilter, asOf time.Time) (*inboxCursor, error) {
	if encoded == "" {
		return nil, nil
	}
	if len(encoded) > maxInboxCursorSize {
		return nil, ErrInvalidInboxCursor
	}
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, ErrInvalidInboxCursor
	}
	var cursor inboxCursor
	decoder := json.NewDecoder(strings.NewReader(string(decoded)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cursor); err != nil || cursor.Version != 1 || cursor.PublicID == "" || len(cursor.PublicID) > maxInboxQueryValue || cursor.ScopeHash != inboxHash(scope) || cursor.FilterHash != inboxFilterHash(filter) || cursor.AsOfUnix != inboxAsOfUnix(asOf) {
		return nil, ErrInvalidInboxCursor
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return nil, ErrInvalidInboxCursor
	}
	return &cursor, nil
}

func inboxAsOfUnix(value time.Time) int64 {
	if value.IsZero() {
		return 0
	}
	return value.UnixNano()
}

func inboxHash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func inboxFilterHash(filter InboxListFilter) string {
	canonical := struct {
		Search, Channel, Assignee, Team, Tag string
		Status                               ConversationStatus
		Priority                             ConversationPriority
		Unread, SLABreached, SLAOverdue      bool
	}{
		Search: strings.ToLower(strings.TrimSpace(filter.Search)), Channel: filter.ChannelPublicID,
		Assignee: filter.AssigneePublicID, Team: filter.TeamPublicID, Tag: strings.ToLower(strings.TrimSpace(filter.Tag)),
		Unread: filter.UnreadOnly, SLABreached: filter.SLABreachedOnly, SLAOverdue: filter.SLAOverdueOnly,
	}
	if filter.Status != nil {
		canonical.Status = *filter.Status
	}
	if filter.Priority != nil {
		canonical.Priority = *filter.Priority
	}
	encoded, _ := json.Marshal(canonical)
	sum := sha256.Sum256(encoded)
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func inboxConversationBefore(left, right *Conversation) bool {
	if !left.LastActivityAt.Equal(right.LastActivityAt) {
		return left.LastActivityAt.After(right.LastActivityAt)
	}
	if left.PublicID != right.PublicID {
		return left.PublicID > right.PublicID
	}
	return left.ID > right.ID
}

func filterAfterInboxCursor(items []*Conversation, cursor *inboxCursor) []*Conversation {
	result := items[:0]
	cursorTime := time.Unix(0, cursor.ActivityUnix)
	for _, value := range items {
		if value.LastActivityAt.Before(cursorTime) || (value.LastActivityAt.Equal(cursorTime) && value.PublicID < cursor.PublicID) {
			result = append(result, value)
		}
	}
	return result
}

func inboxListError(code, message string) InboxListResult {
	return InboxListResult{State: InboxListStateError, Items: make([]*Conversation, 0), Error: &InboxListError{Code: code, Message: message}}
}

// IsInboxListError reports whether a result is an error state without requiring
// adapters to inspect a transport-specific error type.
func IsInboxListError(result InboxListResult) bool { return result.State == InboxListStateError }
