package laravel

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
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
	fmt.Println(dependinces.Require["php"])
}
