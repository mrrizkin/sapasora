package autowired

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (aw *AutoWired) GenerateCode() error {
	sort.Slice(aw.constructors, func(i, j int) bool {
		return aw.constructors[i].PackageName+"."+aw.constructors[i].Name <
			aw.constructors[j].PackageName+"."+aw.constructors[j].Name
	})

	if _, err := os.Stat(filepath.Dir(aw.config.OutputFile)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(aw.config.OutputFile), os.ModePerm); err != nil {
			return err
		}
	}

	file, err := os.Create(aw.config.OutputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	data := aw.buildAdvancedTemplateData()

	return aw.templates.Render(file, "registration.gotmpl", data)
}

func (aw *AutoWired) buildAdvancedTemplateData() map[string]any {
	providers := make([]string, 0)
	providerFuncs := make([]string, 0)
	decorators := make([]string, 0)
	decoratorFuncs := make([]string, 0)
	invokers := make([]string, 0)
	invokerFuncs := make([]string, 0)
	infoMap := make([]map[string]string, 0)
	groups := make(map[string][]string)
	namedProviders := make(map[string]string)

	for _, constructor := range aw.constructors {
		fullName := constructor.Name
		if constructor.PackageName != aw.config.PackageName {
			fullName = constructor.PackageName + "." + constructor.Name
		}

		// Build the fx option string
		var optionBuilder strings.Builder

		switch constructor.Annotation.Type {
		case "provide":

			var asAnotation string
			var resultTagsAnotation string

			// Add interfaces
			if len(constructor.Annotation.As) > 0 {
				for _, iface := range constructor.Annotation.As {
					if strings.Contains(iface, ".") {
						asAnotation = fmt.Sprintf("fx.As(new(%s))", iface)
					} else if constructor.PackageName != aw.config.PackageName {
						asAnotation = fmt.Sprintf("fx.As(new(%s))", constructor.PackageName+"."+iface)
					} else {
						asAnotation = fmt.Sprintf("fx.As(new(%s))", iface)
					}
				}
			}

			// Add name
			if constructor.Annotation.Name != "" {
				resultTagsAnotation = fmt.Sprintf(
					"fx.ResultTags(`name:\"%s\"`)",
					constructor.Annotation.Name,
				)
			}

			// Add group
			if constructor.Annotation.Group != "" {
				resultTagsAnotation = fmt.Sprintf(
					"fx.ResultTags(`group:\"%s\"`)",
					constructor.Annotation.Group,
				)

				groupKey := cases.Title(language.English).String(constructor.Annotation.Group)
				if groups[groupKey] == nil {
					groups[groupKey] = make([]string, 0)
				}
				groups[groupKey] = append(
					groups[groupKey],
					fmt.Sprintf("fx.Provide(fx.Annotate(%s, %s)),", fullName, resultTagsAnotation),
				)
			}

			if asAnotation != "" {
				fmt.Fprintf(
					&optionBuilder,
					"fx.Provide(fx.Annotate(%s, %s)),",
					fullName,
					asAnotation,
				)
			}

			if resultTagsAnotation != "" {
				fmt.Fprintf(
					&optionBuilder,
					"fx.Provide(fx.Annotate(%s, %s)),",
					fullName,
					resultTagsAnotation,
				)
			}

			if asAnotation == "" && resultTagsAnotation == "" {
				fmt.Fprintf(&optionBuilder, "fx.Provide(%s),", fullName)
			}
			providers = append(providers, optionBuilder.String())
			providerFuncs = append(providerFuncs, fullName)

			// Named provider
			if constructor.Annotation.Name != "" {
				namedProviders[cases.Title(language.English).String(constructor.Annotation.Name)] = fullName
			}

		case "decorate":
			var asAnotation string

			if len(constructor.Annotation.As) > 0 {
				for _, iface := range constructor.Annotation.As {
					if strings.Contains(iface, ".") {
						asAnotation = fmt.Sprintf("fx.As(new(%s))", iface)
					} else if constructor.PackageName != aw.config.PackageName {
						asAnotation = fmt.Sprintf("fx.As(new(%s))", constructor.PackageName+"."+iface)
					} else {
						asAnotation = fmt.Sprintf("fx.As(new(%s))", iface)
					}
				}
			}

			if asAnotation != "" {
				fmt.Fprintf(
					&optionBuilder,
					"fx.Decorate(fx.Annotate(%s, %s)),",
					fullName,
					asAnotation,
				)
			} else {
				fmt.Fprintf(&optionBuilder, "fx.Decorate(%s),", fullName)
			}

			decorators = append(decorators, optionBuilder.String())
			decoratorFuncs = append(decoratorFuncs, fullName)

		case "invoke":
			invokers = append(invokers, "fx.Invoke("+fullName+"),")
			invokerFuncs = append(invokerFuncs, fullName)
		}

		// Build info map for debugging
		info := fmt.Sprintf("%s:%s", constructor.FullPath, constructor.Annotation.Type)
		if runtime.GOOS == "windows" {
			info = strings.ReplaceAll(info, "\\", "/")
		}

		if constructor.Comment != "" {
			info += " - " + constructor.Comment
		}
		infoMap = append(infoMap, map[string]string{
			"Key":   fullName,
			"Value": info,
		})
	}

	return map[string]any{
		"PackageName":    aw.config.PackageName,
		"Imports":        aw.imports,
		"Providers":      providers,
		"ProviderFuncs":  providerFuncs,
		"Decorators":     decorators,
		"DecoratorFuncs": decoratorFuncs,
		"Invokers":       invokers,
		"InvokerFuncs":   invokerFuncs,
		"Groups":         groups,
		"NamedProviders": namedProviders,
		"InfoMap":        infoMap,
		"SourceCount":    len(aw.constructors),
		"Comment":        fmt.Sprintf("Auto-generated from %d constructors", len(aw.constructors)),
	}
}
