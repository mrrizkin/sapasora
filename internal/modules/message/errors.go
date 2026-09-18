package message

import "errors"

var (
	// ErrNotFound is returned for a tenant-scoped lookup miss.
	ErrNotFound = errors.New("message resource not found")
	// ErrConflict indicates a duplicate or immutable identity conflict.
	ErrConflict = errors.New("message resource conflict")
	// ErrInvalid indicates a message domain invariant failure.
	ErrInvalid = errors.New("invalid message resource")

	ErrMessageNotFound        = ErrNotFound
	ErrMessageConflict        = ErrConflict
	ErrInvalidMessage         = ErrInvalid
	ErrAttachmentNotFound     = ErrNotFound
	ErrAttachmentConflict     = ErrConflict
	ErrInvalidAttachment      = ErrInvalid
	ErrDeliveryRecordNotFound = ErrNotFound
	ErrDeliveryRecordConflict = ErrConflict
	ErrInvalidDeliveryRecord  = ErrInvalid
	ErrMessageEventNotFound   = ErrNotFound
	ErrMessageEventConflict   = ErrConflict
	ErrInvalidMessageEvent    = ErrInvalid
)
