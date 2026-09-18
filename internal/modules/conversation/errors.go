package conversation

import "errors"

var (
	// ErrNotFound is the stable tenant-scoped lookup miss sentinel.
	ErrNotFound = errors.New("conversation resource not found")
	// ErrConflict indicates a duplicate or immutable public identity conflict.
	ErrConflict = errors.New("conversation resource conflict")
	// ErrInvalid indicates a domain or repository input invariant failure.
	ErrInvalid = errors.New("invalid conversation resource")

	ErrConversationNotFound           = ErrNotFound
	ErrConversationConflict           = ErrConflict
	ErrInvalidConversation            = ErrInvalid
	ErrParticipantNotFound            = ErrNotFound
	ErrParticipantConflict            = ErrConflict
	ErrInvalidParticipant             = ErrInvalid
	ErrConversationAssignmentNotFound = ErrNotFound
	ErrConversationAssignmentConflict = ErrConflict
	ErrInvalidConversationAssignment  = ErrInvalid

	// ErrInvalidInboxQuery and ErrInvalidInboxCursor are safe validation
	// sentinels for the transport-neutral inbox list boundary.
	ErrInvalidInboxQuery  = errors.New("invalid inbox query")
	ErrInvalidInboxCursor = errors.New("invalid inbox cursor")
)
