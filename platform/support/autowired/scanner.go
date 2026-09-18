package autowired

import (
	"bufio"
	"sapasora/platform/support/console"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// Updated method to generate unique aliases
func (aw *AutoWired) addImport(packageName, packagePath string) string {
	// Skip if it's the same package as the generated code
	if packageName == aw.config.PackageName || packagePath == "" {
		return packageName
	}

	// Check if we already have this import path
	if existingAlias, exists := aw.importPaths[packagePath]; exists {
		return existingAlias
	}

	// Generate unique alias
	alias := packageName
	if aw.aliasCount[packageName] > 0 {
		alias = fmt.Sprintf("%s%d", packageName, aw.aliasCount[packageName])
	}

	// Ensure the alias is truly unique (in case of conflicts)
	originalAlias := alias
	counter := aw.aliasCount[packageName]
	for {
		if _, exists := aw.imports[alias]; !exists {
			break
		}
		counter++
		alias = fmt.Sprintf("%s%d", originalAlias, counter)
	}

	// Store the mappings
	aw.imports[alias] = packagePath
	if runtime.GOOS == "windows" {
		aw.imports[alias] = strings.ReplaceAll(packagePath, "\\", "/")
	}

	aw.importPaths[packagePath] = alias
	aw.aliasCount[packageName] = counter + 1

	return alias
}

func (aw *AutoWired) ScanDirectories() error {
	for _, dir := range aw.config.Directories {
		if err := aw.scanDirectory(dir); err != nil {
			return fmt.Errorf("failed to scan directory %s: %w", dir, err)
		}
	}
	return nil
}

func (aw *AutoWired) scanDirectory(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".go") &&
			!strings.HasSuffix(path, "_test.go") &&
			!strings.HasSuffix(path, "_gen.go") {
			if path == aw.config.OutputFile {
				return nil
			}

			return aw.scanFile(path)
		}

		return nil
	})
}

func (aw *AutoWired) scanFile(path string) error {
	// Check exclusion patterns
	if aw.config.ExcludePattern != "" {
		if matched, _ := regexp.MatchString(aw.config.ExcludePattern, path); matched {
			return nil
		}
	}

	// Check inclusion patterns
	if aw.config.IncludePattern != "" {
		if matched, _ := regexp.MatchString(aw.config.IncludePattern, path); !matched {
			return nil
		}
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	packagePath := aw.resolveImportPath(path)

	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if constructor := aw.parseConstructorFunction(fn, node.Name.Name, packagePath, path); constructor != nil {
				aw.constructors = append(aw.constructors, *constructor)
				if aw.config.Verbose {
					if len(constructor.Annotation.As) > 0 {
						console.Info(
							"Found %s: %s.%s as %s",
							constructor.Annotation.Type,
							constructor.PackageName,
							constructor.Name,
							strings.Join(constructor.Annotation.As, ","),
						)
					} else {
						console.Info("Found %s: %s.%s", constructor.Annotation.Type, constructor.PackageName, constructor.Name)
					}
				}
			}
		}

		return true
	})

	return nil
}

func (aw *AutoWired) resolveImportPath(filename string) string {
	// Read go.mod to find module name
	moduleRoot := aw.findModuleRoot(filename)
	if moduleRoot == "" {
		return ""
	}

	moduleName := aw.readModuleName(filepath.Join(moduleRoot, "go.mod"))
	if moduleName == "" {
		return ""
	}

	relPath, err := filepath.Rel(moduleRoot, filepath.Dir(filename))
	if err != nil {
		return ""
	}

	if relPath == "." {
		return moduleName
	}

	return moduleName + "/" + strings.ReplaceAll(relPath, "\\", "/")
}

func (aw *AutoWired) findModuleRoot(startPath string) string {
	dir := filepath.Dir(startPath)

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir { // reached root
			break
		}
		dir = parent
	}

	return ""
}

func (aw *AutoWired) readModuleName(goModPath string) string {
	file, err := os.Open(goModPath)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}

	return ""
}

func (aw *AutoWired) parseConstructorFunction(
	fn *ast.FuncDecl,
	packageName, packagePath, filename string,
) *Constructor {
	if fn.Doc == nil {
		return nil
	}

	annotation := aw.parseAnnotation(fn.Doc.List)
	if annotation == nil {
		return nil
	}

	// Extract comment for documentation
	var comment strings.Builder
	for _, commentLine := range fn.Doc.List {
		text := strings.TrimSpace(strings.TrimPrefix(commentLine.Text, "//"))
		if !strings.HasPrefix(text, "@wired:") {
			if comment.Len() > 0 {
				comment.WriteString("\n")
			}
			comment.WriteString(text)
		}
	}

	uniqueAlias := aw.addImport(packageName, packagePath)

	return &Constructor{
		Name:        fn.Name.Name,
		PackageName: uniqueAlias,
		ImportPath:  packagePath,
		FullPath:    filename,
		Annotation:  *annotation,
		Comment:     comment.String(),
	}
}

func (aw *AutoWired) parseAnnotation(comments []*ast.Comment) *Annotation {
	var annotation *Annotation

	for _, comment := range comments {
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))

		// Parse @wired:provide with options
		if strings.HasPrefix(text, "@wired:") {
			annotation = aw.parseComplexAnnotation(text)
			if annotation != nil {
				break
			}
		}
	}

	return annotation
}

func (aw *AutoWired) parseComplexAnnotation(text string) *Annotation {
	// Regex to parse: @wired:provide(as=Interface1,Interface2 name=myName group=myGroup optional lifecycle=singleton)
	re := regexp.MustCompile(`@wired:(provide|invoke|decorate)(?:\((.*?)\))?`)
	matches := re.FindStringSubmatch(text)

	if len(matches) < 2 {
		return nil
	}

	annotation := &Annotation{Type: matches[1]}

	if len(matches) > 2 && matches[2] != "" {
		// Parse options
		options := matches[2]

		// Parse as=Interface1,Interface2
		if asMatch := regexp.MustCompile(`as=([^\s]+)`).FindStringSubmatch(options); len(
			asMatch,
		) > 1 {
			interfaces := strings.Split(asMatch[1], ",")
			for i, iface := range interfaces {
				interfaces[i] = strings.TrimSpace(iface)
			}
			annotation.As = interfaces
		}

		// Parse name=myName
		if nameMatch := regexp.MustCompile(`name=([^\s]+)`).FindStringSubmatch(options); len(
			nameMatch,
		) > 1 {
			annotation.Name = strings.TrimSpace(nameMatch[1])
		}

		// Parse group=myGroup
		if groupMatch := regexp.MustCompile(`group=([^\s]+)`).FindStringSubmatch(options); len(
			groupMatch,
		) > 1 {
			annotation.Group = strings.TrimSpace(groupMatch[1])
		}

		// Parse optional
		if strings.Contains(options, "optional") {
			annotation.Optional = true
		}

		// Parse lifecycle=singleton
		if lifecycleMatch := regexp.MustCompile(`lifecycle=([^\s]+)`).FindStringSubmatch(options); len(
			lifecycleMatch,
		) > 1 {
			annotation.Lifecycle = strings.TrimSpace(lifecycleMatch[1])
		}
	}

	return annotation
}
