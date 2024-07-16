package laravel

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func Documenting(base_path, destination string) {
	project := Project{BasePath: base_path}
	project.GetMigrationsData()
	project.ValidateProjectStructure()
	project.GetComposerDependincies()
	project.GetEnvFileData()
	writeDocumentationFile(project, destination)

}
func writeProjectMeta(file *os.File, project Project) {
	projectName := string(project.Env["APP_NAME"])
	file.Write([]byte(fmt.Sprintf("# %s\n", projectName)))
	file.Write([]byte("## Requirements\n"))
	file.Write([]byte(fmt.Sprintf("> * php : %s \n", project.Dependinces.Require["php"])))
	file.Write([]byte("## Dependinces\n"))
	for k := range project.Dependinces.Require {
		if k == "php" {
			continue
		}
		file.Write([]byte(fmt.Sprintf("> * %s : %s \n", k, project.Dependinces.Require[k])))
	}
}
func getModelFileData(modelPath string) (tableName string, methods []string) {
	tableName = ""
	file, err := os.Open(modelPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		tablePattern := `\$table = '([^']*)'`
		methodsPattern := `function\s+(\w+)\s*\(`
		reTable := regexp.MustCompile(tablePattern)
		reMethods := regexp.MustCompile(methodsPattern)

		tableMatches := reTable.FindStringSubmatch(text)

		if len(tableMatches) > 1 {
			tableName = tableMatches[1]
			continue
		}
		methodMatches := reMethods.FindStringSubmatch(text)
		if len(methodMatches) > 1 {
			methods = append(methods, methodMatches[1])
			continue
		}
	}
	return tableName, methods

}
func writeProjectModels(file *os.File, project Project) {
	file.Write([]byte("## Models\n"))
	models, err := os.ReadDir(project.BasePath + "\\app\\Models")
	if err != nil {
		panic(err)
	}
	tableName := ""
	var methods []string
	solidTableName := ""
	for modelIndex := range models {
		model := models[modelIndex]
		tableName, methods = getModelFileData(project.BasePath + "\\app\\Models\\" + model.Name())
		modelName := strings.Replace(model.Name(), ".php", "", 1)
		if len(tableName) == 0 {
			tableName, solidTableName = PluralizeCompoundWord(modelName)
		}
		file.Write([]byte(fmt.Sprintf("### %s (%s)\n", modelName, tableName)))
		columns := project.Tables[tableName]
		if len(columns) == 0 {
			columns = project.Tables[solidTableName]
		}
		if len(columns) > 0 {
			file.Write([]byte(fmt.Sprintf("\n**[ %s ]**\n", "Columns")))
		}
		for columnName := range columns {
			file.Write([]byte(fmt.Sprintf("> * %s (%s)\t\n", columnName, columns[columnName])))
		}
		if len(methods) > 0 {
			file.Write([]byte(fmt.Sprintf("\n**[ %s ]**\t\n", "Methods")))
			for methodIndex := range methods {
				file.Write([]byte(fmt.Sprintf("> * %s()\t\n", methods[methodIndex])))
			}
		}
	}
}
func writeDocumentationFile(project Project, destination string) {
	filePath := destination + "\\doc-t.md"
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
	project.readRoutes()
}
