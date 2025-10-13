package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/FriendsOfTerraform/tfdocs-extras"
	"github.com/terraform-docs/terraform-docs/print"
	"github.com/terraform-docs/terraform-docs/terraform"
)

type ExpandedModuleInput struct {
	terraform.Input

	Extras tfdocextras.ObjectGroup `json:"extras,omitempty"`
}

func recordNested(group tfdocextras.ObjectField, record map[string][]tfdocextras.ObjectField) {
	if group.NestedDataType == nil {
		return
	}

	if group.Fields != nil && len(group.Fields) > 0 {
		record[*group.NestedDataType] = group.Fields
	}

	for _, field := range group.Fields {
		recordNested(field, record)
	}
}

func main() {
	config := print.DefaultConfig()
	config.ModuleRoot = os.Args[1]

	module, err := terraform.LoadWithOptions(config)
	if err != nil {
		log.Fatal(err)
	}

	var requiredInputs []ExpandedModuleInput
	var optionalInputs []ExpandedModuleInput
	nestedInputs := make(map[string][]tfdocextras.ObjectField)

	for _, input := range module.Inputs {
		var extras tfdocextras.ObjectGroup
		if input.Type != "" {
			if documented, err := tfdocextras.ParseIntoDocumentedStruct(string(input.Type), input.Name); err == nil && documented != nil {
				extras = *documented
			}
		}

		expanded := ExpandedModuleInput{
			Input:  *input,
			Extras: extras,
		}

		if expanded.Extras.Optional {
			optionalInputs = append(optionalInputs, expanded)
		} else {
			requiredInputs = append(requiredInputs, expanded)
		}

		for _, field := range expanded.Extras.Fields {
			recordNested(field, nestedInputs)
		}
	}

	astModule, _ := json.MarshalIndent(nestedInputs, "", "  ")
	fmt.Printf("%s\n", astModule)
}
