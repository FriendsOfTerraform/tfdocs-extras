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

## Example Usage

Let's say you have the following Terraform variable definition:

```terraform
variable "network_interface" {
  type = object({
    /// List of security group IDs attached to this ENI
    /// @since 1.0.0
    security_group_ids                 = list(string)

    /// Specify the subnet ID this ENI is created on
    /// @since 1.0.0
    subnet_id                          = string

    /// Additional tags for the ENI
    /// @since 1.0.0
    additional_tags                    = optional(map(string), {})

    /// Specify the description of the ENI
    /// @since 1.0.0
    description                        = optional(string)

    /// Enables [elastic fabric adapter][elastic-fabric-adapter]
    /// @since 1.0.0
    enable_elastic_fabric_adapter      = optional(bool, false)

    /// Controls if traffic is routed to the instance when the destination address does not match the instance. Used for NAT or VPNs
    /// @since 1.0.0
    enable_source_destination_checking = optional(bool, true)

    /// Configures custom private IP addresses for the ENI.
    /// @since 1.0.0
    private_ip_addresses = optional(object({
      /// List of private IPv4 addresses to assign to the ENI, the first address will be used as the primary IP address
      /// @since 1.0.0
      ipv4 = optional(list(string))
    }))

    /// Assigns a private CIDR range, either automatically or manually, to the ENI. By assigning [prefixes][ec2-prefixes], you scale and simplify the management of applications, including container and networking applications that require multiple IP addresses on an instance. Network interfaces with prefixes are supported with [instances built on the Nitro System][nitro-system-type].
    /// @since 1.0.0
    prefix_delegation = optional(object({
      /// Configures prefix delegation for IPV4
      /// @since 1.0.0
      ipv4 = optional(object({
        /// Sepcify the number of prefixes AWS chooses from your VPC subnet’s IPv4 CIDR block and assigns it to your network interface. Mutually exclusive to `custom_prefixes`
        /// @since 1.0.0
        auto_assign_count = optional(number)

        /// Specify the prefixes from your VPC subnet’s CIDR block to assign it to your network interface. Mutually exclusive to `auto_assign_count`
        /// @since 1.0.0
        custom_prefixes   = optional(list(string))
      }))
    }))
  })
  description = "Configures the primary network interface"
}
```

Terraform Docs will output the type as:

```terraform
object({
    /// List of security group IDs attached to this ENI
    /// @since 1.0.0
    security_group_ids                 = list(string)

    /// Specify the subnet ID this ENI is created on
    /// @since 1.0.0
    subnet_id                          = string

    /// Additional tags for the ENI
    /// @since 1.0.0
    additional_tags                    = optional(map(string), {})

    /// Specify the description of the ENI
    /// @since 1.0.0
    description                        = optional(string)

    /// Enables [elastic fabric adapter][elastic-fabric-adapter]
    /// @since 1.0.0
    enable_elastic_fabric_adapter      = optional(bool, false)

    /// Controls if traffic is routed to the instance when the destination address does not match the instance. Used for NAT or VPNs
    /// @since 1.0.0
    enable_source_destination_checking = optional(bool, true)

    /// Configures custom private IP addresses for the ENI.
    /// @since 1.0.0
    private_ip_addresses = optional(object({
      /// List of private IPv4 addresses to assign to the ENI, the first address will be used as the primary IP address
      /// @since 1.0.0
      ipv4 = optional(list(string))
    }))

    /// Assigns a private CIDR range, either automatically or manually, to the ENI. By assigning [prefixes][ec2-prefixes], you scale and simplify the management of applications, including container and networking applications that require multiple IP addresses on an instance. Network interfaces with prefixes are supported with [instances built on the Nitro System][nitro-system-type].
    /// @since 1.0.0
    prefix_delegation = optional(object({
      /// Configures prefix delegation for IPV4
      /// @since 1.0.0
      ipv4 = optional(object({
        /// Sepcify the number of prefixes AWS chooses from your VPC subnet’s IPv4 CIDR block and assigns it to your network interface. Mutually exclusive to `custom_prefixes`
        /// @since 1.0.0
        auto_assign_count = optional(number)

        /// Specify the prefixes from your VPC subnet’s CIDR block to assign it to your network interface. Mutually exclusive to `auto_assign_count`
        /// @since 1.0.0
        custom_prefixes   = optional(list(string))
      }))
    }))
  })
```

By calling `ParseIntoDocumentedGroup()` with the above string, you will get a structured representation of the object, including field names, types, optional status, default values, and parsed documentation.

