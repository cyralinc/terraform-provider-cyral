resource "cyral_repository_binding" "all_repo_binding" {
  for_each = local.repos
  enabled       = true
  repository_id = cyral_repository.all_repositories[each.key].id
  sidecar_id    = cyral_sidecar.sidecar.id

  dynamic "listener_binding" {
    for_each = each.value.ports
    content {
      listener_id = cyral_sidecar_listener.sidecar_all_listeners["${each.value.type}_${listener_binding.value}"].listener_id
      node_index = listener_binding.key
    }
  }
}
