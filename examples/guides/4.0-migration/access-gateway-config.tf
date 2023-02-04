resource "cyral_repository_access_gateway" "all_repo_binding_all_access_gateways" {
  for_each = local.repos
  repository_id  = cyral_repository.all_repositories[each.key].id
  sidecar_id  = cyral_sidecar.sidecar.id
  binding_id = cyral_repository_binding.all_repo_binding[each.key].binding_id
}
