resource "cyral_repository_binding" "some_resource_name" {
    enabled = true
    repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
    sidecar_id    = cyral_sidecar.SOME_SIDECAR_RESOURCE_NAME.id
    listener_port = 0
    listener_host = "0.0.0.0"
    sidecar_as_idp_access_gateway = false
}
