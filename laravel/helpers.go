package laravel

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func (project *Project) GetComposerDependincies() {
	file, err := os.Open(project.BasePath + "\\composer.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var dependinces ProjectDependinces
	byteValue, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(byteValue, &dependinces)
	if err != nil {
		panic(err)
	}
	project.Dependinces = dependinces
}
func (project *Project) GetMigrationsData() {
	migrations, err := os.ReadDir(project.BasePath + "\\database\\migrations")
	if err != nil {
		panic(err)
	}
	tables := map[string]map[string]string{}
	for migrationIndex := range migrations {
		migration := migrations[migrationIndex]
		if !migration.IsDir() {
			fmt.Println(migration.Name())
			table, columns := extractMigrationFileData(project.BasePath + "\\database\\migrations\\" + migration.Name())
			tables[table] = columns
		}
	}

}
func extractMigrationFileData(filePath string) (table string, columns map[string]string) {
	file, err := os.Open(filePath)
	columns = map[string]string{}
	if err != nil {
		panic(err)
	}
	defer file.Close()
	columnName := ""
	columnType := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		columnName = ""
		columnType = ""
		input := scanner.Text()
		if strings.Contains(input, "down()") {
			break
		}
		reTableName := regexp.MustCompile(`Schema::create\('([^']*)'`)
		reTableName2 := regexp.MustCompile(`Schema::table\('([^']*)'`)
		reSingleQuotes := regexp.MustCompile(`'([^']*)'`)
		matchSingleQuotes := reSingleQuotes.FindStringSubmatch(input)
		matchTableName2 := reTableName2.FindStringSubmatch(input)
		matchTableName := reTableName.FindStringSubmatch(input)
		if len(matchTableName) > 1 {
			table = matchTableName[1]
			continue
		}
		if len(matchTableName2) > 1 {
			table = matchTableName2[1]
			continue
		}
		if len(matchSingleQuotes) > 1 {
			columnName = matchSingleQuotes[1]
		}

		// Regular expression to match the first chain method name
		reChainMethod := regexp.MustCompile(`->(\w+)\(`)
		matchChainMethod := reChainMethod.FindStringSubmatch(input)
		if len(matchChainMethod) > 1 {
			columnType = matchChainMethod[1]
		}
		if columnName != "" && columnType != "" {
			if columnName == "" {
				columnName = columnType
			}
			// fmt.Printf("%s : %s \n", columnName, columnType)
			columns[columnName] = columnType
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return table, columns
}
func (project *Project) GetEnvFileData() {
	file, err := os.Open(project.BasePath + "\\.env")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	env := map[string]string{}
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		splited := strings.Split(text, "=")
		env[splited[0]] = splited[1]
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	project.Env = env
	// fmt.Println(env["APP_NAME"])
}
