// Package policies provides the policies for the application
package policies

import (
	"sapasora/internal/modules/permission"
	"sapasora/internal/modules/role"
	"sapasora/platform/satpam"
)

type CanGetAvatarGateway struct{}

func (p *CanGetAvatarGateway) Name() string { return "can_get_avatar_gateway" }
func (p *CanGetAvatarGateway) Authorize(sub *satpam.Attribute, res map[string]*satpam.Attribute) {
	permission, ok := sub.GetAttribute("permission").(*permission.Permission)
	if !ok {
		if role, ok := sub.GetAttribute("role").(*role.Role); ok {
			permission = role.Permissions
		}
	}

	satpam.Assert(permission != nil, "you does not have any permission")
	satpam.Assert(
		permission.CanGetAvatarGateway,
		"you does not have permission to get avatar gateway",
	)
}
