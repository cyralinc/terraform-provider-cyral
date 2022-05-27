data "cyral_sidecar_instance_ids" "this" {
  sidecar_id = cyral_sidecar.some_sidecar_resource.id
}

output "sidecar_instance_ids" {
  value = data.cyral_sidecar_instance_ids.this.instance_ids
}
