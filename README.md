> [!WARNING]
>
> This library is stable enough for experimental use; however, it is not recommended for production code just yet.

> [!IMPORTANT]
>
> I am a frontend developer, and besides knowing what Terraform is, I have no experience using it. I built this library to solve a problem we have managing documentation in our FriendsOfTerraform organization.
>
> My goal is to eventually integrate this library into [Terraform Docs](https://github.com/terraform-docs/terraform-docs) either as a plugin or built-in functionality. Considering that I barely started learning Go about a month ago, I have no idea what the best approach for integrating this library into Terraform Docs is; feedback and direction are highly appreciated.
>
> Like everybody else, I have multiple things in my life that require my attention, so if you'd like to encourage me to work on this, here's how you can do so:
> 
> - Provide feedback on this library (e.g. missing features, bug reports, etc.)
> - Star this repository / follow me on GitHub
> - [Sponsor me on GitHub](https://github.com/sponsors/allejo)
>
> \- [@allejo](https://github.com/allejo)

# Terraform Documentation Extras (tfdocs-extras)

A Go library for parsing an `object()` Terraform type definition string into a documented structure.  Support for [documenting nested objects has been a feature request dating back to April 2020](https://github.com/terraform-docs/terraform-docs/issues/242). The biggest challenge is that [Terraform Docs](https://github.com/terraform-docs/terraform-docs) does not parse the `object()` type definition itself and returns it as a raw string; this library fills that gap.

After writing an [RFC for features that would be useful to the community](https://github.com/FriendsOfTerraform/modules/issues/38), this library implements the core functionality of parsing documented `object()` type definitions. This library will be used by a plugin for Terraform Docs to improve documentation generation.

## Installation

```bash
go get github.com/FriendsOfTerraform/tfdocs-extras
```

## Public API

The library exposes only **one main function** and its related types:

### Main Function

```go
func ParseIntoDocumentedStruct(input string, name string) (*ObjectGroup, error)
```

Parses a Terraform type definition string into a documented object group. This is the **only exported function** in the library.

### Public Types

These types are exported so you can work with the results:

- `ObjectGroup` - The main result type containing parsed fields and metadata
- `ObjectField` - Individual field within an object structure  
- `VariableMetadata` - Metadata about variables (name, type, optional status, etc.)
- `FieldDocBlock` - Parsed documentation for a field
- `DocDirective` - Documentation directives like `@since`, `@param`, etc.

### Method

- `(*ObjectGroup).GetObjectName() string` - Returns CamelCase name for the object group

## Example Usage

Let's say you have the following Terraform variable definition:

```terraform
variable "default_capacity_provider_strategy" {
  type = map(object({
    /// The relative percentage of the total number of launched tasks that should use the specified capacity provider.
    /// `weight` is taken into consideration only after the `base` count of tasks has been satisfied.
    ///
    /// @since 1.0.0
    base = optional(number, 0)

    /**
     * The number of tasks, at a minimum, to run on the specified capacity provider. Only one capacity provider in a
     * capacity provider strategy can have `base` defined. Defaults to `0`.
     *
     * @since 1.0.0
     */
    weight = number
  }))
  description = "Specify the default capacity provider strategy that is used when creating services in the cluster"
  default     = {}
}
```

Terraform Docs will output the type as:

```terraform
map(object({
    /// The relative percentage of the total number of launched tasks that should use the specified capacity provider.
    /// `weight` is taken into consideration only after the `base` count of tasks has been satisfied.
    ///
    /// @since 1.0.0
    base = optional(number, 0)

    /**
     * The number of tasks, at a minimum, to run on the specified capacity provider. Only one capacity provider in a
     * capacity provider strategy can have `base` defined. Defaults to `0`.
     *
     * @since 1.0.0
     */
    weight = number
}))
```

By calling `ParseIntoDocumentedGroup()` with the above string, you will get a structured representation of the object, including field names, types, optional status, default values, and parsed documentation.

```json
{
  "Name": "root_object",
  "Documentation": {
    "Content": [],
    "Directives": []
  },
  "DataTypeStr": "map(object(RootObject))",
  "Optional": false,
  "DefaultValue": null,
  "Fields": [
    {
      "Name": "base",
      "Documentation": {
        "Content": [
          "The relative percentage of the total number of launched tasks that should use the specified capacity provider.",
          "`weight` is taken into consideration only after the `base` count of tasks has been satisfied."
        ],
        "Directives": [
          {
            "Name": "since",
            "Content": "1.0.0"
          }
        ]
      },
      "DataTypeStr": "number",
      "Optional": true,
      "DefaultValue": "0",
      "NestedDataType": null
    },
    {
      "Name": "weight",
      "Documentation": {
        "Content": [
          "The number of tasks, at a minimum, to run on the specified capacity provider. Only one capacity provider in a",
          "capacity provider strategy can have `base` defined. Defaults to `0`."
        ],
        "Directives": [
          {
            "Name": "since",
            "Content": "1.0.0"
          }
        ]
      },
      "DataTypeStr": "number",
      "Optional": false,
      "DefaultValue": null,
      "NestedDataType": null
    }
  ],
  "ParentDataType": "map"
}
```

## Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/FriendsOfTerraform/tfdocs-extras"
)

func main() {
    input := `optional(object({
        /// The user's name
        /// @since 1.0.0
        name = string
        
        /// The user's age
        /// @default 18
        age = optional(number, 18)
        
        /// User's address
        address = object({
            street = string
            city = string
        })
    }))`
    
    documented, err := tfdocextras.ParseIntoDocumentedStruct(input, "root_object")
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    astJSON, _ := json.MarshalIndent(documented, "", "  ")
    fmt.Printf("%s\n", astJSON)
}
```

## License

MIT
