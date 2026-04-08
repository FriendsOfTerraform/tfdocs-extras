# Kitchen Sink Terraform Module

This module exists purely as fixture content to exercise parsing and rendering edge cases.

<!-- TFDOCS_EXTRAS_START -->






## Inputs

### Required



    

    

    

    

    
<table><thead><tr><th>Type</th><th align="left" width="100%">Name</th><th>Default&nbsp;Value</th></tr></thead><tbody>
        <tr>
    <td><code>string</code></td>
    <td width="100%">regex_validated_id</td><td><em>n/a</em></td>
</tr>
<tr><td colspan="3">

Identifier that follows a strict pattern.

    

    
**Regex Pattern:**
```
^[a-zA-Z0-9_-]{5}$
```


        
Example Matches:
- `abcd4`
- `efgh_`
- `ijkl-`

    

    

    
**Since:** 1.2.0
        


</td></tr>
<tr>
    <td><code>string</code></td>
    <td width="100%">required_string</td><td><em>n/a</em></td>
</tr>
<tr><td colspan="3">

Required primitive string.

    

    

    
**Examples:**
- [Production]("production")

    

    
**Since:** 1.0.0
        


</td></tr>
</tbody></table>


### Optional



    

    

    

    

    
<table><thead><tr><th>Type</th><th align="left" width="100%">Name</th><th>Default&nbsp;Value</th></tr></thead><tbody>
        <tr>
    <td><code>object(<a href="#code_fence_examples">code_fence_examples</a>)</code></td>
    <td width="100%">code_fence_examples</td><td><code>{
  "enabled": false,
  "labels": []
}</code></td>
</tr>
<tr><td colspan="3">

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

    

    

    

    

    
**Since:** 1.1.0
        


</td></tr>
<tr>
    <td><code>object(<a href="#deeply_nested_object">deeply_nested_object</a>)</code></td>
    <td width="100%">deeply_nested_object</td><td><code>{
  "enabled": true,
  "name": "example"
}</code></td>
</tr>
<tr><td colspan="3">

Deeply nested object that mixes required and optional fields.

    

    

    

    

    
**Since:** 3.0.0
        


</td></tr>
<tr>
    <td><code>string</code></td>
    <td width="100%">enum_like_value</td><td><code>"staging"</code></td>
</tr>
<tr><td colspan="3">

Enum-like string values.

    
**Allowed Values:**
- `development`
- `staging`
- `production`

    

    

    

    


</td></tr>
<tr>
    <td><code>list(object(<a href="#list_of_objects_mixed_docs">list_of_objects_mixed_docs</a>))</code></td>
    <td width="100%">list_of_objects_mixed_docs</td><td><code>[]</code></td>
</tr>
<tr><td colspan="3">

List of objects where some properties are documented and others are not.

    

    

    
**Examples:**
- [Minimal item]("[{ name = \"api\", port = 8080 }]")

    

    


</td></tr>
<tr>
    <td><code>list(string)</code></td>
    <td width="100%">list_of_strings_optional</td><td><code>[]</code></td>
</tr>
<tr><td colspan="3">

Optional list(string) represented via default empty list.

    

    

    

    

    
**Since:** 1.0.0
        


</td></tr>
<tr>
    <td><code>map(object(<a href="#map_of_objects_documented">map_of_objects_documented</a>))</code></td>
    <td width="100%">map_of_objects_documented</td><td><code>{}</code></td>
</tr>
<tr><td colspan="3">

Map of documented object definitions.

    

    

    

    

    
**Since:** 2.0.0
        


</td></tr>
<tr>
    <td><code>list(string)</code></td>
    <td width="100%">named_links</td><td><code>[]</code></td>
</tr>
<tr><td colspan="3">

