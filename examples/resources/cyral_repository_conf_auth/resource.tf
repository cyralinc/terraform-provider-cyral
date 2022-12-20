resource "cyral_repository_conf_auth" "some_resource_name" {
    repository_id = ""
    allow_native_auth = true
    client_tls = "enable|disable|enabledAndVerifyCertificate"
    identity_provider = ""
    repo_tls = "enable|disable|enabledAndVerifyCertificate"
}
