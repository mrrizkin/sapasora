package role

import (
	"sapasora/internal/modules/permission"
	"sapasora/internal/modules/role"
)

type RoleResponse *role.Role // @name RoleController.RoleResponse

type RoleListResponse *role.Pagination[*role.Role] // @name RoleController.RoleListResponse

type RoleStoreRequest struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Permissions permission.Permission `json:"permissions"`
} // @name RoleController.RoleStoreRequest

type RoleUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
} // @name RoleController.RoleUpdateRequest

type RoleUpdatePermissionsRequest struct {
	Permissions permission.Permission `json:"permissions"`
} // @name RoleController.RoleUpdatePermissionsRequest
