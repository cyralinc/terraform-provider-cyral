terraform {
  required_providers {
    cyral = {
      source  = "cyralinc/cyral"
      version = "~> 3.0"
    }
  }
}

provider "cyral" {
    # Follow the instructions in the Cyral Terraform Provider page to set up the
    # credentials: https://registry.terraform.io/providers/cyralinc/cyral/latest/docs
    client_id     = "sa/default/af295f09-3c48-409e-aada-79d04eb200af"
    client_secret = "2dcda1e3-b802-45ef-989d-97a42c280a99"

    control_plane = "hbf1812-tfmigrate-a01-ctl.k8-sandbox.gcp.cyral.com/"

}

resource "cyral_sidecar" "pg_sidecar" {
  name              = "MainSidecar"
  deployment_method = "terraform"
}

resource "cyral_sidecar_credentials" "sidecar_credentials" {
    sidecar_id = cyral_sidecar.pg_sidecar.id
}

resource "cyral_repository" "pg_repo" {
  name = "pg_repo"
  type = "postgresql"
  host = "postgresql.mycompany.com"
  port = 5432
}

resource "cyral_repository_binding" "pg_repo_binding" {
  repository_id = cyral_repository.pg_repo.id
  sidecar_id    = cyral_sidecar.pg_sidecar.id
  listener_port = 5432
}

resource "cyral_repository" "mongodb_repo" {
  name = "mongodb_repo"
  type = "mongodb"

  # Specify the address or hostname of the endpoint of one node in the
  # MongoDB replica set. Cyral will automatically/dynamically identify
  # the remaining nodes of the replication cluster.
  host = "mycluster-shard-00-01.example.mongodb.net"

  port = 27017
  properties {
    mongodb_replica_set {
      max_nodes = 3

      # Specify the replica set identifier, a string value that
      # identifies the MongoDB replica set cluster. To find your
      # replica set ID, see our article:
      #
      # * https://cyral.freshdesk.com/a/solutions/articles/44002241594
      replica_set_id = "my-replica-set-id"
    }
  }
}

resource "cyral_repository_binding" "mongodb_repo_binding" {
  repository_id                 = cyral_repository.mongodb_repo.id
  sidecar_id                    = cyral_sidecar.mongodb_sidecar.id
  listener_port                 = 27017
  sidecar_as_idp_access_gateway = true
}

resource "cyral_sidecar" "mongodb_sidecar" {
  name              = "MongoDBSidecar"
  deployment_method = "terraform"
}
