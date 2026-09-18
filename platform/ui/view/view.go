package view

import (
	"bytes"
	"context"
	"io"

	"sapasora/platform/ui/bundler"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type View struct {
	viewCtx *ViewCtx
}

type ViewIn struct {
	fx.In

	Bundler *bundler.Bundler
}

func NewView(in ViewIn) *View {
	ctx := NewViewCtx()
	ctx.SetUserValue("bundler", in.Bundler)

	return &View{
		viewCtx: ctx,
	}
}

func (v *View) Context() context.Context {
	return v.viewCtx
}

func (v *View) SetSharedViewData(key any, value any) {
	v.viewCtx.SetUserValue(key, value)
}

func (v *View) SetViewData(key any, value any) *View {
	viewCtx := v.viewCtx.Clone()
	viewCtx.SetUserValue(key, value)
	newView := &View{
		viewCtx: viewCtx,
	}
	return newView
}

func (v *View) Render(ctx *fiber.Ctx, component templ.Component) error {
	ctx.Set("Content-Type", "text/html")
	return v.renderTempl(v.injectContext(ctx.Context()), component, ctx.Response().BodyWriter())
}

func (v *View) RenderBytes(component templ.Component) []byte {
	buf := new(bytes.Buffer)
	v.renderTempl(v.Context(), component, buf)
	return buf.Bytes()
}

func (v *View) WithContext(ctxs ...context.Context) context.Context {
	return v.viewCtx.WithContext(ctxs...)
}

func (v *View) renderTempl(ctx context.Context, component templ.Component, writer io.Writer) error {
	return component.Render(ctx, writer)
}

func (v *View) injectContext(ctxs ...context.Context) context.Context {
	return v.viewCtx.WithContext(ctxs...)
}
