package laravel

func Documenting(base_path string) {
	project := Project{BasePath: base_path}
	project.ValidateProjectStructure()
	project.GetComposerDependincies()

}