```json
{
  "name": "root_object",
  "documentation": {
    "content": [],
    "directives": []
  },
  "dataType": "object(RootObject)",
  "optional": false,
  "fields": [
    {
      "name": "security_group_ids",
      "documentation": {
        "content": [
          "List of security group IDs attached to this ENI"
        ],
        "directives": [
          {
            "name": "since",
            "content": "1.0.0"
          }
        ]
      },
      "dataType": "list(string)",
      "optional": false
    },
    {
      "name": "subnet_id",
      "documentation": {
        "content": [
          "Specify the subnet ID this ENI is created on"
        ],
        "directives": [
          {
            "name": "since",
            "content": "1.0.0"
          }
        ]
      },
      "dataType": "string",
      "optional": false
    },
    {
      "name": "additional_tags",
      "documentation": {
        "content": [
          "Additional tags for the ENI"
        ],
        "directives": [
          {
            "name": "since",
            "content": "1.0.0"
          }
        ]
      },
      "dataType": "map(string)",
      "optional": true
    },
    {
      "name": "description",
      "documentation": {
        "content": [
          "Specify the description of the ENI"
        ],
        "directives": [
          {
            "name": "since",
            "content": "1.0.0"
          }
        ]
      },
      "dataType": "string",
      "optional": true
    },
    {
      "name": "enable_elastic_fabric_adapter",
      "documentation": {
        "content": [
          "Enables [elastic fabric adapter][elastic-fabric-adapter]"
        ],
        "directives": [
          {
            "name": "since",
            "content": "1.0.0"
          }
        ]
      },
      "dataType": "bool",
      "optional": true,
      "defaultValue": "false"
    },
    {
      "name": "enable_source_destination_checking",
      "documentation": {
        "content": [
          "Controls if traffic is routed to the instance when the destination address does not match the instance. Used for NAT or VPNs"
        ],
        "directives": [
          {
            "name": "since",
            "content": "1.0.0"
          }
        ]
      },
      "dataType": "bool",
      "optional": true,
      "defaultValue": "true"
    },
    {
      "name": "private_ip_addresses",
      "documentation": {
        "content": [
          "Configures custom private IP addresses for the ENI."
        ],
        "directives": [
          {
            "name": "since",
            "content": "1.0.0"
          }
        ]
      },
      "dataType": "object(PrivateIpAddresses)",
      "optional": true,
      "nestedDataType": [
        {
          "name": "ipv4",
          "documentation": {
            "content": [
              "List of private IPv4 addresses to assign to the ENI, the first address will be used as the primary IP address"
            ],
            "directives": [
              {
                "name": "since",
                "content": "1.0.0"
              }
            ]
          },
          "dataType": "list(string)",
          "optional": true
        }
      ]
    },
    {
      "name": "prefix_delegation",
      "documentation": {
        "content": [
          "Assigns a private CIDR range, either automatically or manually, to the ENI. By assigning [prefixes][ec2-prefixes], you scale and simplify the management of applications, including container and networking applications that require multiple IP addresses on an instance. Network interfaces with prefixes are supported with [instances built on the Nitro System][nitro-system-type]."
        ],
        "directives": [
          {
            "name": "since",
            "content": "1.0.0"
          }
        ]
      },
      "dataType": "object(PrefixDelegation)",
      "optional": true,
      "nestedDataType": [
        {
          "name": "ipv4",
          "documentation": {
            "content": [
              "Configures prefix delegation for IPV4"
            ],
            "directives": [
              {
                "name": "since",
                "content": "1.0.0"
              }
            ]
          },
          "dataType": "object(Ipv4)",
          "optional": true,
          "nestedDataType": [
            {
              "name": "auto_assign_count",
              "documentation": {
                "content": [
                  "Sepcify the number of prefixes AWS chooses from your VPC subnet’s IPv4 CIDR block and assigns it to your network interface. Mutually exclusive to `custom_prefixes`"
                ],
                "directives": [
                  {
                    "name": "since",
                    "content": "1.0.0"
                  }
                ]
              },
              "dataType": "number",
              "optional": true
            },
            {
              "name": "custom_prefixes",
              "documentation": {
                "content": [
                  "Specify the prefixes from your VPC subnet’s CIDR block to assign it to your network interface. Mutually exclusive to `auto_assign_count`"
                ],
                "directives": [
                  {
                    "name": "since",
                    "content": "1.0.0"
                  }
                ]
              },
              "dataType": "list(string)",
              "optional": true
            }
          ]
        }
      ]
    }
  ],
  "parentDataType": ""
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
