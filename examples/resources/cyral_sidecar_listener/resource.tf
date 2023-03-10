resource "cyral_sidecar" "sidecar" {
  name              = "sidecar"
  deployment_method = "docker"
}

# Plain listener
resource "cyral_sidecar_listener" "listener" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["mongodb"]
  network_address {
    port = 27017
  }
}

# Listener with MySQL Settings
resource "cyral_sidecar_listener" "listener_mysql" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["mysql"]
  network_address {
    port = 3306
  }

  mysql_settings {
    db_version    = "3.4.0"
    character_set = "ujis_japanese_ci"
  }
}

# Listener for S3 CLI
resource "cyral_sidecar_listener" "listener_s3_cli" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["s3"]
  network_address {
    port = 443
  }
}

# Listener for S3 browser (using port 444 assuming port 443 is used for CLI)
resource "cyral_sidecar_listener" "listener_s3_cli" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["s3"]
  network_address {
    port = 444
  }
  s3_settings {
    proxy_mode = true
  }
}


# Listener with DynamoDB Settings
resource "cyral_sidecar_listener" "listener_dynamodb" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["dynamodb"]
  network_address {
    port = 8000
  }
  dynamodb_settings {
    proxy_mode = true
  }
}
