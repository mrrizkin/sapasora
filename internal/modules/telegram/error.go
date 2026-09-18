package telegram

import "errors"

var (
	ErrDeviceNotFound         = errors.New("device not found")
	ErrNoAvatarFound          = errors.New("no avatar found")
	ErrNotConnected           = errors.New("not connected")
	ErrAlreadyLoggedIn        = errors.New("already logged in")
	ErrNotLoggedIn            = errors.New("not logged in")
	ErrNoSession              = errors.New("no session")
	ErrClientAlreadyConnected = errors.New("client already connected")
	ErrFailedToConnect        = errors.New("failed to connect")
	ErrInvalidPhoneNumber     = errors.New("invalid phone number")
	ErrEmptyBody              = errors.New("body cannot be empty")
	ErrMissingStanzaID        = errors.New("missing stanza id in contextinfo")
	ErrMissingParticipant     = errors.New("missing participant in contextinfo")
)
