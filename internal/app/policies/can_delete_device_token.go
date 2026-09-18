// Package policies provides the policies for the application
package policies

import (
	"sapasora/internal/modules/permission"
	"sapasora/internal/modules/role"
	"sapasora/platform/satpam"
)

type CanDeleteDeviceToken struct{}

func (p *CanDeleteDeviceToken) Name() string { return "can_delete_device_token" }
func (p *CanDeleteDeviceToken) Authorize(sub *satpam.Attribute, res map[string]*satpam.Attribute) {
	permission, ok := sub.GetAttribute("permission").(*permission.Permission)
	if !ok {
		if role, ok := sub.GetAttribute("role").(*role.Role); ok {
			permission = role.Permissions
		}
	}

	satpam.Assert(permission != nil, "you does not have any permission")
	satpam.Assert(
		permission.CanDeleteDeviceToken,
		"you does not have permission to delete device token",
	)
}
