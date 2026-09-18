package config

// MustLoad is a helper that panics on config load errors (useful for required configs)
func MustLoad[T any](config Config, name string, target T) T {
	if err := config.LoadStruct(name, target); err != nil {
		panic(err)
	}
	return target
}

// LoadOrDefault loads config or returns default instance on error
func LoadOrDefault[T any](config Config, name string, target T, defaultVal T) T {
	if err := config.LoadStruct(name, target); err != nil {
		return defaultVal
	}
	return target
}
