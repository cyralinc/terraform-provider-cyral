### Minimal Repository
resource "cyral_repository" "minimal_repo" {
    type = "mongodb"
    name = "minimal_repo"

    repo_node {
        name = "node-1"
        host = "mongodb.cyral.com"
        port = 27017
    }
}

### Repository with Connection Draining, Preferred Access Gateway, and Labels
resource "cyral_repository" "repo_with_conn_draining" {
    type = "mongodb"
    name = "repo_with_conn_draining"
    labels = [ "single-node", "us-east-1" ]

    repo_node {
        name = "node-1"
        host = "mongodb.cyral.com"
        port = 27017
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
resource "cyral_repository" "multi_node_mongo_repo" {
    type = "mongodb"
    name = "multi_node_mongo_repo"
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
