resource "cyral_sidecar" "sidecar" {
  name              = "sidecar"
  deployment_method = "docker"
}

// Plain listener
resource "cyral_sidecar_listener" "listener" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["mongodb"]
  network_address {
    port = 27017
  }
}

// Listener with MySQL Settings
resource "cyral_sidecar_listener" "listener_mysql" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["mysql"]
  network_address {
    port = 3306
  }

  mysql_settings {
    db_version    = "8.0.4"
    character_set = "utf8mb4_0900_ai_ci"
  }
}

# Listener for S3 CLI and AWS SDK
resource "cyral_sidecar_listener" "listener_s3_cli" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["s3"]
  network_address {
    port = 443
  }
  s3_settings {
    proxy_mode = true
  }
}

# Listener for S3 browser (using port 444 assuming port
# 443 is used for CLI)
resource "cyral_sidecar_listener" "listener_s3_cli" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["s3"]
  network_address {
    port = 444
  }
  s3_settings {
    // may be omitted for s3 browser as it defaults to `false`
    proxy_mode = false
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
    // must be true if repo_type is either `dynamodb` or `dynamodbstreams`
    proxy_mode = true
  }
}
