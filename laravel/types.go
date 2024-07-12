package laravel

type Project struct {
	BasePath    string
	Env         map[string]string
	Dependinces ProjectDependinces
	Tables      map[string]map[string]string
}
type ProjectDependinces struct {
	Name       string
	Require    map[string]string
	RequireDev map[string]string
}
