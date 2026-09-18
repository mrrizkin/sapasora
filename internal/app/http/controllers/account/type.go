package account

import (
	"sapasora/internal/modules/account"
)

type AccountResponse *account.Account // @name AccountController.AccountResponse

type AccountListResponse *account.Pagination[*account.Account] // @name AccountController.AccountListResponse

type AccountLoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=3,max=255"`
} // @name AccountController.AccountLoginRequest

type AccountStoreRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Username string `json:"username" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=3,max=255"`
	RoleID   string `json:"role_id" validate:"required,min=1,max=255"`
} // @name AccountController.AccountStoreRequest

type AccountUpdateRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Username string `json:"username" validate:"required,min=3,max=255"`
	RoleID   string `json:"role_id" validate:"required,min=1,max=255"`
} // @name AccountController.AccountUpdateRequest

type AccountUpdatePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=3,max=255"`
	NewPassword string `json:"new_password" validate:"required,min=3,max=255"`
} // @name AccountController.AccountUpdatePasswordRequest
