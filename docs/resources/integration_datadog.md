---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cyral_integration_datadog Resource - terraform-provider-cyral"
subcategory: ""
description: |-
    ~> DEPRECATED If configuring Datadog for logging purposes, use resource cyral_integration_logging instead.
---

# cyral_integration_datadog (Resource)

~> **DEPRECATED** If configuring Datadog for logging purposes, use resource `cyral_integration_logging` instead.

## Example Usage

```terraform
resource "cyral_integration_datadog" "some_resource_name" {
    name = ""
    api_key = ""
}
```

<!-- schema generated by tfplugindocs -->

## Schema

### Required

-   `api_key` (String, Sensitive) Datadog API key.
-   `name` (String) Integration name that will be used internally in the control plane.

### Read-Only

-   `id` (String) ID of this resource in Cyral environment
