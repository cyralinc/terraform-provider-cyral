resource "cyral_sidecar" "sidecar" {
  name = "my-sidecar"
  deployment_method = "terraform"
}

resource "cyral_repository" "repo_mongodb" {
  name = "my-mongodb-repo"
  type = "mongodb"
  repo_node {
    name = "single-node-mongo"
    host = "mongodb.cyral.com"
    port = 27017
  }
  mongodb_settings {
    server_type = "standalone"
  }
}

resource "cyral_repository" "repo_pg" {
  name = "my-pg-repo"
  type = "postgresql"
  repo_node {
    host = "pg.cyral.com"
    port = 5432
  }
}

resource "cyral_sidecar_listener" "listener_mongodb" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["mongodb"]
  network_address {
    port = 27017
  }
}

resource "cyral_sidecar_listener" "listener_pg" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["postgresql"]
  network_address {
    port = 5434 // clients will connect to the pg repo using sidecar port 5434
  }
}

resource "cyral_repository_binding" "binding_mongodb" {
  sidecar_id = cyral_sidecar.sidecar.id
  repository_id = cyral_repository.repo_mongodb.id
  listener_binding {
    listener_id = cyral_sidecar_listener.listener_mongodb.listener_id
    node_index = 0 // optional if there is only one repo node
  }
}

resource "cyral_repository_binding" "binding_pg" {
  sidecar_id = cyral_sidecar.sidecar.id
  repository_id = cyral_repository.repo_pg.id
  listener_binding {
    listener_id = cyral_sidecar_listener.listener_pg.listener_id
  }
}
