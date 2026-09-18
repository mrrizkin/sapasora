// Package role provides the Role controllers
package role

import (
	"sapasora/internal/app/http/controllers"
	"sapasora/internal/app/policies"
	"sapasora/internal/modules/role"
	"sapasora/platform/satpam"
	"sapasora/platform/support/hash"

	"github.com/gofiber/fiber/v2"
)

type RoleController struct {
	*controllers.Controller

	roleService role.RoleService
}

// NewRoleController creates a new Rolecontrollers
// @wired:provide
func NewRoleController(
	controller *controllers.Controller,
	roleService role.RoleService,
) *RoleController {
	return &RoleController{
		Controller:  controller,
		roleService: roleService,
	}
}

// List godoc
// @Summary      List role
// @Description  List role
// @Tags         Role
// @Produce      json
// @Param        page    query  int    false "Page"
// @Param        limit   query  int    false "Limit"
// @Param        search  query  string false "Search"
// @Success      200 {object} RoleListResponse
// @Security     Authorization
// @Router       /api/v1/role [get]
func (c *RoleController) List(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanListRole{})
	subject, err := c.GetSubject(ctx, "account")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrForbidden
	}

	gate.AuthorizeAllPermissions(subject)

	params, err := c.ParseListQuery(ctx)
	if err != nil {
		return err
	}
	page := params.Page
	limit := params.Limit
	search := params.Search

	roleList, err := c.roleService.ListRole(ctx.Context(), search, page, limit)
	if err != nil {
		return err
	}

	return ctx.JSON(RoleListResponse(roleList))
}

// Get godoc
// @Summary      Get role
// @Description  Get role
// @Tags         Role
// @Produce      json
// @Param        id  path  string true "ID"
// @Success      200 {object} RoleResponse
// @Security     Authorization
// @Router       /api/v1/role/{id} [get]
func (c *RoleController) Get(ctx *fiber.Ctx) error {
	id, err := c.PublicIDParam(ctx)
	if err != nil {
		return err
	}
	role, err := c.roleService.GetRoleByPublicID(ctx.Context(), id)
	if err != nil {
		return err
	}

	gate := satpam.New(&policies.CanGetRole{}).AddResource("role", role)
	subject, err := c.GetSubject(ctx, "account")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}

	gate.AuthorizeAllPermissions(subject)

	return ctx.JSON(RoleResponse(role))
}

// Store godoc
// @Summary      Store role
// @Description  Store role
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        role  body  RoleStoreRequest true "Body"
// @Success      200 {object} RoleResponse
// @Security     Authorization
// @Router       /api/v1/role [post]
func (c *RoleController) Store(ctx *fiber.Ctx) error {
	gate := satpam.New(&policies.CanStoreRole{})
	subject, err := c.GetSubject(ctx, "account")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}
	gate.AuthorizeAllPermissions(subject)

	var payload RoleStoreRequest
	if err := c.BodyParserValidate(ctx, &payload); err != nil {
		return err
	}

	role := role.Role{
		PublicID:    hash.NanoID(),
		Name:        payload.Name,
		Description: payload.Description,
		Permissions: &payload.Permissions,
	}

	if err := c.roleService.CreateRole(ctx.Context(), &role); err != nil {
		return err
	}

	return ctx.JSON(RoleResponse(&role))
}

// Update godoc
// @Summary      Update role
// @Description  Update role
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        id  path  string true "ID"
// @Param        role  body  RoleUpdateRequest true "Body"
// @Success      200 {object} RoleResponse
// @Security     Authorization
// @Router       /api/v1/role/{id} [put]
func (c *RoleController) Update(ctx *fiber.Ctx) error {
	id, err := c.PublicIDParam(ctx)
	if err != nil {
		return err
	}
	var payload RoleUpdateRequest
	if err := c.BodyParserValidate(ctx, &payload); err != nil {
		return err
	}

	role, err := c.roleService.GetRoleByPublicID(ctx.Context(), id)
	if err != nil {
		return err
	}

	gate := satpam.New(&policies.CanUpdateRole{}).AddResource("role", role)
	subject, err := c.GetSubject(ctx, "account")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}

	gate.AuthorizeAllPermissions(subject)

	role.Name = payload.Name
	role.Description = payload.Description

	if err := c.roleService.UpdateRole(ctx.Context(), role); err != nil {
		return err
	}

	return ctx.JSON(RoleResponse(role))
}

// UpdatePermissions godoc
// @Summary      Update role permissions
// @Description  Update role permissions
// @Tags         Role
// @Produce      json
// @Param        id  path  string true "ID"
// @Param        role body RoleUpdatePermissionsRequest true "Body"
// @Success      200 {object} RoleResponse
// @Security     Authorization
// @Router       /api/v1/role/{id}/permissions [put]
func (c *RoleController) UpdatePermissions(ctx *fiber.Ctx) error {
	id, err := c.PublicIDParam(ctx)
	if err != nil {
		return err
	}

	var payload RoleUpdatePermissionsRequest
	if err := c.BodyParserValidate(ctx, &payload); err != nil {
		return err
	}

	role, err := c.roleService.GetRoleByPublicID(ctx.Context(), id)
	if err != nil {
		return err
	}

	gate := satpam.New(&policies.CanUpdatePermissionRole{}).AddResource("role", role)
	subject, err := c.GetSubject(ctx, "account")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}

	gate.AuthorizeAllPermissions(subject)

	role.Permissions = &payload.Permissions

	if err := c.roleService.UpdateRole(ctx.Context(), role); err != nil {
		return err
	}

	return ctx.JSON(RoleResponse(role))
}

// Destroy godoc
// @Summary      Destroy role
// @Description  Destroy role
// @Tags         Role
// @Produce      json
// @Param        id  path  string true "ID"
// @Success      200 {object} RoleResponse
// @Security     Authorization
// @Router       /api/v1/role/{id} [delete]
func (c *RoleController) Destroy(ctx *fiber.Ctx) error {
	id, err := c.PublicIDParam(ctx)
	if err != nil {
		return err
	}
	role, err := c.roleService.GetRoleByPublicID(ctx.Context(), id)
	if err != nil {
		return err
	}

	gate := satpam.New(&policies.CanDeleteRole{}).AddResource("role", role)
	subject, err := c.GetSubject(ctx, "account")
	if err != nil {
		return err
	}
	if subject == nil {
		return fiber.ErrUnauthorized
	}

	gate.AuthorizeAllPermissions(subject)

	if err := c.roleService.DeleteRole(ctx.Context(), role); err != nil {
		return err
	}

	return ctx.JSON(RoleResponse(role))
}
