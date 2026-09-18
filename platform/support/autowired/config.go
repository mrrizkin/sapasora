package autowired

type Config struct {
	Directories    []string
	OutputFile     string
	PackageName    string
	Verbose        bool
	Watch          bool
	ConfigFile     string
	ExcludePattern string
	IncludePattern string
}
