resource "cyral_repository" "repo" {
  name = "tf-account-repo"
  type          = "mongodb"
  repo_node {
        name = "single-node-mongo"
        host = "mongodb.cyral.com"
        port = 27017
  }
}

resource "cyral_sidecar" "sidecar" {
  name = "tf-account-sidecar"
  deployment_method = "docker"
}

resource "cyral_sidecar_listener" "listener" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["mongodb"]
  network_address {
    host          = "mongodb.cyral.com"
    port          = 27017
  }
}

resource "cyral_repository_binding" "binding" {
  sidecar_id = cyral_sidecar.sidecar.id
  repository_id = cyral_repository.repo.id
  enabled = true
  listener_binding {
    listener_id = cyral_sidecar_listener.listener.listener_id
    node_index = 0
  }
}
