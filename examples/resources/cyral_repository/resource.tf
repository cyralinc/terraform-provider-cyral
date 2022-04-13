### Single repository
resource "cyral_repository" "some_resource_name" {
    host = ""
    port = 0
    type = ""
    name = ""
}

### Multiple repositories using a local variable
locals {
    repos = {
        mymongodb = {
            host = "mongodb.cyral.com"
            port = 27017
            type = "mongodb"
        }
        mymariadb = {
            host = "mariadb.cyral.com"
            port = 3310
            type = "mariadb"
        }
        mypostgresql = {
            host = "postgresql.cyral.com"
            port = 5432
            type = "postgresql"
        }
    }
}

resource "cyral_repository" "repositories" {
    for_each = local.repos
    name  = each.key
    type  = each.value.type
    host  = each.value.host
    port  = each.value.port
}