Demonstrates named and reference links.

    

    

    
**Examples:**
- [Functions docs](#terraform-functions)

    
**Links:**
- [Terraform Type Constraints](https://developer.hashicorp.com/terraform/language/expressions/type-constraints)

    


</td></tr>
<tr>
    <td><code>object(<a href="#object_with_block_comment_docs">object_with_block_comment_docs</a>)</code></td>
    <td width="100%">object_with_block_comment_docs</td><td><code>{
  "id": "example-id"
}</code></td>
</tr>
<tr><td colspan="3">

Uses block-style docs in object property type definitions.

    

    

    

    

    


</td></tr>
<tr>
    <td><code>number</code></td>
    <td width="100%">optional_number_with_default</td><td><code>30</code></td>
</tr>
<tr><td colspan="3">

Optional number with an explicit default.

    

    

    

    

    
**Since:** 1.0.0
        
**Deprecated:** Use `request_timeout_seconds` instead.
        


</td></tr>
</tbody></table>


## Outputs



    

    

    

    

    
<table><thead><tr><th>Type</th><th align="left" width="100%">Name</th><th>Sensitive</th></tr></thead><tbody>
        <tr>
    <td><code>map(object(<a href="#documented_map">documented_map</a>))</code></td>
    <td width="100%">documented_map</td><td></td>
</tr>
<tr><td colspan="3">

Pass-through for documented map(object(...)) input.

    

    

    

    

    
**Since:** 2.0.0
        


</td></tr>
<tr>
    <td><code>string</code></td>
    <td width="100%">required_string</td><td></td>
</tr>
<tr><td colspan="3">

Echoes the required string input.

    

    

    

    

    
**Since:** 1.0.0
        


</td></tr>
<tr>
    <td><code>string</code></td>
    <td width="100%">sensitive_summary</td><td>Yes</td>
</tr>
<tr><td colspan="3">

Sensitive summary value for output handling tests.

    

    

    

    

    


</td></tr>
<tr>
    <td><code>list(number)</code></td>
    <td width="100%">service_ports</td><td></td>
</tr>
<tr><td colspan="3">

Extracted list of service ports.

    

    

    

    

    


</td></tr>
</tbody></table>



## Objects



#### autoscaling_profiles

Nested map(object(...)) for autoscaling profiles.

    

    

    

    

    
<table><thead><tr><th>Type</th><th align="left" width="100%">Name</th><th>Required</th><th>Default&nbsp;Value</th></tr></thead><tbody>
        <tr>
    <td><code>number</code></td>
    <td width="100%">min</td><td>Yes</td><td><em>n/a</em></td>
</tr>
<tr><td colspan="4">



    

    

    

    

    


</td></tr>
<tr>
    <td><code>number</code></td>
    <td width="100%">max</td><td>Yes</td><td><em>n/a</em></td>
</tr>
<tr><td colspan="4">



    

    

    

    

    


</td></tr>
<tr>
    <td><code>number</code></td>
    <td width="100%">target_cpu</td><td></td><td><code>70</code></td>
</tr>
<tr><td colspan="4">



    

    

    

    

    


</td></tr>
</tbody></table>



#### code_fence_examples



    

    

    

    

    
<table><thead><tr><th>Type</th><th align="left" width="100%">Name</th><th>Required</th><th>Default&nbsp;Value</th></tr></thead><tbody>
        <tr>
    <td><code>bool</code></td>
    <td width="100%">enabled</td><td>Yes</td><td><em>n/a</em></td>
</tr>
<tr><td colspan="4">



    

    

    

    

    


</td></tr>
<tr>
    <td><code>list(string)</code></td>
    <td width="100%">labels</td><td>Yes</td><td><em>n/a</em></td>
</tr>
<tr><td colspan="4">



    

    

    

    

    


</td></tr>
</tbody></table>



#### deeply_nested_object



    

    

    

    

    
<table><thead><tr><th>Type</th><th align="left" width="100%">Name</th><th>Required</th><th>Default&nbsp;Value</th></tr></thead><tbody>
        <tr>
    <td><code>string</code></td>
    <td width="100%">name</td><td>Yes</td><td><em>n/a</em></td>
</tr>
<tr><td colspan="4">

Required top-level name.

    

    

    

    

    


</td></tr>
<tr>
    <td><code>bool</code></td>
    <td width="100%">enabled</td><td></td><td><code>true</code></td>
</tr>
<tr><td colspan="4">

Optional top-level toggle.

    

    

    

    

    


</td></tr>
<tr>
    <td><code>object(<a href="#deployment">deployment</a>)</code></td>
    <td width="100%">deployment</td><td></td><td><code>null</code></td>
</tr>
<tr><td colspan="4">

Optional nested object with null default semantics.

    

    

    

    

    


</td></tr>
</tbody></table>



#### deployment

Optional nested object with null default semantics.

    

    

    

    

    
<table><thead><tr><th>Type</th><th align="left" width="100%">Name</th><th>Required</th><th>Default&nbsp;Value</th></tr></thead><tbody>
        <tr>
    <td><code>string</code></td>
    <td width="100%">region</td><td>Yes</td><td><em>n/a</em></td>
</tr>
<tr><td colspan="4">

Required region for deployment.

    

    

    

    

    


</td></tr>
<tr>
    <td><code>number</code></td>
    <td width="100%">replicas</td><td></td><td><code>2</code></td>
</tr>
<tr><td colspan="4">

Optional replica count.

    

    

    

    

    


</td></tr>
<tr>
    <td><code>list(string)</code></td>
    <td width="100%">zones</td><td></td><td><code>[]</code></td>
</tr>
<tr><td colspan="4">

Optional list of zones.

    

    

    

    

    


</td></tr>
<tr>
    <td><code>map(object(<a href="#autoscaling_profiles">autoscaling_profiles</a>))</code></td>
    <td width="100%">autoscaling_profiles</td><td></td><td><code>{}</code></td>
</tr>
<tr><td colspan="4">

Nested map(object(...)) for autoscaling profiles.

    

    

    

    

    


</td></tr>
</tbody></table>



#### details



    

    

    

    

    
<table><thead><tr><th>Type</th><th align="left" width="100%">Name</th><th>Required</th><th>Default&nbsp;Value</th></tr></thead><tbody>
        <tr>
    <td><code>string</code></td>
    <td width="100%">title</td><td>Yes</td><td><em>n/a</em></td>
</tr>
<tr><td colspan="4">



    

    

    

    

    


</td></tr>
<tr>
    <td><code>string</code></td>
    <td width="100%">body</td><td>Yes</td><td><em>n/a</em></td>
</tr>
<tr><td colspan="4">



    

    

    

    

    


</td></tr>
</tbody></table>



#### documented_map



    

    

    

    

    
<table><thead><tr><th>Type</th><th align="left" width="100%">Name</th><th>Required</th><th>Default&nbsp;Value</th></tr></thead><tbody>
        <tr>
    <td><code>string</code></td>
    <td width="100%">display_name</td><td>Yes</td><td><em>n/a</em></td>
</tr>
<tr><td colspan="4">

Human-friendly display name.

    

    

    

    

    


</td></tr>
<tr>
    <td><code>string</code></td>
    <td width="100%">tier</td><td>Yes</td><td><em>n/a</em></td>
</tr>
<tr><td colspan="4">

Environment tier.

    

    

    

    

    


</td></tr>
<tr>
    <td><code>list(string)</code></td>
    <td width="100%">tags</td><td></td><td><code>[]</code></td>
</tr>
<tr><td colspan="4">

Optional tags for each object.

    

    

    

    

    


</td></tr>
<tr>
    <td><code>map(string)</code></td>
    <td width="100%">metadata</td><td></td><td><code>{}</code></td>
</tr>
<tr><td colspan="4">

Optional metadata map.

    

    

    

    

    


</td></tr>
</tbody></table>



#### health_check



    

    

    

    

    
<table><thead><tr><th>Type</th><th align="left" width="100%">Name</th><th>Required</th><th>Default&nbsp;Value</th></tr></thead><tbody>
        <tr>
    <td><code>string</code></td>
    <td width="100%">path</td><td></td><td><code>"/healthz"</code></td>
</tr>
<tr><td colspan="4">



    

    

    

    

    


</td></tr>
<tr>
    <td><code>number</code></td>
    <td width="100%">interval_seconds</td><td></td><td><code>30</code></td>
</tr>
<tr><td colspan="4">



    

    

    

    

    


</td></tr>
<tr>
    <td><code>number</code></td>
    <td width="100%">timeout_seconds</td><td></td><td><code>5</code></td>
</tr>
<tr><td colspan="4">



    

    

    

    

    


</td></tr>
<tr>
    <td><code>list(number)</code></td>
    <td width="100%">expected_statuses</td><td></td><td><code>[200]</code></td>
</tr>
<tr><td colspan="4">



    

    

    

    

    


</td></tr>
<tr>
    <td><code>bool</code></td>
    <td width="100%">enable_status_check</td><td></td><td><code>true</code></td>
</tr>
<tr><td colspan="4">



    

    

    

    

    


</td></tr>
</tbody></table>



#### list_of_objects_mixed_docs



    

    

    

    

    
<table><thead><tr><th>Type</th><th align="left" width="100%">Name</th><th>Required</th><th>Default&nbsp;Value</th></tr></thead><tbody>
        <tr>
    <td><code>string</code></td>
    <td width="100%">name</td><td>Yes</td><td><em>n/a</em></td>
</tr>
<tr><td colspan="4">

Name of the service.

    

    

    

    

    


</td></tr>
<tr>
    <td><code>number</code></td>
    <td width="100%">port</td><td></td><td><code>8080</code></td>
</tr>
<tr><td colspan="4">

Public port for the service.

    

    

    

    

    


</td></tr>
<tr>
    <td><code>string</code></td>
    <td width="100%">protocol</td><td></td><td><code>"http"</code></td>
</tr>
<tr><td colspan="4">



    

    

    

    

    


</td></tr>
<tr>
    <td><code>object(<a href="#health_check">health_check</a>)</code></td>
    <td width="100%">health_check</td><td></td><td><code>null</code></td>
</tr>
<tr><td colspan="4">



    

    

    

    

    


</td></tr>
</tbody></table>



#### map_of_objects_documented



    

    

    

    

    
<table><thead><tr><th>Type</th><th align="left" width="100%">Name</th><th>Required</th><th>Default&nbsp;Value</th></tr></thead><tbody>
        <tr>
    <td><code>string</code></td>
    <td width="100%">display_name</td><td>Yes</td><td><em>n/a</em></td>
</tr>
<tr><td colspan="4">

Human-friendly display name.

    

    

    

    

    
**Since:** 1.0.0
        


</td></tr>
<tr>
    <td><code>string</code></td>
    <td width="100%">tier</td><td>Yes</td><td><em>n/a</em></td>
</tr>
<tr><td colspan="4">

Environment tier.

    
**Allowed Values:**
- `dev`
- `stage`
- `prod`

    

    

    

    


</td></tr>
<tr>
    <td><code>list(string)</code></td>
    <td width="100%">tags</td><td></td><td><code>[]</code></td>
</tr>
<tr><td colspan="4">

Optional tags for each object.

    

    

    

    

    


</td></tr>
<tr>
    <td><code>map(string)</code></td>
    <td width="100%">metadata</td><td></td><td><code>{}</code></td>
</tr>
<tr><td colspan="4">

Optional metadata map.

    

    

    

    

    


</td></tr>
</tbody></table>



#### object_with_block_comment_docs



    

    

    

    

    
<table><thead><tr><th>Type</th><th align="left" width="100%">Name</th><th>Required</th><th>Default&nbsp;Value</th></tr></thead><tbody>
        <tr>
    <td><code>string</code></td>
    <td width="100%">id</td><td>Yes</td><td><em>n/a</em></td>
</tr>
<tr><td colspan="4">

Primary key.

    

    

    

    

    
**Since:** 1.0.0
        


</td></tr>
<tr>
    <td><code>string</code></td>
    <td width="100%">summary</td><td></td><td><code>null</code></td>
</tr>
<tr><td colspan="4">

Optional summary.

    

    

    

    

    
**Deprecated:** Prefer `details`.
        


</td></tr>
<tr>
    <td><code>object(<a href="#details">details</a>)</code></td>
    <td width="100%">details</td><td></td><td><code>null</code></td>
</tr>
<tr><td colspan="4">



    

    

    

    

    


</td></tr>
</tbody></table>





[terraform-functions]: https://developer.hashicorp.com/terraform/language/functions


<!-- TFDOCS_EXTRAS_END -->
