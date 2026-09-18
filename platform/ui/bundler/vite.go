package bundler

import "context"

func Vite(ctx context.Context) *Bundler {
	return ctx.Value("bundler").(*Bundler)
}
