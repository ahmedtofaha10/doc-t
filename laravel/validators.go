package laravel

import (
	"log"
	"os"
)

func checkDirOrFileExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) { // check dir
		log.Fatalf("[%s] does not exists :/", path)
	}
}
func (project *Project) ValidateProjectStructure() {
	checkDirOrFileExists(project.BasePath)                     // check base dir
	checkDirOrFileExists(project.BasePath + "\\app")           // check app dir
	checkDirOrFileExists(project.BasePath + "\\composer.json") // check composer reqs
}
