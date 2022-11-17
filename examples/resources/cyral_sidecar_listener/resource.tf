### Plain mySQL listener
resource "cyral_sidecar_listener" "plain_mysql" {
  sidecar_id = "some-sidecar-id"
  network_address {
    port = 3306
  }
  repo_types =["mysql"]
}

### MySQL listener with Settings
resource "cyral_sidecar_listener" "mysql_with_settings" {
  sidecar_id = "some-sidecar-id"
  network_address {
    port = 3307
    host = "some.mysqldb.com"
  }
  mysql_settings {
    db_version = "5.7.0"
    character_set = "latin1_german1_ci"
  }
  repo_types =["mysql"]
}


### S3 listener with Proxy Mode
resource "cyral_sidecar_listener" "s3_proxy" {
  sidecar_id = "some-sidecar-id"
  network_address {
    port = 443
  }
  s3_settings {
    proxy_mode = true
  }
  repo_types =["s3"]
}

### DynamoDB listener with Proxy Mode
resource "cyral_sidecar_listener" "dynamodb_proxy" {
  sidecar_id = "some-sidecar-id"
  network_address {
    port = 8000
  }
  dynamodb_settings {
    proxy_mode = true
  }
  repo_types =["dynamodb"]
}
