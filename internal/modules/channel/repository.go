package channel

import "context"

// ChannelAccountFilter limits tenant-scoped list operations. Nil fields mean
// no filter. The returned order is stable by internal creation ID.
type ChannelAccountFilter struct {
	Type     *ChannelType
	Provider *Provider
	Status   *ConnectionState
}

// ChannelAccountRepository is the persistence boundary for the channel
// domain. Every lookup carries a tenant ID; public IDs are never sufficient
// to cross that boundary by themselves.
type ChannelAccountRepository interface {
	ListChannelAccounts(context.Context, string, ChannelAccountFilter) ([]*ChannelAccount, error)
	CreateChannelAccount(context.Context, *ChannelAccount) error
	GetChannelAccountByPublicID(context.Context, string, string) (*ChannelAccount, error)
	UpdateChannelAccount(context.Context, *ChannelAccount) error
	DeleteChannelAccount(context.Context, string, string) error
	RecordConnectionEvent(context.Context, *ConnectionEvent) error
	ListConnectionEvents(context.Context, string, string) ([]*ConnectionEvent, error)
}

// Repository is a short descriptive alias for the channel account repository.
type Repository = ChannelAccountRepository
type ChannelRepository = ChannelAccountRepository
type ChannelAccountListFilter = ChannelAccountFilter
