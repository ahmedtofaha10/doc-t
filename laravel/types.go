package laravel

type Project struct {
	BasePath    string
	Env         map[string]string
	Dependinces ProjectDependinces
}
type ProjectDependinces struct {
	Name       string
	Require    map[string]string
	RequireDev map[string]string
}
