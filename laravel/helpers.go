package laravel

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
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
