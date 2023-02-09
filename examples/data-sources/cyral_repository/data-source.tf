resource "cyral_repository" "mongo-repository" {
  type = "mongodb"
  name = "tf-provider-mongo-repository"

  repo_node {
    host = "mongodb.cyral.com"
    port = 27017
  }
  mongodb_settings {
    server_type = "standalone"
  }
}

resource "cyral_repository" "mysql-repository" {
  type = "mysql"
  name = "tf-provider-mysql-repository"

  repo_node {
    host = "mysql.com"
    port = 3306
  }
}

data "cyral_repository" "search-for-mysql-repo" {
  depends_on = [
    cyral_repository.mongo-repository,
    cyral_repository.mysql-repository
  ]
  name = "tf-provider-mysql-repository"
  type = "mysql"
}

output "mysql_repo_id" {
  value = data.cyral_repository.search-for-mysql-repo.id
}
