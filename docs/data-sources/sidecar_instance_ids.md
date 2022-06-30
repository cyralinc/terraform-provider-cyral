---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cyral_sidecar_instance_ids Data Source - cyral"
subcategory: ""
description: |-
  Retrieves the IDs of all the current instances of a given sidecar.
---

# cyral_sidecar_instance_ids (Data Source)

Retrieves the IDs of all the current instances of a given sidecar.

## Example Usage

```terraform
data "cyral_sidecar_instance_ids" "this" {
  sidecar_id = cyral_sidecar.some_sidecar_resource.id
}

output "sidecar_instance_ids" {
  value = data.cyral_sidecar_instance_ids.this.instance_ids
}
```

<!-- schema generated by tfplugindocs -->

## Schema

### Required

- `sidecar_id` (String) The ID of the sidecar.

### Read-Only

- `id` (String) Computed ID for this data source (locally computed to be used in Terraform state).
- `instance_ids` (List of String) All the current instance IDs of the sidecar.