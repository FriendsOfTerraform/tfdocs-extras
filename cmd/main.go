package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/allejo/tfdocs-extra"
)

func main() {
	data, err := os.ReadFile("multiline_string.txt")
	if err != nil {
		log.Fatal(err)
	}

	documented, err := tfdocextras.ParseIntoDocumentedStruct(string(data), "root_object")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	astJSON, _ := json.MarshalIndent(documented, "", "  ")
	fmt.Printf("%s\n", astJSON)
}
