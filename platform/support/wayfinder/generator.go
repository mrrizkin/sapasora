package wayfinder

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sapasora/platform/support/arr"
	"slices"
	"sort"
	"strings"
)

// create a regex that is not \w or \d
var nonWordOrDigitRegex = regexp.MustCompile(`[^\w\d]`)

// Generate generates both routes and controllers
func (w *Wayfinder) Generate() error {
	if err := w.generateWayfinder(); err != nil {
		return err
	}
	if err := w.generateRoutes(); err != nil {
		return err
	}
	if err := w.generateControllers(); err != nil {
		return err
	}
	return nil
}

// generateRoutes generates TypeScript route definitions organized by router structure
func (w *Wayfinder) generateRoutes() error {
	if err := w.prepareOutputPath(w.routeOutputPath); err != nil {
		w.log.Error("Failed to prepare route output path: %v", err)
		return err
	}

	for folder, routes := range w.groupRoutesByRouter() {
		if err := w.generateRouteContent(folder, routes); err != nil {
			w.log.Error("Failed to generate route folder %q: %v", folder, err)
			return err
		}
	}
	return nil
}

// generateControllers generates TypeScript route definitions organized by controller structure
func (w *Wayfinder) generateControllers() error {
	if err := w.prepareOutputPath(w.controllerOutputPath); err != nil {
		w.log.Error("Failed to prepare controller output path: %v", err)
		return err
	}

	groups := w.groupRoutesByController()
	for folder, files := range groups {
		if err := w.generateControllerContent(folder, files); err != nil {
			w.log.Error("Failed to generate controller folder %q: %v", folder, err)
			return err
		}
	}

	tree, err := w.convertRouteMapToTree(groups)
	if err != nil {
		w.log.Error("Failed to convert route map to tree: %v", err)
		return err
	}
	if err := w.processTree(tree, w.controllerOutputPath, "actions"); err != nil {
		w.log.Error("Failed to process tree: %v", err)
		return err
	}
	return nil
}

// generateWayfinder generates TypeScript wayfinder helpers
func (w *Wayfinder) generateWayfinder() error {
	if err := os.MkdirAll(w.wayfinderOutputPath, 0755); err != nil {
		w.log.Error("Failed to create wayfinder output path: %v", err)
		return err
	}

	file, err := os.Create(filepath.Join(w.wayfinderOutputPath, "wayfinder.ts"))
	if err != nil {
		w.log.Error("Failed to create wayfinder output file: %v", err)
		return err
	}
	defer file.Close()

	err = w.templates.Render(file, "wayfinder.ts.gotmpl", nil)
	if err != nil {
		w.log.Error("Failed to generate wayfinder helpers: %v", err)
		return err
	}
	return nil
}

// prepareOutputPath cleans the output path if configured to do so
func (w *Wayfinder) prepareOutputPath(outputPath string) error {
	if !w.cleanOutputPath {
		return nil
	}

	w.log.Info("Cleaning output path: %s", outputPath)
	return os.RemoveAll(outputPath)
}

// generateRouteContent generates TypeScript files for a route folder
func (w *Wayfinder) generateRouteContent(folder string, routes RouterGroup) error {
	folderPath := filepath.Join(w.routeOutputPath, folder)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return fmt.Errorf("failed to create folder: %w", err)
	}

	file, err := os.Create(filepath.Join(folderPath, "index.ts"))
	if err != nil {
		return err
	}
	defer file.Close()

	templateData := RoutesTemplateData{
		Routes:          w.buildRoutesTemplateData(routes),
		Name:            nonWordOrDigitRegex.ReplaceAllString(folder, "_"),
		IsDefaultExport: len(folder) > 0,
		IsForm:          true,
	}

	return w.templates.Render(file, "routes.ts.gotmpl", templateData)
}

