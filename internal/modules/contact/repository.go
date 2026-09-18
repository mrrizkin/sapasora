package contact

import "context"

// ContactFilter limits tenant-scoped contact listing.
type ContactFilter struct {
	Status *ContactStatus
	Source *ContactSource
}

// ContactRepository is the persistence boundary for the contact domain.
// Every operation carries a tenant ID directly or through a tenant-owned
// model; public IDs and normalized identities are never global lookup grants.
type ContactRepository interface {
	CreateContact(context.Context, *Contact) error
	GetContactByPublicID(context.Context, string, string) (*Contact, error)
	GetContactByID(context.Context, string, uint64) (*Contact, error)
	ListContacts(context.Context, string, ContactFilter) ([]*Contact, error)
	UpdateContact(context.Context, *Contact) error

	CreateContactAddress(context.Context, *ContactAddress) error
	GetContactAddressByPublicID(context.Context, string, string) (*ContactAddress, error)
	GetContactAddressByID(context.Context, string, uint64) (*ContactAddress, error)
	FindContactAddressByIdentity(context.Context, string, AddressIdentity) (*ContactAddress, error)
	ListContactAddresses(context.Context, string, uint64) ([]*ContactAddress, error)
	UpdateContactAddress(context.Context, *ContactAddress) error
}

// Repository is a concise compatibility alias for the domain boundary.
type Repository = ContactRepository
