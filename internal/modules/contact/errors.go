package contact

import "errors"

var (
	// ErrNotFound is the stable tenant-scoped lookup miss sentinel.
	ErrNotFound = errors.New("contact resource not found")
	// ErrConflict indicates a duplicate public ID or normalized address
	// identity, or an immutable identity change.
	ErrConflict = errors.New("contact resource conflict")
	// ErrInvalid indicates a model or repository input invariant failure.
	ErrInvalid = errors.New("invalid contact resource")
	// ErrMergeConfirmationRequired prevents a preview from being treated as an
	// implicit authorization to mutate contacts.
	ErrMergeConfirmationRequired = errors.New("explicit merge confirmation is required")

	ErrContactNotFound                  = ErrNotFound
	ErrContactAddressNotFound           = ErrNotFound
	ErrContactConflict                  = ErrConflict
	ErrContactAddressConflict           = ErrConflict
	ErrInvalidContact                   = ErrInvalid
	ErrInvalidContactAddress            = ErrInvalid
	ErrContactMergeConfirmationRequired = ErrMergeConfirmationRequired
)
