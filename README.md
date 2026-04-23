# Terraform Documentation Extras (tfdocs-extras)

[![GitHub Release](https://img.shields.io/github/v/release/FriendsOfTerraform/tfdocs-extras?include_prereleases)](https://github.com/FriendsOfTerraform/tfdocs-extras/releases/latest) [![GitHub License](https://img.shields.io/github/license/FriendsOfTerraform/tfdocs-extras)](/LICENSE) [![Continuous Integration](https://github.com/FriendsOfTerraform/tfdocs-extras/actions/workflows/ci.yml/badge.svg)](https://github.com/FriendsOfTerraform/tfdocs-extras/actions/workflows/ci.yml)

> [!NOTE]
> This library has reached release candidate status and is ready for use in production; we're using it in our [FriendsOfTerraform modules](https://github.com/FriendsOfTerraform/modules). The Go library API and documentation specification syntax are finalized and will only be receiving non-breaking changes. This repo has a few different components, so be sure to read our [Backward Compatibility Promise](#backward-compatibility-promise) for what this means.

A Go library for parsing an `object()` Terraform type definition string into a documented structure. Support for [documenting nested objects has been a feature request dating back to April 2020](https://github.com/terraform-docs/terraform-docs/issues/242). The biggest challenge is that [Terraform Docs](https://github.com/terraform-docs/terraform-docs) does not parse the `object()` type definition itself and returns it as a raw string; this library fills that gap.

This repository houses a Go library that can parse a Terraform `object()` type definition string (including nested objects) into a structured representation that includes field names, types, optional status, default values, and parsed documentation (including support for doc directives like `@since`, `@example`, etc.). Additionally, it houses a simple CLI tool that uses this project's API for reading a Terraform variable file and outputting the parsed documentation in a GitHub-friendly Markdown format. Last, it also includes a GitHub Action that automatically updates documentation in your Terraform module READMEs whenever you push changes to your Terraform files.

- Usage
  - [as a Go library](#usage-as-a-go-library)
  - [as a CLI tool](usage-as-a-cli-tool)
  - [as a GitHub Action](#usage-as-a-github-action)
- [Documentation Syntax](#documentation-specification)

> [!TIP]
> Check out our [kitchen sink Terraform module](./kitchen-sink-module/), showcasing the capabilities of this library and what it renders.

![Markup preview in a README](./.github/screenshots/tfdocs-extras-preview.png)

## Disclaimer

I am a frontend developer, and besides knowing what Terraform is, I have no experience using it. I built this library to solve a problem we have managing documentation in our [FriendsOfTerraform](https://github.com/FriendsOfTerraform) organization.

The goal is to eventually integrate this library into [Terraform Docs](https://github.com/terraform-docs/terraform-docs) either as a plugin or built-in functionality. I've only recently started learning Go, so I have no idea what the best approach for integrating this library into Terraform Docs is; feedback and direction are highly appreciated.

Like everybody else, I have multiple things in my life that require my attention, so if you'd like to encourage me to work on this, here's how you can do so:

- Provide feedback on this library (e.g. missing features, bug reports, etc.)
- Star this repository / follow me on GitHub
- [Sponsor me on GitHub](https://github.com/sponsors/allejo)

\- [@allejo](https://github.com/allejo)

## Usage as a Go Library

To integrate this library with your own tool, you can install it as a dependency via `go get`.

```bash
go get github.com/FriendsOfTerraform/tfdocs-extras
```

This library provides extra functionality on top of Terraform Docs. It exposes its main functions for parsing Terraform type definitions and generating documentation manifests.

### ParseModuleArgsIntoManifest

The primary function for processing both inputs and outputs from a Terraform module. It returns a `ModuleManifest` containing parsed and structured documentation for required inputs, optional inputs, outputs, and nested object structures.

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/FriendsOfTerraform/tfdocs-extras"
    "github.com/terraform-docs/terraform-docs/print"
    "github.com/terraform-docs/terraform-docs/terraform"
)

func main() {
    // Use terraform-docs to load the module
    config := print.DefaultConfig()
    config.ModuleRoot = "/path/to/tf-module-folder"
    module, _ := terraform.LoadWithOptions(config)

    // Parse both inputs and outputs into a complete manifest
    manifest := tfdocsextras.ParseModuleArgsIntoManifest(module.Inputs, module.Outputs)

    // The manifest contains:
    // - manifest.RequiredInputs: Required input variables
    // - manifest.OptionalInputs: Optional input variables
    // - manifest.Outputs: Module outputs
    // - manifest.Objects: Nested object type definitions

    // Output the manifest as JSON for demonstration purposes
    astJSON, _ := json.MarshalIndent(manifest, "", "  ")
    fmt.Printf("%s\n", astJSON)
}
```

### ParseIntoDocumentedStruct

Parse a Terraform type definition string (such as `object({...})`, `map(object({...}))`, etc.) into a structured representation with documentation. This is useful when you need to parse type definitions directly without going through the full module loading process.

```go
package main

import (
    "encoding/json"
    "fmt"
    "strings"

    "github.com/FriendsOfTerraform/tfdocs-extras"
)

func main() {
    // Parse a type definition string
    typeDefinition := `map(object({
        /// The name of the configuration
        /// @since 1.0.0
        name = string

        /// The port number
        /// @since 1.0.0
        port = optional(number, 8080)

        /// Nested configuration settings
        settings = optional(object({
            /// Enable debug mode
            debug = bool

            /// Timeout in seconds
            timeout = number
        }))
    }))`

    // Parse the type definition
    documented, err := tfdocsextras.ParseIntoDocumentedStruct(typeDefinition, "server_config")
    if err != nil {
        panic(err)
    }

    // Access the parsed structure
    fmt.Printf("Type: %s\n", documented.DataTypeStr)
    fmt.Printf("Nested Type: %s\n", *documented.NestedDataType)

    // Iterate through fields
    for _, field := range documented.Properties {
        fmt.Printf("\nField: %s\n", field.Name)
        fmt.Printf("  Type: %s\n", field.DataTypeStr)
        fmt.Printf("  Optional: %v\n", field.Optional)
        fmt.Printf("  Description: %s\n", strings.Join(field.Documentation.Content, " "))

        // Access directives
        for _, directive := range field.Documentation.Directives {
            fmt.Printf("  @%s: %s\n", directive.Name, directive.RawContent)
        }
    }

    // Convert to JSON
    jsonOutput, _ := json.MarshalIndent(documented, "", "  ")
    fmt.Printf("\nJSON Output:\n%s\n", jsonOutput)
}
```

## Usage as a CLI Tool

This project includes a rudimentary CLI tool that reads a Terraform module folder and outputs the parsed variable documentation in Markdown format.

[Download the latest build from GitHub](https://github.com/FriendsOfTerraform/tfdocs-extras/releases/latest)

> [!IMPORTANT]
>
> This CLI tool is not intended to be a replacement for Terraform Docs and will become obsolete once its functionality is integrated into Terraform Docs either as a plugin or built-in feature.

The tool accepts a single argument specifying the path to a Terraform module folder. It will read the module folder using Terraform Docs, parse the variable definitions using this library, and write the output to a `README.md` file in the module folder.

```bash
./tfdocs-extras /path/to/TerraformModules/aws/route53
```

The README requires specific markers to identify where to insert the generated documentation. The generated Markdown will be inserted between the following markers:

```html
<!-- TFDOCS_EXTRAS_START -->

<!-- TFDOCS_EXTRAS_END -->
```

### JSON Output

Pass the `-json` flag to print a JSON representation of the parsed module manifest to stdout instead of updating the README. This is useful for piping the output into other tools or for debugging.

```bash
./tfdocs-extras -json /path/to/TerraformModules/aws/route53
```

## Usage as a GitHub Action

This repository is also published as a GitHub Action. It automatically downloads the appropriate binary for the runner's OS and architecture, scans the specified directories for `README.md` files containing both the `TFDOCS_EXTRAS_START`/`TFDOCS_EXTRAS_END` markers, and processes each one.

Each entry in `directories` is either an explicit path or a glob pattern:

| Pattern                         | Behavior                                                            |
| ------------------------------- | ------------------------------------------------------------------- |
| `./modules/aws/vpc`             | Process this single directory.                                      |
| `./modules/aws/*`               | Process every direct subdirectory of `aws/`.                        |
| `./modules/{aws,azure,vault}/*` | Process every direct subdirectory of `aws/`, `azure/` and `vault/`. |

Directories that do not contain a `README.md` with both the `TFDOCS_EXTRAS_START`/`TFDOCS_EXTRAS_END` markers are silently skipped.

```yaml
- uses: FriendsOfTerraform/tfdocs-extras@main
  with:
    directories: |
      ./modules/aws/vpc
      ./modules/aws/s3
      ./modules/gcp/*
      ./modules/{azure,vault}/*
```

> [!TIP]
>
> This GitHub Action is distributed as part of the tool's main repository and therefore does not provide any Git references (e.g., tags like `@v1`) like most other actions. Any tags that exist in this repository (and therefore, this action) will correlate to the software and **not** the action. For ease of use, we recommend using `@main` to always pull the latest action and executable.
>
> However, if you'd like to harden your workflow to avoid [ref confusion][ref-confusion] pin the workflow to an exact hash. Additionally, to guarantee replicable behavior in addition to pinning your `uses:`, specify the exact version of the CLI tool you'd like to use with the `version` argument of this action.

### Inputs

| Input              | Required | Default                                                              | Description                                                                                                         |
| ------------------ | -------- | -------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| `directories`      | Yes      |                                                                      | Newline-separated list of Terraform module directories to process. Supports `*` as a wild card pattern.             |
| `version`          | No       | `latest`                                                             | Version of tfdocs-extras to download (e.g. `v0.1.0`).                                                               |
| `token`            | No       | `${{ github.token }}`                                                | GitHub token used to download the release binary.                                                                   |
| `commit`           | No       | `false` (boolean)                                                    | Commit any README.md changes after processing. Mutually exclusive with `fail_on_diff`.                              |
| `commit_message`   | No       | `chore: update tfdocs-extras documentation`                          | Commit message. Only used when `commit` is `true`.                                                                  |
| `commit_author`    | No       | `github-actions[bot] <github-actions[bot]@users.noreply.github.com>` | Commit author in `Name <email>` format. Only used when `commit` is `true`.                                          |
| `commit_branch`    | No       | Current branch                                                       | Branch to push the commit to. Only used when `commit` is `true`.                                                    |
| `fail_on_diff`     | No       | `false` (boolean)                                                    | Exit with a non-zero status if any README.md files were modified. Mutually exclusive with `commit`.                 |
| `json_output_file` | No       |                                                                      | Path to write the aggregated JSON output. Use this for large outputs that exceed GitHub Actions' output size limit. |

### Outputs

| Output   | Description                                                                                                                                                                                                                  |
| -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `result` | JSON array of module manifests, one object per processed directory. Only set if the total size is below ~1MB (GitHub Actions output limit). For larger outputs, use the `json_output_file` input to write to a file instead. |

### Examples

#### Auto-commit updated documentation

Automatically regenerate and commit documentation whenever Terraform files change on the main branch.

```yaml
on:
  push:
    branches: [main]
    paths: ["**.tf"]

jobs:
  docs:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v6
      - uses: FriendsOfTerraform/tfdocs-extras@main
        with:
          directories: |
            ./modules/aws/vpc
            ./modules/aws/s3
          commit: true
          commit_author: "github-actions[bot] <github-actions[bot]@users.noreply.github.com>"
```

#### Enforce up-to-date documentation in pull requests

Fail the CI check if a pull request contains Terraform changes without updated documentation.

```yaml
on:
  pull_request:
    paths: ["**.tf"]

jobs:
  docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
      - uses: FriendsOfTerraform/tfdocs-extras@main
        with:
          directories: |
            ./modules/aws/vpc
            ./modules/aws/s3
            ./modules/azure/*
            ./modules/vault/*
          fail_on_diff: true
```

#### Use the JSON output in a downstream step

```yaml
- uses: actions/checkout@v6
- uses: FriendsOfTerraform/tfdocs-extras@main
  id: tfdocs
  with:
    directories: ./modules/aws/vpc
- run: echo '${{ steps.tfdocs.outputs.result }}'
```

#### Write JSON output to a file for large module sets

When processing many modules or modules with large manifests, the aggregated JSON may exceed GitHub Actions' output size limit (~1MB). Use `json_output_file` to write the output to a file instead.

```yaml
- uses: actions/checkout@v6
- uses: FriendsOfTerraform/tfdocs-extras@main
  with:
    directories: |
      ./modules/*
    json_output_file: ./tfdocs-manifests.json
- name: Upload manifest
  uses: actions/upload-artifact@v4
  with:
    name: tfdocs-manifests
    path: ./tfdocs-manifests.json
```

## Documentation Specification

The goal of this library is to support Terraform module creators to document their nested variables inline using comments instead of needing to maintain the documentation separately. We introduce two main features:

- Doc Blocks (denoted by `///` or `/** ... */`)
- Doc Directives (denoted by `@directive-name`)

Here's an example of how to document a nested object variable in Terraform:

```terraform
# variables.tf

variable "access_points" {
  type = map(object({
    /// Configures the permissions EFS use to create the specified root
    /// directory if the directory does not already exist
    ///
    /// @since 1.0.0
    root_directory_creation_permissions = optional(object({
      /// Owner group ID for the access point's root directory, if the directory
      /// does not already exist. Valid value: `0 - 4294967295`
      ///
      /// @since 1.0.0
      owner_group_id = number

      /// Owner user ID for the access point's root directory, if the directory
      /// does not already exist. Valid value: `0 - 4294967295`
      ///
      /// @since 1.0.0
      owner_user_id = number
    }))

    /// Path on the EFS file system to expose as the root directory to NFS
    /// clients using the access point. A path can have up to four
    /// subdirectories; `root_directory_creation_permissions` must be
    /// specified if the root path does not exist.
    ///
    /// @since 1.0.0
    root_directory_path = optional(string, "/")
  }))
  description = <<EOT
    Configures [access points][efs-access-point].

    @link {efs-access-point} https://docs.aws.amazon.com/efs/latest/ug/efs-access-points.html
    @example "Access Points" #access-points
    @since 1.0.0
  EOT
  default = {}
}
```

### Doc Blocks

Doc blocks are multi-line comments that start with `///` or are enclosed within `/** ... */`. The triple-slash style is supported to allow for a more compact syntax and to differentiate from regular comments (i.e. `//`).

```terraform
root_directory_creation_permissions = optional(object({
  /**
   * Owner group ID for the access point's root directory, if the directory
   * does not already exist. Valid value: `0 - 4294967295`
   *
   * @since 1.0.0
   */
  owner_group_id = number

  /// Owner user ID for the access point's root directory, if the directory
  /// does not already exist. Valid value: `0 - 4294967295`
  ///
  /// @since 1.0.0
  owner_user_id = number
}))
```

### Directives

Directives are special annotations within doc blocks that provide additional metadata about the documented field. They start with an `@` symbol followed by the directive name and content. Unless otherwise specified, the content of the directive is treated as a single line. However, some directives (e.g., `@type`) support multi-line content.

> [!TIP]
> Using HEREDOC syntax for a variable's `description` attribute allows you to use `@`-directives for a top-level variable.

#### `@deprecated`

Use this directive to mark a field as deprecated and optionally provide migration guidance. This directive's only argument is the deprecation message; everything after the directive name is considered part of the deprecation message. The presence of this directive will mark the field as deprecated in the generated documentation.

```text
@deprecated Use `new_field_name` instead.
```

#### `@enum`

When a field can accept only a specific set of values, you can document the allowed values using the `@enum` directive. The different values are delimited by a vertical pipe (i.e. `|`); spaces around the pipe are optional.

```text
@enum value1|value2|value3
```

#### `@example`

The `@example` directive allows you to provide usage examples for the documented field. They will be listed alongside the field's documentation under the "Examples" section. It accepts two space-separated parameters: a title surrounded by double quotes and a link (URL or anchor).

```text
@example "Advanced Usage Example" #heading-id
@example "Basic Usage Example" https://example.com/usage-example
```

#### `@link`

There are two types of links you can create using the `@link` directive: named links and reference links.

##### Named Links

A named link allows you to create a link with a custom display name and will be displayed alongside the field's documentation in a special "Links" section. It accepts two space-separated parameters: a title surrounded by required double quotes and a URL.

```text
@link "Some Resource Documentation" https://example.com/some/resource
```

##### Reference Links

A reference link uses curly braces (i.e. `{}`) to create [link reference definitions](https://spec.commonmark.org/0.31.2/#link-reference-definition). Reference links are useful when you want to reuse the same link multiple times in your documentation without repeating the URL. Another use case is when the URL is long and would clutter the documentation if displayed inline.

```text
A link to [Some Resource Documentation][resource-id].

@link {resource-id} https://example.com/some/resource
```

> [!IMPORTANT]
> If you have manually created a link reference definition in your README, you may reference it in your documentation using the standard Markdown syntax for reference links without needing to use the `@link` directive.
>
> Keep in mind that following this pattern will require you to manage the link reference definitions yourself and ensure that they are included in the README, as this library will not be aware of manually created link reference definitions and will not include them in the generated documentation.

> [!WARNING]
> When using reference links, make sure to use unique identifiers (per module) for each link reference definition to avoid conflicts. If two `@link` directives use the same identifier, they will overwrite each other, and only the last one will be included in the generated documentation.

#### `@regex`

The `@regex` directive allows you to specify a regular expression between `/` delimiters that the field's value must match. After the pattern, you can provide example values that conform to the regex. Double quotes are supported for example values that contain spaces and can be escaped with a backslash (`\`).

```text
@regex /(Average|Minimum|Maximum) (<=|<|>=|>) (\d+)/ "Average >= 20" "Minimum < 10" "Maximum <= 100"
```

#### `@since`

The version when the field was introduced.

```text
@since 1.0.0
```

#### `@type`

The `@type` directive allows you to specify an output's data type where the type information is not automatically available. This directive supports multiline type definitions for complex nested structures.

```terraform
output "nat_gateways" {
  description = <<EOT
    Map of default NAT gateways. The key of the map is the NAT gateway's name

    @since 1.0.0
    @type map(object({
      /// The availability zone of the NAT gateway
      /// @since 1.0.0
      availability_zone = string

      /// The association ID of the Elastic IP address that's associated with the NAT Gateway
      /// @since 1.0.0
      association_id = string
    }))
  EOT
  value = aws_nat_gateway.this
}
```

> [!NOTE]
> The `@type` directive supports multiline definitions. When the directive contains an opening bracket (`(` or `{`), the parser will automatically continue reading subsequent lines until all brackets are balanced.

> [!WARNING]
> The `@type` directive is only respected for outputs. For variables, the type information is already available in the variable definition and does not require a directive, therefore usage of `@type` in an input is ignored.

## Backward Compatibility Promise

The [semantic versioning](https://semver.org/) used in this repository represents the version of the Go library itself, the sole focus of this project. At the moment, I have no interest in making the provided CLI tool and GitHub Action their own standalone projects; these are provided as convenient implementations of the library's functionality until this project is integrated with Terraform Docs.

The GitHub Action downloads the appropriate binary of our CLI to quickly integrate it into your project and therefore will always be compatible with each other. However, our GitHub Action does not and will not have any version numbering. We recommend pinning the action to `@main` to always get the latest version of the action (and therefore CLI tool) OR pinning to the exact hash to avoid [ref confusion][ref-confusion].

I will try my absolute best to avoid introducing breaking changes to the CLI tool and GitHub Action, but I cannot guarantee it. If there is enough interest in supporting the CLI tool as a light-weight alternative to Terraform Docs for users who only want the documentation parsing features without needing the full functionality of Terraform Docs, I may consider making the CLI tool its own standalone project with its own versioning and backward compatibility promise. The GitHub Action will likely remain as part of this repository since it's just a thin wrapper around the CLI tool, but it will be versioned alongside the CLI tool to ensure compatibility. Let me know if you would like to see this happen by [creating an issue](/issues) and providing your use case for why you want a standalone CLI tool. The more use cases and interest I see, the more likely I am to make and maintain a standalone CLI tool.

### What is covered by the backward compatibility promise?

The current stable and production-ready version of the library is the 0.0.x series. This library's 1.0.0 release will likely happen when it's successfully integrated with Terraform Docs and has received sufficient exposure and usage in the community. Until then, we will only introduce non-breaking changes in 0.x releases to the following:

- The public API of the Go library, i.e.,
  - [`ParseModuleArgsIntoManifest`](#parsemoduleargsintomanifest) and its return type (`ModuleManifest`)
  - [`ParseIntoDocumentedStruct`](#parseintodocumentedstruct) and its return type (`DocumentedStruct`) and its nested types (`DocumentedField`, `Directive`, etc.)
- The JSON schema of the structs in this library (e.g., `ModuleManifest`, `DocumentedStruct`, etc.) so you can safely use the JSON output in your own tools
- The [documentation specification](#documentation-specification), i.e.,
  - The comment syntax for doc blocks
  - Directive syntax should be treated as function signatures, so any changes to directive syntax will be non-breaking as long as the old syntax is still supported (e.g., if we change `@enum value1|value2|value3` to `@enum [value1, value2, value3]`, we will still support the old pipe-delimited syntax for backward compatibility)
  - The set of supported directives and their semantics (e.g., `@deprecated` will always mark a field as deprecated in the generated documentation regardless of any syntax changes to the directive)

## License

[MIT](./LICENSE)

[ref-confusion]: https://docs.zizmor.sh/audits/#ref-confusion
