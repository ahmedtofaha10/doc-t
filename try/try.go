package main

import (
	"fmt"
	"regexp"
)

func main() {
	// The input string
	input := `$table->enum('publish',['Active','Not Active'])->default('Not Active');`

	// Regular expression to match the name between the first set of single quotes
	reSingleQuotes := regexp.MustCompile(`'([^']*)'`)
	matchSingleQuotes := reSingleQuotes.FindStringSubmatch(input)
	if len(matchSingleQuotes) > 1 {
		fmt.Println("Name between first single quotes:", matchSingleQuotes[1])
	} else {
		fmt.Println("No match found for single quotes")
	}

	// Regular expression to match the first chain method name
	reChainMethod := regexp.MustCompile(`->(\w+)\(`)
	matchChainMethod := reChainMethod.FindStringSubmatch(input)
	if len(matchChainMethod) > 1 {
		fmt.Println("First chain method name:", matchChainMethod[1])
	} else {
		fmt.Println("No match found for chain method")
	}
}
