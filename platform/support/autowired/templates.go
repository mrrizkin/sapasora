package autowired

import (
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"text/template"
)

type Templates struct {
	tmpl *template.Template
}

func NewTemplates() *Templates {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("templates.go: cannot determine caller path")
	}
	dir := filepath.Dir(currentFile)
	pattern := filepath.Join(dir, "templates", "*.gotmpl")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		panic(fmt.Errorf("failed to parse templates: %w", err))
	}

	return &Templates{
		tmpl: tmpl,
	}
}

func (t *Templates) Render(writer io.Writer, name string, data any) error {
	return t.tmpl.ExecuteTemplate(writer, name, data)
}
