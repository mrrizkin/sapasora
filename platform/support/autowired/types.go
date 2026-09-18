package autowired

type Annotation struct {
	Type      string   // provide, invoke, decorate
	As        []string // interfaces to provide as
	Name      string   // named dependency
	Group     string   // group name
	Optional  bool     // optional dependency
	Lifecycle string   // singleton, request, etc.
}

type Constructor struct {
	Name        string
	PackageName string
	ImportPath  string
	FullPath    string
	Annotation  Annotation
	Comment     string
}

type AutoWired struct {
	constructors []Constructor
	imports      map[string]string
	importPaths  map[string]string
	aliasCount   map[string]int
	config       Config
	templates    *Templates
}
