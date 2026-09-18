package kuproy

import (
	"sapasora/platform/support/arr"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
)

var tagRegex = regexp.MustCompile(`^(?P<name>\w+)\((?P<value>[^)]+)\)$`)

func parseArgs(cli *Kuproy) (*Generator, error) {
	gen := Generator{
		ModulePath: getModuleName(),
	}
	args := cli.Args
	switch cli.Task {
	default:
		if len(cli.Args) < 4 {
			return nil, fmt.Errorf("invalid number of arguments")
		}

		gen.Module = args[0]
		gen.Entity = args[1]
		gen.TableName = args[2]

		args = args[3:]

	case "policy":
		if len(args) != 2 {
			return nil, fmt.Errorf("invalid number of arguments")
		}

		gen.Entity = args[0]
		gen.Module = args[1]

		// ignore rest of args
		args = []string{}

	case "migration":
		if len(args) < 1 {
			return nil, fmt.Errorf("invalid number of arguments")
		}

		gen.Entity = args[0]
		timestamp := time.Now().Format("2006_01_02_150405")
		filename := fmt.Sprintf("%s_%s.go", timestamp, gen.Entity)
		gen.MigrationName = filename

		args = args[1:]
	}

	fields, err := parseFields(args)
	if err != nil {
		return nil, err
	}
	gen.Fields = fields

	return &gen, nil
}

func parseFields(args []string) ([]Field, error) {
	var parsed []Field
	for _, f := range args {
		parts := strings.Split(f, ":")
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "Warning: skipping invalid field format '%s'\n", f)
			continue
		}

		field := Field{
			Name: parts[0],
			Type: parts[1],
		}

		// Handle attributes like :unique, :required
		if len(parts) > 2 {
			field.Attributes = arr.Map(parts[2:], parseTag)

			if !slices.Contains(parts[2:], "json") {
				field.Attributes = arr.Merge(
					[]string{fmt.Sprintf(`json:"%s"`, strcase.ToSnake(field.Name))},
					field.Attributes,
				)
			}
		}

		parsed = append(parsed, field)
	}

	return parsed, nil
}

func parseTag(arg string) string {
	matches := tagRegex.FindStringSubmatch(arg)
	if matches == nil {
		return ""
	}
	return fmt.Sprintf(`%s:"%s"`, matches[1], matches[2])
}
