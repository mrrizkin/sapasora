package kuproy

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Templates struct {
	*template.Template
}

func NewTemplates() *Templates {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("templates.go: cannot determine caller path")
	}
	dir := filepath.Dir(currentFile)
	pattern := filepath.Join(dir, "templates", "*.gotmpl")
	title := cases.Title(language.English)
	tmpl, err := template.New("kuproy").
		Funcs(template.FuncMap{
			"lower": strings.ToLower,
			"title": title.String,
			"json": func(v string) string {
				return strcase.ToSnake(v)
			},
			"camel_case": func(v string) string {
				return strcase.ToCamel(v)
			},
			"snake_case": func(v string) string {
				return strcase.ToSnake(v)
			},
		}).
		ParseGlob(pattern)
	if err != nil {
		panic(fmt.Errorf("failed to parse templates: %w", err))
	}

	return &Templates{
		Template: tmpl,
	}
}

func (t *Templates) Write(filePath, name string, data any) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.ExecuteTemplate(f, name, data)
}
