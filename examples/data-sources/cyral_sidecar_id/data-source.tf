resource "cyral_sidecar" "sidecar_1" {
  name = "tf-provider-sidecar-1"
  deployment_method = "cloudFormation"
}

data "cyral_sidecar_id" "this" {
  sidecar_name = cyral_sidecar.sidecar_1.name
}

output "sidecar_id" {
  value = data.cyral_sidecar_id.this.id
}
