package inertia

import "sapasora/platform/config"

type InertiaConfig struct {
	ManifestPath   string
	ContainerID    string
	EncryptHistory bool
	WithSSR        bool
	SSRURL         string
}

func NewInertiaConfig(cfg config.Config) *InertiaConfig {
	return &InertiaConfig{
		ManifestPath:   cfg.GetString("inertia.manifest_path"),
		ContainerID:    cfg.GetString("inertia.container_id"),
		EncryptHistory: cfg.GetBool("inertia.encrypt_history"),
		WithSSR:        cfg.GetBool("inertia.ssr"),
		SSRURL:         cfg.GetString("inertia.ssr_url"),
	}
}
