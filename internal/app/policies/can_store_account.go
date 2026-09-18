// Package policies provides the policies for the application
package policies

import (
	"sapasora/internal/modules/permission"
	"sapasora/internal/modules/role"
	"sapasora/platform/satpam"
)

type CanStoreAccount struct{}

func (p *CanStoreAccount) Name() string { return "can_store_account" }
func (p *CanStoreAccount) Authorize(sub *satpam.Attribute, res map[string]*satpam.Attribute) {
	permission, ok := sub.GetAttribute("permission").(*permission.Permission)
	if !ok {
		if role, ok := sub.GetAttribute("role").(*role.Role); ok {
			permission = role.Permissions
		}
	}

	satpam.Assert(permission != nil, "you does not have any permission")
	satpam.Assert(
		permission.CanStoreAccount,
		"you does not have permission to store account",
	)
}
