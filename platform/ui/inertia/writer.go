package inertia

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type responseWriter struct {
	body       []byte
	statusCode int
	header     http.Header
}

func write(ctx *fiber.Ctx, w *responseWriter) error {
	for key, values := range w.Header() {
		for _, value := range values {
			ctx.Set(key, value)
		}
	}

	ctx = ctx.Status(w.StatusCode())

	body := w.Body()
	if len(body) > 0 {
		return ctx.Send(body)
	}

	return nil
}

func newResponseWriter() *responseWriter {
	return &responseWriter{
		header: http.Header{},
	}
}

func (w *responseWriter) Header() http.Header {
	return w.header
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return len(b), nil
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *responseWriter) Body() []byte {
	return w.body
}

func (w *responseWriter) StatusCode() int {
	if w.statusCode == 0 {
		return http.StatusOK
	}
	return w.statusCode
}
