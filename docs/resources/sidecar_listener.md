# cyral_sidecar_listener (Resource)

Create new [sidecar listeners](https://cyral.com/docs/sidecars/sidecar-listeners).

-> **NOTE** Import ID syntax is `{sidecar_id}-{listener_id}`.

## Example Usage

```terraform
### plain mySQL listener
resource "cyral_sidecar_listener" "plain_mysql" {
    sidecar_id = cyral_sidecar.SOME_SIDECAR_RESOURCE_NAME.id
    tcp_listener_port = 3306
    repo_types =["mysql"]
}

### multiplexed mySQL listener
resource "cyral_sidecar_listener" "multiplex_mysql" {
    sidecar_id = cyral_sidecar.SOME_SIDECAR_RESOURCE_NAME.id
    tcp_listener_port = 3307
    multiplexed = true
    mysql_settings_db_version = "5.7"
    repo_types =["mysql"]
}

### S3 listener, using proxy mode
resource "cyral_sidecar_listener" "s3_proxy" {
    sidecar_id = cyral_sidecar.SOME_SIDECAR_RESOURCE_NAME.id
    tcp_listener_port = 443
    s3_settings_proxy_mode = true
    repo_types =["s3"]
}

### mariaDB using unix socket istener
resource "cyral_sidecar_listener" "file_mariadb" {
    sidecar_id = cyral_sidecar.SOME_SIDECAR_RESOURCE_NAME.id
    unix_listener_file = "/var/run/mysqld/mysql.sock"
    repo_types =["mariadb"]
}
```
