package inertia

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type InertiaContext struct {
	fasthttp *fasthttp.RequestCtx
	userCtx  context.Context
}

func NewInertiaContext(ctx *fiber.Ctx) context.Context {
	return &InertiaContext{
		fasthttp: ctx.Context(),
		userCtx:  ctx.UserContext(),
	}
}

func (i *InertiaContext) Deadline() (deadline time.Time, ok bool) {
	return i.fasthttp.Deadline()
}

func (i *InertiaContext) Done() <-chan struct{} {
	return i.fasthttp.Done()
}

func (i *InertiaContext) Err() error {
	return i.fasthttp.Err()
}

func (i *InertiaContext) Value(key any) any {
	if val := i.userCtx.Value(key); val != nil {
		return val
	}
	return i.fasthttp.UserValue(key)
}
