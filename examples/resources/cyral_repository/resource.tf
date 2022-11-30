### Minimal Repository
resource "cyral_repository" "some_resource_name" {
    type = "mongodb"
    name = "some_repo_name"

    repo_node {
        name = "node-1"
        host = "mongodb.cyral.com"
        port = 27017
    }
}

### Repository with Connection Draining, Preferred Access Gateway, and Labels
resource "cyral_repository" "some_resource_name" {
    type = "mongodb"
    name = "some_repo_name"
    labels = [ "single-node", "us-east-1" ]

    repo_node {
        name = "node-1"
        host = "mongodb.cyral.com"
        port = 0
    }

    connection_draining {
      auto = true
      wait_time = 30
    }

    preferred_access_gateway {
      sidecar_id = "some-sidecar-id"
      binding_id = "some-binding-id"
    }
}

### Multi-Node MongoDB Repository with Replicaset
resource "cyral_repository" "some_resource_name" {
    type = "mongodb"
    name = "some_repo_name"
    labels = [ "multi-node", "us-east-2" ]

    repo_node {
        name = "node-1"
        host = "mongodb-node1.cyral.com"
        port = 27017
    }

    repo_node {
        name = "node-2"
        host = "mongodb-node2.cyral.com"
        port = 27017
    }

    repo_node {
        name = "node-3"
        dynamic = true
    }

    mongodb_settings {
      replica_set_name = "some-replica-set"
      server_type = "replicaset"
    }
}
