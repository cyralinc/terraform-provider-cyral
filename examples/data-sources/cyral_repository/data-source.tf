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

resource "cyral_repository" "mysql-repository1" {
  type = "mysql"
  name = "tf-provider-mysql-repository1"

  repo_node {
    host = "mysql.com"
    port = 3306
  }
}

resource "cyral_repository" "mysql-repository2" {
  type = "mysql"
  name = "tf-provider-mysql-repository2"

  repo_node {
    host = "mysql2.com"
    port = 3306
  }
}

data "cyral_repository" "specific-mysql-repo" {
  depends_on = [
    cyral_repository.mongo-repository,
    cyral_repository.mysql-repository1,
    cyral_repository.mysql-repository2,
  ]
  # As we have more than one MySQL repos, we need to provide
  # the name and type or just the name (repo names are unique)
  name = "tf-provider-mysql-repository1"
  type = "mysql"
}

data "cyral_repository" "all-mysql-repos" {
  depends_on = [
    cyral_repository.mongo-repository,
    cyral_repository.mysql-repository1,
    cyral_repository.mysql-repository2,
  ]
  type = "mysql"
}

output "mysql1_repo_id" {
  # Because our search is targeting a specific name that we
  # know it exists referencing index 0 is safe.
  value = data.cyral_repository.specific-mysql-repo.repository_list[0].id
}

output "all_mysql_repo_ids" {
  value = data.cyral_repository.all-mysql-repos.repository_list
}
