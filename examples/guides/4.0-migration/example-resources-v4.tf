resource "cyral_repository" "mongo_repo" {
  name = "mongo-repo"
  type = "mongodb"
  repo_node {
    host = "mongodb.cyral.com"
    port = 27017
  }
  mongodb_settings {
    server_type = "standalone"
  }
}

resource "cyral_repository_binding" "binding" {
  sidecar_id  = cyral_sidecar.sidecar.id
  repository_id = cyral_repository.mongo_repo.id
  enabled   = true
  listener_binding {
    listener_id = cyral_sidecar_listener.listener.listener_id
    node_index  = 0
  }
}

resource "cyral_sidecar_listener" "listener" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["mongodb"]
  network_address {
    port = 27020
  }
}

resource "cyral_repository_access_gateway" "access_gateway" {
  repository_id = cyral_repository.mongo_repo.id
  sidecar_id  = cyral_sidecar.sidecar.id
  binding_id  = cyral_repository_binding.binding.binding_id
}
