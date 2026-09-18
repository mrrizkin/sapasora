package bootstrap

import (
	"sapasora/platform/ui/inertia"
	"sapasora/platform/ui/view"
	"sapasora/resources/views"
)

// Inertia is the middleware for the inertia
// @wired:decorate
func Inertia(v *view.View, inertia *inertia.Inertia) *inertia.Inertia {
	content := v.RenderBytes(views.App())
	inertia.SetEntry(content)
	return inertia
}
