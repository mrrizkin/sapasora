package bootstrap

import (
	"sapasora/platform/ui/bundler"
)

// Bundler decorates the bundler
// @wired:decorate
func Bundler(bundler *bundler.Bundler) *bundler.Bundler {
	vite := bundler.GetVite()
	cuncurrency := 6
	vite.Prefetch(&cuncurrency, "")
	return bundler
}
