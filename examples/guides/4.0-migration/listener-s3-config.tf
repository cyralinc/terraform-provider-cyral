resource "cyral_sidecar_listener" "sidecar_all_listeners" {
  for_each = local.type_port_map
  repo_types = [each.value.type]
  sidecar_id = cyral_sidecar.sidecar.id

  network_address {
    host          = "0.0.0.0"
    port          = each.value.port
  }

  dynamic "s3_settings" {
    for_each = each.value.type == "s3" ? [""] : []
    content {
      proxy_mode = true
    }
  }
}
