package laravel

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"
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
			// fmt.Println(migration.Name())
			table, columns := extractMigrationFileData(project.BasePath + "\\database\\migrations\\" + migration.Name())
			tables[table] = columns
		}
	}
	project.Tables = tables

}
func splitCamelCase(word string) []string {
	var words []string
	var currentWord strings.Builder

	for i, r := range word {
		if unicode.IsUpper(r) && i != 0 {
			words = append(words, currentWord.String())
			currentWord.Reset()
		}
		currentWord.WriteRune(r)
	}
	words = append(words, currentWord.String())

	return words
}

// Pluralize a word by appending "s" or "es"
func Pluralize(word string) string {
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "sh") || strings.HasSuffix(word, "ch") || strings.HasSuffix(word, "x") || strings.HasSuffix(word, "z") {
		return strings.ToLower(word + "es")
	}
	return strings.ToLower(word + "s")
}

// Pluralize a compound word
func PluralizeCompoundWord(compoundWord string) (string, string) {
	words := splitCamelCase(compoundWord)
	if len(words) == 0 {
		return compoundWord, ""
	}
	solidTableName := strings.ToLower(strings.Join(words, "_"))

	// Pluralize the last word of the compound word
	words[len(words)-1] = Pluralize(words[len(words)-1])
	tableName := strings.ToLower(strings.Join(words, "_"))
	return tableName, solidTableName
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
		if columnName == "" {
			columnName = columnType
		}
		if columnName != "" && columnType != "" {
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
func (project *Project) readRoutes() []Route {
	// apiRoutes := []map[string]string{}
	apiFile, err := os.Open(project.BasePath + "\\routes\\api.php")
	if err != nil {
		panic(err)
	}
	defer apiFile.Close()
	scanner := bufio.NewScanner(apiFile)
	regexMethod := regexp.MustCompile(`::([^']*)\(`)
	currentGroups := []map[string]string{}
	routes := []Route{}
	for scanner.Scan() {
		methodMatches := regexMethod.FindStringSubmatch(scanner.Text())
		if strings.Contains(scanner.Text(), "});") {
			if len(currentGroups) > 0 {
				currentGroups = currentGroups[0 : len(currentGroups)-1]
			}
		}
		if len(methodMatches) > 1 {
			method := methodMatches[1]

			if method == "middleware" {
				regexMiddleware := regexp.MustCompile(`middleware\(\[\s*(.*?)\s*\]\)`)
				middlewareMatches := regexMiddleware.FindStringSubmatch(scanner.Text())
				// fmt.Println(middlewareMatches)
				if len(middlewareMatches) > 1 {
					newGroup := map[string]string{}
					newGroup["type"] = "middleware"
					newGroup["value"] = middlewareMatches[1]
					currentGroups = append(currentGroups, newGroup)
				}
			} else if method == "group" {
				regexPrefix := regexp.MustCompile(`'prefix'\s*=>\s*'([^']*)'`)
				prefixMatches := regexPrefix.FindStringSubmatch(scanner.Text())
				// fmt.Println(middlewareMatches)
				if len(prefixMatches) > 1 {
					newGroup := map[string]string{}
					newGroup["type"] = "prefix"
					newGroup["value"] = prefixMatches[1]
					currentGroups = append(currentGroups, newGroup)
				}
			} else {
				regexUri := regexp.MustCompile(`::` + method + `\('([^']*)',`)
				uriMatches := regexUri.FindStringSubmatch(scanner.Text())
				if len(uriMatches) > 1 {
					uri := uriMatches[1]
					// fmt.Println(method, uri, currentGroups)
					middlewares := []string{}
					prefixes := []string{}
					fullUri := "api/"
					for groupIndex := range currentGroups {
						if currentGroups[groupIndex]["type"] == "middleware" {
							middlewares = append(middlewares, currentGroups[groupIndex]["value"])
						} else if currentGroups[groupIndex]["type"] == "prefix" {
							prefixes = append(prefixes, currentGroups[groupIndex]["value"])
							fullUri += currentGroups[groupIndex]["value"]
						}
					}
					fullUri += uri
					route := Route{}
					// TODO get and add the action also
					route.Middlewares = middlewares
					route.Prefixes = prefixes
					route.Uri = uri
					route.Method = method
					route.FullUri = strings.Replace(fullUri, "//", "/", -1)
					routes = append(routes, route)
				}
			}

		}
	}
	return routes
}
