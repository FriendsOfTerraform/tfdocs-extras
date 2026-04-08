variable "required_string" {
  description = <<-EOT
    Required primitive string.

    @since 1.0.0
    @example "Production" "production"
  EOT

  type = string
}

variable "optional_number_with_default" {
  description = <<-EOT
    Optional number with an explicit default.

    @since 1.0.0
    @deprecated Use `request_timeout_seconds` instead.
  EOT

  type    = number
  default = 30
}

variable "enum_like_value" {
  description = <<-EOT
    Enum-like string values.

    @enum development|staging|production
    @default staging
  EOT

  type    = string
  default = "staging"
}

variable "regex_validated_id" {
  description = <<-EOT
    Identifier that follows a strict pattern.

    @regex /^[a-zA-Z0-9_-]{5}$/ abcd4 efgh_ ijkl-
    @since 1.2.0
  EOT

  type = string

  validation {
    condition     = can(regex("^[a-zA-Z0-9_-]{5}$", var.regex_validated_id))
    error_message = "regex_validated_id must be 5 characters and match ^[a-zA-Z0-9_-]{5}$"
  }
}

variable "code_fence_examples" {
  description = <<-EOT
    Demonstrates fenced code blocks in descriptions.

    ```terraform
    code_fence_examples = {
      enabled = true
      labels  = ["alpha", "beta"]
    }
    ```

    ~~~terraform
    code_fence_examples = {
      enabled = false
      labels  = []
    }
    ~~~

    @since 1.1.0
  EOT

  type = object({
    enabled = bool
    labels  = list(string)
  })

  default = {
    enabled = false
    labels  = []
  }
}

variable "named_links" {
  description = <<-EOT
    Demonstrates named and reference links.

    @link "Terraform Type Constraints" https://developer.hashicorp.com/terraform/language/expressions/type-constraints
    @link {terraform-functions} https://developer.hashicorp.com/terraform/language/functions
    @example "Functions docs" #terraform-functions
  EOT

  type    = list(string)
  default = []
}

variable "list_of_strings_optional" {
  description = <<-EOT
    Optional list(string) represented via default empty list.

    @since 1.0.0
    @type optional(list(string), [])
  EOT

  type    = list(string)
  default = []
}

variable "map_of_objects_documented" {
  description = <<-EOT
    Map of documented object definitions.

    @type map(object({ ... }))
    @since 2.0.0
  EOT

  type = map(object({
    /// Human-friendly display name.
    /// @since 1.0.0
    display_name = string

    /// Environment tier.
    /// @enum dev|stage|prod
    tier = string

    /// Optional tags for each object.
    tags = optional(list(string), [])

    /// Optional metadata map.
    metadata = optional(map(string), {})
  }))

  default = {}
}

variable "list_of_objects_mixed_docs" {
  description = <<-EOT
    List of objects where some properties are documented and others are not.

    @type list(object({ ... }))
    @example "Minimal item" "[{ name = \"api\", port = 8080 }]"
  EOT

  type = list(object({
    /// Name of the service.
    /// @required true
    name = string

    /// Public port for the service.
    /// @default 8080
    port = optional(number, 8080)

    protocol = optional(string, "http")

    # Intentionally undocumented nested object.
    health_check = optional(object({
      path                = optional(string, "/healthz")
      interval_seconds    = optional(number, 30)
      timeout_seconds     = optional(number, 5)
      expected_statuses   = optional(list(number), [200])
      enable_status_check = optional(bool, true)
    }), null)
  }))

  default = []
}

variable "deeply_nested_object" {
  description = <<-EOT
    Deeply nested object that mixes required and optional fields.

    @type object({ ... })
    @since 3.0.0
  EOT

  type = object({
    /// Required top-level name.
    name = string

    /// Optional top-level toggle.
    enabled = optional(bool, true)

    /// Optional nested object with null default semantics.
    deployment = optional(object({
      /// Required region for deployment.
      region = string

      /// Optional replica count.
      replicas = optional(number, 2)

      /// Optional list of zones.
      zones = optional(list(string), [])

      /// Nested map(object(...)) for autoscaling profiles.
      autoscaling_profiles = optional(map(object({
        min = number
        max = number
        # Undocumented property by design.
        target_cpu = optional(number, 70)
      })), {})
    }), null)
  })

  default = {
    name    = "example"
    enabled = true
  }
}

variable "object_with_block_comment_docs" {
  description = <<-EOT
    Uses block-style docs in object property type definitions.

    @type object({ ... })
  EOT

  type = object({
    /**
     * Primary key.
     * @since 1.0.0
     */
    id = string

    /**
     * Optional summary.
     * @deprecated Prefer `details`.
     */
    summary = optional(string, null)

    details = optional(object({
      title = string
      body  = string
    }), null)
  })

  default = {
    id = "example-id"
  }
}
