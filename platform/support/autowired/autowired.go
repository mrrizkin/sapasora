// Package autowired contains the auto register constructors for go fx
package autowired

func DefaultConfig() Config {
	return Config{
		Directories:    []string{"internal"},
		OutputFile:     "internal/app/autowired/autowired.go",
		PackageName:    "autowired",
		Verbose:        false,
		Watch:          false,
		ConfigFile:     "",
		ExcludePattern: "",
		IncludePattern: "",
	}
}

func NewAutoWired(config ...Config) *AutoWired {
	c := DefaultConfig()
	if len(config) > 0 {
		conf := config[0]
		if conf.Directories != nil {
			c.Directories = conf.Directories
		}
		if conf.OutputFile != "" {
			c.OutputFile = conf.OutputFile
		}
		if conf.PackageName != "" {
			c.PackageName = conf.PackageName
		}
		if conf.Verbose {
			c.Verbose = conf.Verbose
		}
		if conf.Watch {
			c.Watch = conf.Watch
		}
		if conf.ConfigFile != "" {
			c.ConfigFile = conf.ConfigFile
		}
		if conf.ExcludePattern != "" {
			c.ExcludePattern = conf.ExcludePattern
		}
		if conf.IncludePattern != "" {
			c.IncludePattern = conf.IncludePattern
		}
	}

	return &AutoWired{
		imports:      make(map[string]string),
		importPaths:  make(map[string]string),
		aliasCount:   make(map[string]int),
		constructors: make([]Constructor, 0),
		config:       c,
		templates:    NewTemplates(),
	}
}
