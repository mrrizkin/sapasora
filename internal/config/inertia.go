package config

type Inertia struct {
	ManifestPath   string `name:"manifest_path"   env:"INERTIA_VITE_MANIFEST_PATH,default=public/build/manifest.json"`
	ContainerID    string `name:"container_id"    env:"INERTIA_CONTAINER_ID,default=app"`
	EncryptHistory bool   `name:"encrypt_history" env:"INERTIA_ENCRYPT_HISTORY,default=true"`
	WithSSR        bool   `name:"ssr"             env:"INERTIA_WITH_SSR,default=false"`
	SSRURL         string `name:"ssr_url"         env:"INERTIA_SSR_URL,default=http://127.0.0.1:13714"`
}
