package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	data, err := os.ReadFile("multiline_string.txt")
	if err != nil {
		log.Fatal(err)
	}

	ast, err := ParseAst(string(data))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Print the AST structure more clearly
	astJSON, _ := json.MarshalIndent(ast, "", "  ")
	fmt.Printf("%s\n", astJSON)
}
