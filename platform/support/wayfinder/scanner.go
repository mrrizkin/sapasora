package wayfinder

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var (
	actionsRegex       = regexp.MustCompile(`.*\/app\/http\/controllers.(.*)-fm`)
	pathRegex          = regexp.MustCompile(`.*\/app\/http\/controllers\/([^.]*)`)
	pointerStructRegex = regexp.MustCompile(`\(\*(.*)\)`)
	functionNameRegex  = regexp.MustCompile(`\.([a-zA-Z0-9_]*)-fm`)
)

func groupRoutes(routes []fiber.Route) map[string][]Route {
	grouped := make(map[string][]Route)

	for _, route := range routes {
		if route.Name == "" {
			continue
		}

		folderPath := extractFolderPath(route.Name)
		handler := extractHandler(route)
		name := extractRouteName(route.Name)

		grouped[folderPath] = append(grouped[folderPath], Route{
			Name:    name,
			Route:   route,
			Handler: handler,
		})
	}

	return grouped
}

func extractFolderPath(routeName string) string {
	paths := strings.Split(routeName, ".")
	for i, path := range paths {
		paths[i] = strings.ReplaceAll(path, "/", "_")
	}
	return filepath.Join(paths[:len(paths)-1]...)
}

func extractRouteName(routeName string) string {
	paths := strings.Split(routeName, ".")
	return paths[len(paths)-1]
}

func extractHandler(route fiber.Route) Handler {
	for _, handler := range route.Handlers {
		if path, functionName, structPointer, location := parseHandlerFunction(handler); functionName != "" {
			return Handler{
				Name:          functionName,
				Path:          filepath.Join("app/http/controllers", path),
				StructPointer: structPointer,
				Location:      location,
			}
		}
	}
	return Handler{}
}

func parseHandlerFunction(
	handler fiber.Handler,
) (path, functionName, structPointer, location string) {
	ptr := reflect.ValueOf(handler).Pointer()
	fn := runtime.FuncForPC(ptr)
	if fn == nil {
		return // couldn't get function info
	}

	funcPath := fn.Name()

	// Validate that the function path contains expected patterns before applying regex
	if !strings.Contains(funcPath, "app/http/controllers") {
		return // not a controller function
	}

	// Extract path using regex with validation
	if match := actionsRegex.FindStringSubmatch(funcPath); len(match) > 1 {
		// Use the match if available
	} else {
		return // couldn't match expected pattern
	}

	if match := pathRegex.FindStringSubmatch(funcPath); len(match) > 1 {
		path = strings.ReplaceAll(match[1], ".", "/") // convert dots to path separators
	}

	if match := functionNameRegex.FindStringSubmatch(funcPath); len(match) > 1 {
		functionName = match[1]
	}

	if match := pointerStructRegex.FindStringSubmatch(funcPath); len(match) > 1 {
		structPointer = match[1]
	}

	// Only look for location if we have both struct pointer and function name
	if structPointer != "" && functionName != "" {
		searchDir := filepath.Join("internal/app/http/controllers", path)
		file, line, err := findMethodLocation(structPointer, functionName, searchDir)
		if err != nil {
			// Log the error but continue without location information
		} else {
			location = filepath.Join(path, file) + ":" + strconv.Itoa(line)
		}
	}

	return path, functionName, structPointer, location
}

func findMethodLocation(
	structName, methodName, dirPath string,
) (filename string, line int, err error) {
	// Get all .go files in the specified directory
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return "", 0, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	fset := token.NewFileSet()

	// Process each .go file individually
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())

		// Parse individual file
		astFile, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err != nil {
			// Skip files that can't be parsed, but don't fail entirely
			continue
		}

		// Check if this file's package name is relevant (skip main packages)
		if astFile.Name.Name == "main" {
			continue
		}

		// Iterate through all declarations in the file
		for _, decl := range astFile.Decls {
			// Check if the declaration is a function declaration
			if fn, ok := decl.(*ast.FuncDecl); ok {
				// Check if it's a method (has a receiver) and matches our criteria
				if fn.Recv != nil && len(fn.Recv.List) > 0 {
					receiver := fn.Recv.List[0].Type

					// Check if receiver type matches structName
					// Handle both pointer and non-pointer receivers
					receiverTypeName := getReceiverTypeName(receiver)
					if receiverTypeName == structName && fn.Name.Name == methodName {
						// Found the method, return file and line number
						pos := fset.Position(fn.Pos())
						relPath, err := filepath.Rel(".", filePath)
						if err != nil {
							relPath = filePath
						}
						return relPath, pos.Line, nil
					}
				}
			}
		}
	}

	return "", 0, fmt.Errorf("method %s.%s not found in %s", structName, methodName, dirPath)
}

// Helper function to get the type name from a receiver
func getReceiverTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr: // pointer receiver *Type
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name
		}
	case *ast.Ident: // value receiver Type
		return t.Name
	}
	return ""
}
