package account

import (
	"sapasora/internal/modules/account"
)

type AccountResponse *account.Account // @name AccountController.AccountResponse

type AccountListResponse *account.Pagination[*account.Account] // @name AccountController.AccountListResponse

type AccountLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
} // @name AccountController.AccountLoginRequest

type AccountStoreRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	RoleID   string `json:"role_id"`
} // @name AccountController.AccountStoreRequest

type AccountUpdateRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	RoleID   string `json:"role_id"`
} // @name AccountController.AccountUpdateRequest

type AccountUpdatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
} // @name AccountController.AccountUpdatePasswordRequest
