resource "cyral_sidecar" "sidecar" {
  name = "tf-account-sidecar"
  deployment_method = "docker"
}

# Plain listener
resource "cyral_sidecar_listener" "listener" {
    sidecar_id = cyral_sidecar.sidecar.id
    repo_types = ["mongodb"]
    network_address {
        host          = "mongodb.cyral.com"
        port          = 27017
    }
}

# Listener with MySQL Settings
resource "cyral_sidecar_listener" "listener" {
    sidecar_id = cyral_sidecar.sidecar.id
    repo_types = ["mysql"]
    network_address {
        host          = "mysql.cyral.com"
        port          = 443
    }

    mysql_settings {
        db_version = "3.4.0"
        character_set = "ujis_japanese_ci"
    }
}

# Listener with S3 Settings
resource "cyral_sidecar_listener" "listener" {
    sidecar_id = cyral_sidecar.sidecar.id
    repo_types = ["s3"]
    network_address {
        host          = "s3.cyral.com"
        port          = 8002
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
        host          = "dynamodb.cyral.com"
        port          = 1234
    }
    dynamodb_settings {
        proxy_mode = true
    }
}
