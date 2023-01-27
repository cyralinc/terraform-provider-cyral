resource "cyral_sidecar" "sidecar" {
  name              = "tf-account-sidecar"
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
resource "cyral_sidecar_listener" "listener" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["mysql"]
  network_address {
    # Specify the network address if the sidecar software is deployed on a host with multiple network
    # interfaces but we want the sidecar to only accept connections on a specific one. Otherwise, if you
    # leave host empty, the sidecar will listen on all network interfaces.
    host = "0.0.0.0"
    port = 3306
  }

  mysql_settings {
    db_version    = "3.4.0"
    character_set = "ujis_japanese_ci"
  }
}

# Listener with S3 Settings
resource "cyral_sidecar_listener" "listener" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["s3"]
  network_address {
    port = 443
  }
  s3_settings {
    proxy_mode = true
  }
}

# Listener with DynamoDB Settings
resource "cyral_sidecar_listener" "listener" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["dynamodb"]
  network_address {
    port = 8000
  }
  dynamodb_settings {
    proxy_mode = true
  }
}
