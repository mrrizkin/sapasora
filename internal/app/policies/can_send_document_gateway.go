// Package policies provides the policies for the application
package policies

import (
	"sapasora/internal/modules/permission"
	"sapasora/internal/modules/role"
	"sapasora/platform/satpam"
)

type CanSendDocumentGateway struct{}

func (p *CanSendDocumentGateway) Name() string { return "can_send_document_gateway" }

func (p *CanSendDocumentGateway) Authorize(
	sub *satpam.Attribute,
	res map[string]*satpam.Attribute,
) {
	permission, ok := sub.GetAttribute("permission").(*permission.Permission)
	if !ok {
		if role, ok := sub.GetAttribute("role").(*role.Role); ok {
			permission = role.Permissions
		}
	}

	satpam.Assert(permission != nil, "you does not have any permission")
	satpam.Assert(
		permission.CanSendDocumentGateway,
		"you does not have permission to send document gateway",
	)
}
