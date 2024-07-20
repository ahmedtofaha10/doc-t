package laravel

type Project struct {
	BasePath    string
	Env         map[string]string
	Dependinces ProjectDependinces
	Tables      map[string]map[string]string
}
type Controller struct {
	Name    string
	Traits  string
	Methods []string
}
type ProjectDependinces struct {
	Name       string
	Require    map[string]string
	RequireDev map[string]string
}
type Route struct {
	Uri         string
	Method      string
	FullUri     string
	Action      string
	Middlewares []string
	Prefixes    []string
}
type MI map[string]interface{}
