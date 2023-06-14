resource "cyral_sidecar" "sidecar" {
  name               = "my-sidecar"
  deployment_method  = "terraform"
}

resource "cyral_repository" "mysql_1" {
  name     = "mysql-1"
  type     = "mysql"

  repo_node {
    host = "your-mysql-1-db-host"
    // This database accepts incoming connections in 3307,
    // but all of them could be using the same port as long
    // as they live in different hosts.
    port = 3307
  }
}

resource "cyral_repository" "mysql_2" {
  name     = "mysql-2"
  type     = "mysql"

  repo_node {
    host = "your-mysql-2-db-host"
    // This database accepts incoming connections in 3309,
    // but all of them could be using the same port as long
    // as they live in different hosts.
    port = 3309
  }
}

resource "cyral_sidecar_listener" "listener" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["mysql"]
  // Clients will connect to both MySQL repos through
  // the sidecar using 3306
  network_address {
    port = 3306
  }
  # MySQL version that will be shown to clients
  # connecting to both MySQL instances
  mysql_settings {
    db_version    = "8.0.4"
  }
}

resource "cyral_repository_binding" "binding_mysql_1" {
  repository_id = cyral_repository.mysql_1.id
  sidecar_id    = cyral_sidecar.sidecar.id

  listener_binding {
    listener_id = cyral_sidecar_listener.listener.listener_id
  }
}

resource "cyral_repository_binding" "binding_mysql_2" {
  repository_id = cyral_repository.mysql_2.id
  sidecar_id    = cyral_sidecar.sidecar.id

  listener_binding {
    listener_id = cyral_sidecar_listener.listener.listener_id
  }
}
