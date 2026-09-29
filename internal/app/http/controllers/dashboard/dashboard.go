// Package dashboard provides the Dashboard controllers
package dashboard

import (
	"sapasora/internal/app/http/controllers"
	"sapasora/internal/modules/apikey"
	"sapasora/internal/modules/device"

	"github.com/gofiber/fiber/v2"
)

type DashboardController struct {
	*controllers.Controller

	deviceService device.DeviceService
	apikeyService apikey.APIKeyService
}

// NewDashboardController creates a new Dashboardcontrollers
// @wired:provide
func NewDashboardController(
	controller *controllers.Controller,
	deviceService device.DeviceService,
	apikeyService apikey.APIKeyService,
) *DashboardController {
	return &DashboardController{
		Controller:    controller,
		deviceService: deviceService,
		apikeyService: apikeyService,
	}
}

func (c *DashboardController) Index(ctx *fiber.Ctx) error {
	ownerID, err := c.GetOwnerID(ctx, "account")
	if err != nil {
		return err
	}

	// Small dashboards only: pull a generous page of the owner's own devices
	// and tally status counts in-memory instead of adding dedicated count
	// queries. This keeps the summary owner-scoped, matching every other
	// resource in the app.
	deviceList, err := c.deviceService.ListDeviceForUser(ctx.Context(), ownerID, "", 1, 1000)
	if err != nil {
		return err
	}

	var connected, disconnected, inactive int
	for _, d := range deviceList.Data {
		switch d.Status {
		case device.DeviceStatusConnected:
			connected++
		case device.DeviceStatusDisconnected:
			disconnected++
		default:
			inactive++
		}
	}

	apikeyList, err := c.apikeyService.ListAPIKeyForUser(ctx.Context(), ownerID, 1, 1000)
	if err != nil {
		return err
	}

	return c.Inertia(ctx, "dashboard/index", fiber.Map{
		"stats": fiber.Map{
			"total_devices":        len(deviceList.Data),
			"connected_devices":    connected,
			"disconnected_devices": disconnected,
			"inactive_devices":     inactive,
			"total_api_keys":       len(apikeyList.Data),
		},
	})
}
