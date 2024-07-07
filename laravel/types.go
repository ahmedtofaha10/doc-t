package laravel

type Project struct {
	BasePath string
}
type ProjectDependinces struct {
	Name       string
	Require    map[string]string
	RequireDev map[string]string
}
