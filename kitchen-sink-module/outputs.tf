output "required_string" {
  description = <<-EOT
    Echoes the required string input.

    @since 1.0.0
    @type string
  EOT

  value = var.required_string
}

output "service_ports" {
  description = <<-EOT
    Extracted list of service ports.

    @type list(number)
  EOT

  value = [for svc in var.list_of_objects_mixed_docs : svc.port]
}

output "documented_map" {
  description = <<-EOT
    Pass-through for documented map(object(...)) input.

    @since 2.0.0
    @type map(object({
      /// Human-friendly display name.
      display_name = string

      /// Environment tier.
      tier = string

      /// Optional tags for each object.
      tags = optional(list(string), [])

      /// Optional metadata map.
      metadata = optional(map(string), {})
    }))
  EOT

  value = var.map_of_objects_documented
}

output "sensitive_summary" {
  description = <<-EOT
    Sensitive summary value for output handling tests.

    @type string
  EOT

  value     = "${var.required_string}:${var.enum_like_value}"
  sensitive = true
}
