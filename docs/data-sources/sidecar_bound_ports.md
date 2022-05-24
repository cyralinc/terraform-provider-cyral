# Sidecar Bound Ports Data Source

Retrieves all the ports of a given sidecar that are currently bound to repositories.

## Example Usage

```hcl
resource "cyral_sidecar" "sidecar_1" {
  name = "tf-provider-sidecar-1"
  deployment_method = "cloudFormation"
}

resource "cyral_repository" "repo_1" {
  name = "tf-provider-repo-1"
  type = "mysql"
  host = "mysql.com"
  port = 3306
}

resource "cyral_repository_binding" "repo_binding_1" {
  repository_id = cyral_repository.repo_1.id
  sidecar_id = cyral_sidecar.sidecar_1.id
  listener_port = 3306
  enabled = true
}

resource "cyral_repository" "repo_2" {
  name = "tf-provider-repo-2"
  type = "mongodb"
  host = "mongodb.com"
  port = 27017
}

resource "cyral_repository_binding" "repo_binding_2" {
  repository_id = cyral_repository.repo_2.id
  sidecar_id = cyral_sidecar.sidecar_1.id
  listener_port = 27017
  enabled = true
}

data "cyral_sidecar_bound_ports" "this" {
  depends_on = [
    cyral_repository_binding.repo_binding_1,
    cyral_repository_binding.repo_binding_2
  ]
  sidecar_id = cyral_sidecar.sidecar_1.id
}

output "sidecar_bound_ports" {
  value = data.cyral_sidecar_bound_ports.this.bound_ports
}
```

-> Notice that, in this case, the `depends_on` argument will be needed if you want to retrieve the sidecar bound ports only after the bindings are created/updated. Otherwise, if `depends_on` is omitted, the data source will retrieve the bound ports before creating/updating the bindings, which in this case would return zero ports.

## Argument Reference

- `sidecar_id` - (Required) The ID of the sidecar.

## Attribute Reference

- `id` - Computed ID for this data source (locally computed to be used in Terraform state).
- `bound_ports` - All the sidecar ports that are currently bound to repositories.
