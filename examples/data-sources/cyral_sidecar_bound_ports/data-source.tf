resource "cyral_sidecar" "sidecar" {
  name = "tf-account-sidecar"
  deployment_method = "docker"
}

resource "cyral_repository" "repo-1" {
  name = "tf-account-repo-1"
  type          = "mongodb"
  repo_node {
        host = "mongodb.cyral.com"
        port = 27017
  }
}

resource "cyral_sidecar_listener" "listener-1" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["mongodb"]
  network_address {
    host          = "mongodb.cyral.com"
    port          = 27017
  }
}

resource "cyral_repository_binding" "binding-1" {
  sidecar_id = cyral_sidecar.sidecar.id
  repository_id = cyral_repository.repo-1.id
  enabled = true
  listener_binding {
    listener_id = cyral_sidecar_listener.listener-1.listener_id
    node_index = 0
  }
}

resource "cyral_repository" "repo-2" {
  name = "tf-account-repo-2"
  type          = "mongodb"
  repo_node {
        host = "mongodb.cyral.com"
        port = 27018
  }
}

resource "cyral_sidecar_listener" "listener-2" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["mongodb"]
  network_address {
    host          = "mongodb.cyral.com"
    port          = 27017
  }
}

resource "cyral_repository_binding" "binding-2" {
  sidecar_id = cyral_sidecar.sidecar.id
  repository_id = cyral_repository.repo-2.id
  enabled = true
  listener_binding {
    listener_id = cyral_sidecar_listener.listener-2.listener_id
    node_index = 0
  }
}

data "cyral_sidecar_bound_ports" "this" {
  # Notice that in this case the `depends_on` argument will be
  # needed if you want to retrieve the sidecar bound ports only
  # after the bindings are created/updated. Otherwise, if
  # `depends_on` is omitted, the data source will retrieve the
  # bound ports before creating/updating the bindings, which in
  # this case would return zero ports.
  depends_on = [
    cyral_repository_binding.binding_1,
    cyral_repository_binding.binding_2
  ]
  sidecar_id = cyral_sidecar.sidecar.id
}

output "sidecar_bound_ports" {
  value = data.cyral_sidecar_bound_ports.this.bound_ports
}
