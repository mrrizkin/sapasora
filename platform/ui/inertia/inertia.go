package inertia

import (
	"fmt"
	"net/http"

	"sapasora/platform/session"
	"sapasora/platform/ui/bundler"
	"sapasora/platform/ui/view"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/romsar/gonertia"
)

type Inertia struct {
	core    *gonertia.Inertia
	view    *view.View
	session *session.Session

	options []gonertia.Option
}

func NewInertia(
	config *InertiaConfig,
	view *view.View,
	bundler *bundler.Bundler,
	session *session.Session,
) (*Inertia, error) {
	options := []gonertia.Option{
		gonertia.WithFlashProvider(newFlash()),
	}

	if config.ContainerID != "" {
		options = append(options, gonertia.WithContainerID(config.ContainerID))
	}

	if config.ManifestPath != "" {
		options = append(
			options,
			gonertia.WithVersion(getVersionFromManifest(config.ManifestPath)),
		)
	}

	if config.EncryptHistory {
		options = append(options, gonertia.WithEncryptHistory(config.EncryptHistory))
	}

	if config.WithSSR {
		options = append(options, gonertia.WithSSR(config.SSRURL))
	}

	inertia, err := gonertia.New(
		fmt.Sprintf(inertiaEntry, bundler.Entry("resources/js/app.ts")),
		options...,
	)
	if err != nil {
		return nil, err
	}

	return &Inertia{
		core:    inertia,
		view:    view,
		session: session,
		options: options,
	}, nil
}

func (i *Inertia) SetEntry(entry []byte) {
	inertia, err := gonertia.NewFromBytes(entry, i.options...)
	if err != nil {
		panic(err)
	}
	i.core = inertia
}

func (i *Inertia) Middleware(c *fiber.Ctx) error {
	if session, err := i.session.Get(c); err == nil {
		c.Locals("session_id", session.ID())
	}
	return adaptor.HTTPMiddleware(i.core.Middleware)(c)
}

func (i *Inertia) Redirect(ctx *fiber.Ctx, url string, status ...int) error {
	r, err := i.convertRequest(ctx)
	if err != nil {
		return err
	}

	w := newResponseWriter()
	i.core.Redirect(w, r, url, status...)
	return write(ctx, w)
}

func (i *Inertia) Location(ctx *fiber.Ctx, url string, status ...int) error {
	r, err := i.convertRequest(ctx)
	if err != nil {
		return err
	}

	w := newResponseWriter()
	i.core.Location(w, r, url, status...)
	return write(ctx, w)
}

func (i *Inertia) Back(ctx *fiber.Ctx, status ...int) error {
	r, err := i.convertRequest(ctx)
	if err != nil {
		return err
	}

	w := newResponseWriter()
	i.core.Back(w, r, status...)
	return write(ctx, w)
}

func (i *Inertia) Render(ctx *fiber.Ctx, component string, props ...fiber.Map) error {
	r, err := i.convertRequest(ctx)
	if err != nil {
		return err
	}

	inertiaProps := make([]gonertia.Props, len(props))
	for i, prop := range props {
		inertiaProps[i] = gonertia.Props(prop)
	}

	w := newResponseWriter()
	if err := i.core.Render(w, r, component, inertiaProps...); err != nil {
		return err
	}

	ctx.Set("Content-Type", "text/html")
	return write(ctx, w)
}

func (i *Inertia) ShareProp(key string, value any) {
	i.core.ShareProp(key, value)
}

func (i *Inertia) ShareTemplateFunc(key string, fn any) error {
	return i.core.ShareTemplateFunc(key, fn)
}

func (i *Inertia) ShareTemplateData(key string, value any) {
	i.core.ShareTemplateData(key, value)
}

func (i *Inertia) convertRequest(ctx *fiber.Ctx) (*http.Request, error) {
	r, err := adaptor.ConvertRequest(ctx, true)
	if err != nil {
		return nil, err
	}
	r = r.WithContext(i.view.WithContext(ctx.Context(), ctx.UserContext()))
	return r, nil
}
