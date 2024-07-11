package laravel

import (
	"fmt"
	"os"
	"strings"
)

func Documenting(base_path string) {
	project := Project{BasePath: base_path}
	project.GetMigrationsData()
	// project.ValidateProjectStructure()
	// project.GetComposerDependincies()
	// project.GetEnvFileData()
	// writeDocumentationFile(project)

}
func writeProjectMeta(file *os.File, project Project) {
	projectName := string(project.Env["APP_NAME"])
	file.Write([]byte(fmt.Sprintf("# %s\n", projectName)))
	file.Write([]byte("## Requirements\n"))
	file.Write([]byte(fmt.Sprintf("> php : %s \n", project.Dependinces.Require["php"])))
	file.Write([]byte("## Dependinces\n"))
	for k := range project.Dependinces.Require {
		if k == "php" {
			continue
		}
		file.Write([]byte(fmt.Sprintf("> %s : %s \n", k, project.Dependinces.Require[k])))
	}
}
func writeProjectModels(file *os.File, project Project) {
	file.Write([]byte("## Models\n"))
	models, err := os.ReadDir(project.BasePath + "\\app\\Models")
	if err != nil {
		panic(err)
	}
	for modelIndex := range models {
		model := models[modelIndex]
		modelName := strings.Replace(model.Name(), ".php", "", 1)
		file.Write([]byte(fmt.Sprintf("### %s\n", modelName)))
	}
}
func writeDocumentationFile(project Project) {
	filePath := project.BasePath + "\\doc-t.md"
	if _, err := os.Stat(filePath); !os.IsNotExist(err) { // check file
		os.Remove(filePath)
	}
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writeProjectMeta(file, project)
	writeProjectModels(file, project)

}