// generateControllerContent generates TypeScript files for a controller folder
func (w *Wayfinder) generateControllerContent(folder string, files StructPointerGroup) error {
	folderPath := filepath.Join(w.controllerOutputPath, folder)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return fmt.Errorf("failed to create folder: %w", err)
	}

	for fileName, functions := range files {
		templateData := RoutesTemplateData{
			Routes:          w.buildRoutesTemplateData(functions),
			Name:            fileName,
			IsDefaultExport: len(fileName) > 0,
			IsForm:          true,
		}

		if fileName == "" {
			fileName = "index"
		}

		file, err := os.Create(filepath.Join(folderPath, fileName+".ts"))
		if err != nil {
			return err
		}
		defer file.Close()

		if err := w.templates.Render(file, "routes.ts.gotmpl", templateData); err != nil {
			return fmt.Errorf("failed to generate controller file %s: %w", fileName, err)
		}
	}

	return nil
}

func (w *Wayfinder) buildRoutesTemplateData(routes RouterGroup) []RouteData {
	var data []RouteData
	for name, routeGroup := range routes {
		if len(routeGroup) == 0 {
			continue
		}

		data = append(data, RouteData{
			Location: routeGroup[0].Handler.Location,
			Path:     routeGroup[0].Route.Path,
			Methods:  arr.Map(routeGroup, func(r Route) string { return r.Route.Method }),
			Name:     name,
			Params:   routeGroup[0].Route.Params,
		})
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Name < data[j].Name
	})

	return data
}

func (w *Wayfinder) groupRoutesByRouter() StructPointerGroup {
	grouped := make(StructPointerGroup)
	for folder, routes := range w.routes {
		groupByName := arr.GroupBy(routes, func(route Route) string {
			return route.Name
		})
		grouped[folder] = groupByName
	}
	return grouped
}

func (w *Wayfinder) groupRoutesByController() PathGroup {
	grouped := make(PathGroup)

	for _, routeGroup := range w.routes {
		for _, route := range routeGroup {
			handler := route.Handler
			if handler.Name == "" {
				continue
			}

			if grouped[handler.Path] == nil {
				grouped[handler.Path] = make(StructPointerGroup)
			}

			if grouped[handler.Path][handler.StructPointer] == nil {
				grouped[handler.Path][handler.StructPointer] = make(RouterGroup)
			}

			grouped[handler.Path][handler.StructPointer][handler.Name] = append(
				grouped[handler.Path][handler.StructPointer][handler.Name],
				route,
			)
		}
	}

	return grouped
}

func (w *Wayfinder) convertRouteMapToTree(input PathGroup) (map[string]any, error) {
	root := make(map[string]any)

	for path, controllers := range input {
		segments := strings.Split(path, "/")
		current := root

		for _, seg := range segments {
			if _, exists := current[seg]; !exists {
				current[seg] = make(map[string]any)
			}
			current = current[seg].(map[string]any)
		}

		for controller := range controllers {
			current[controller] = make(map[string]any)
		}
	}

	return root, nil
}

func (w *Wayfinder) processTree(node map[string]any, folderPath, branchName string) error {
	keys := arr.Keys(node)
	if len(keys) == 0 {
		return nil
	}
	slices.Sort(keys)

	if err := os.MkdirAll(folderPath, 0o755); err != nil {
		return fmt.Errorf("failed to create folder %s: %w", folderPath, err)
	}

	if _, err := os.Stat(filepath.Join(folderPath, "index.ts")); err == nil {
		return nil
	}

	file, err := os.Create(filepath.Join(folderPath, "index.ts"))
	if err != nil {
		return fmt.Errorf("failed to create index.ts in %s: %w", folderPath, err)
	}
	defer file.Close()

	data := IndexTemplateData{
		Actions: keys,
		Name:    branchName,
	}
	if err := w.templates.Render(file, "index.ts.gotmpl", data); err != nil {
		return fmt.Errorf("failed to render template in %s: %w", folderPath, err)
	}

	for key, val := range node {
		child, ok := val.(map[string]any)
		if !ok {
			continue
		}
		childPath := filepath.Join(folderPath, key)
		if err := w.processTree(child, childPath, key); err != nil {
			return err
		}
	}

	return nil
}
