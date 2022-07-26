resource "cyral_sidecar" "some_resource_name" {
    name = ""
    deployment_method = "someValidMethod"
    labels = ["label1", "label2"]
    user_endpoint = ""
    bypass_mode = "failover"
}
