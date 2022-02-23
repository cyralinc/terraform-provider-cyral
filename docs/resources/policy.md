# Policy Resource

Provides a resource to handle [policies](https://cyral.com/docs/reference/policy). See also: [Policy Rule](./policy_rule.md)

## Example Usage

```hcl
resource "cyral_policy" "some_resource_name" {
  data = [""]
  description = ""
  enabled = true | false
  name = ""
  tags = [""]
}
```

## Argument Reference

- `data` - (Optional) List that specify which data fields a policy manages. Each field is represented by the LABEL you established for it in your data map. The actual location of that data (the names of fields, columns, or databases that hold it) is listed in the data map.
- `description` - (Optional) String that describes the policy (ex: `your_policy_description`).
- `enabled` - (Optional) Boolean that causes a policy to be enabled or disabled.
- `name` - (Required) Policy name that will be used internally in Control Plane (ex: `your_policy_name`).
- `properties` - (Optional) Policy properties requiring a `name` and a `description`.
- `tags` - (Optional) Tags that can be used to organize and/or classify your policies (ex: `[your_tag1, your_tag2]`).

For more information, see the [Policy Guide](https://cyral.com/docs/policy#policy).

## Attribute Reference

- `id` - The ID of this resource.
- `created` - Policy creation timestamp.
- `last_updated` - Last update timestamp.
- `version` - Incremental counter for every update on the policy.
- `type` - Policy type.
