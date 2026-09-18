// Package policies provides the policies for the application
package policies

import (
	"sapasora/internal/modules/permission"
	"sapasora/internal/modules/role"
	"sapasora/platform/satpam"
)

type CanGetQRGateway struct{}

func (p *CanGetQRGateway) Name() string { return "can_get_qr_gateway" }
func (p *CanGetQRGateway) Authorize(sub *satpam.Attribute, res map[string]*satpam.Attribute) {
	permission, ok := sub.GetAttribute("permission").(*permission.Permission)
	if !ok {
		if role, ok := sub.GetAttribute("role").(*role.Role); ok {
			permission = role.Permissions
		}
	}

	satpam.Assert(permission != nil, "you does not have any permission")
	satpam.Assert(
		permission.CanGetQRGateway,
		"you does not have permission to get QR gateway",
	)
}
