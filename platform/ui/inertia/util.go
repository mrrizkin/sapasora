package inertia

import (
	"crypto/md5"
	"encoding/hex"
	"os"

	"sapasora/platform/support/hash"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/romsar/gonertia"
)

func EncryptHistory(c *fiber.Ctx, encrypt ...bool) {
	c.SetUserContext(gonertia.SetEncryptHistory(c.UserContext(), encrypt...))
}

func ClearHistory(c *fiber.Ctx) {
	c.SetUserContext(gonertia.ClearHistory(c.UserContext()))
}

func EncryptHistoryMiddleware(c *fiber.Ctx) error {
	EncryptHistory(c)
	return c.Next()
}

func IsInertiaRequest(c *fiber.Ctx) bool {
	r, err := adaptor.ConvertRequest(c, true)
	if err != nil {
		return false
	}

	return gonertia.IsInertiaRequest(r)
}

func AddValidationError(c *fiber.Ctx, errMap fiber.Map) error {
	c.SetUserContext(gonertia.AddValidationErrors(
		c.UserContext(),
		gonertia.ValidationErrors(errMap),
	))
	return nil
}

func SharePropContext(ctx *fiber.Ctx, props fiber.Map) {
	ctx.SetUserContext(gonertia.SetProps(ctx.UserContext(), gonertia.Props(props)))
}

func getVersionFromManifest(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return hash.NanoID()
	}

	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
