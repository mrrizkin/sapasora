// Package security provides server security utilities
package security

import (
	"fmt"
	"strings"
)

// CSPBuilder helps construct Content Security Policy strings
type CSPBuilder struct {
	directives map[string][]string
}

// NewCSPBuilder creates a new CSP builder
func NewCSPBuilder() *CSPBuilder {
	return &CSPBuilder{
		directives: make(map[string][]string),
	}
}

// DefaultSrc sets the default-src directive
func (b *CSPBuilder) DefaultSrc(sources ...string) *CSPBuilder {
	b.directives["default-src"] = sources
	return b
}

// ScriptSrc sets the script-src directive
func (b *CSPBuilder) ScriptSrc(sources ...string) *CSPBuilder {
	b.directives["script-src"] = sources
	return b
}

// StyleSrc sets the style-src directive
func (b *CSPBuilder) StyleSrc(sources ...string) *CSPBuilder {
	b.directives["style-src"] = sources
	return b
}

// FontSrc sets the font-src directive
func (b *CSPBuilder) FontSrc(sources ...string) *CSPBuilder {
	b.directives["font-src"] = sources
	return b
}

// ImgSrc sets the img-src directive
func (b *CSPBuilder) ImgSrc(sources ...string) *CSPBuilder {
	b.directives["img-src"] = sources
	return b
}

// ConnectSrc sets the connect-src directive
func (b *CSPBuilder) ConnectSrc(sources ...string) *CSPBuilder {
	b.directives["connect-src"] = sources
	return b
}

// FrameSrc sets the frame-src directive
func (b *CSPBuilder) FrameSrc(sources ...string) *CSPBuilder {
	b.directives["frame-src"] = sources
	return b
}

// ObjectSrc sets the object-src directive
func (b *CSPBuilder) ObjectSrc(sources ...string) *CSPBuilder {
	b.directives["object-src"] = sources
	return b
}

// MediaSrc sets the media-src directive
func (b *CSPBuilder) MediaSrc(sources ...string) *CSPBuilder {
	b.directives["media-src"] = sources
	return b
}

// BaseURI sets the base-uri directive
func (b *CSPBuilder) BaseURI(sources ...string) *CSPBuilder {
	b.directives["base-uri"] = sources
	return b
}

// FormAction sets the form-action directive
func (b *CSPBuilder) FormAction(sources ...string) *CSPBuilder {
	b.directives["form-action"] = sources
	return b
}

// FrameAncestors sets the frame-ancestors directive
func (b *CSPBuilder) FrameAncestors(sources ...string) *CSPBuilder {
	b.directives["frame-ancestors"] = sources
	return b
}

// AddSource adds a source to an existing directive
func (b *CSPBuilder) AddSource(directive string, sources ...string) *CSPBuilder {
	if existing, ok := b.directives[directive]; ok {
		b.directives[directive] = append(existing, sources...)
	} else {
		b.directives[directive] = sources
	}
	return b
}

// Build constructs the CSP string
func (b *CSPBuilder) Build() string {
	if len(b.directives) == 0 {
		return ""
	}

	var parts []string
	// Order matters for readability, so we'll build in a logical order
	order := []string{
		"default-src",
		"script-src",
		"style-src",
		"img-src",
		"font-src",
		"connect-src",
		"media-src",
		"object-src",
		"frame-src",
		"base-uri",
		"form-action",
		"frame-ancestors",
	}

	for _, directive := range order {
		if sources, ok := b.directives[directive]; ok && len(sources) > 0 {
			parts = append(parts, fmt.Sprintf("%s %s", directive, strings.Join(sources, " ")))
			delete(b.directives, directive)
		}
	}

	// Add any remaining directives not in the ordered list
	for directive, sources := range b.directives {
		if len(sources) > 0 {
			parts = append(parts, fmt.Sprintf("%s %s", directive, strings.Join(sources, " ")))
		}
	}

	return strings.Join(parts, "; ")
}
