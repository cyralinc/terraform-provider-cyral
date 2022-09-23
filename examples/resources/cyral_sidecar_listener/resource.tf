
### plain mySQL listener
resource "cyral_sidecar_listener" "plain_mysql" {
  sidecar_id = "2F1rBhVT7nX3GCzXGEWOHGcmEzP"
  tcp_listener {
    port = 3306
  }
  repo_types =["mysql"]
}

### multiplexed mySQL listener
resource "cyral_sidecar_listener" "multiplex_mysql" {
  sidecar_id = "2F1rBhVT7nX3GCzXGEWOHGcmEzP"
  tcp_listener {
    port = 3307
  }
  multiplexed = true
  mysql_settings {
    db_version = "5.7"
  }
  repo_types =["mysql"]
}


### S3 listener, using proxy mode
resource "cyral_sidecar_listener" "s3_proxy" {
  sidecar_id = "2F1rBhVT7nX3GCzXGEWOHGcmEzP"
  tcp_listener {
    port = 443
  }
  s3_settings {
    proxy_mode = true
  }
  repo_types =["s3"]
}

### mariaDB using unix socket listener
resource "cyral_sidecar_listener" "file_mariadb" {
  sidecar_id = "2F1rBhVT7nX3GCzXGEWOHGcmEzP"
  unix_listener {
    file = "/var/run/mysqld/mysql.sock"
  }
  repo_types =["mariadb"]
}
