# cyral_repository (Resource)

Manages [repositories](https://cyral.com/docs/how-to/track-repos/).

-> Import ID syntax is `{repository_id}`.

## Example Usage

More complex examples using `cyral_repository` resource are available in the `Guides` section:

-   [Create an AWS EC2 sidecar to protect PostgreSQL and MySQL databases](https://registry.terraform.io/providers/cyralinc/cyral/latest/docs/guides/setup_cp_and_deploy_sidecar)
-   [Setup SSO access to MongoDB cluster using Okta IdP](https://registry.terraform.io/providers/cyralinc/cyral/latest/docs/guides/mongodb_cluster_okta_idp)

```terraform
### Minimal Repository
resource "cyral_repository" "minimal_repo" {
    type = "mongodb"
    name = "minimal_repo"

    repo_node {
        name = "node-1"
        host = "mongodb.cyral.com"
        port = 27017
    }

    mongodb_settings {
      server_type = "standalone"
    }
}

### Repository with Connection Draining and Labels
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

    mongodb_settings {
      server_type = "standalone"
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
```

<!-- schema generated by tfplugindocs -->

## Schema

### Required

-   `name` (String) Repository name that will be used internally in the control plane (ex: `your_repo_name`).
-   `repo_node` (Block List, Min: 1) List of nodes for this repository. (see [below for nested schema](#nestedblock--repo_node))
-   `type` (String) Repository type. List of supported types:
    -   `denodo`
    -   `dremio`
    -   `dynamodb`
    -   `dynamodbstreams`
    -   `galera`
    -   `mariadb`
    -   `mongodb`
    -   `mysql`
    -   `oracle`
    -   `postgresql`
    -   `redshift`
    -   `s3`
    -   `snowflake`
    -   `sqlserver`

### Optional

-   `connection_draining` (Block Set, Max: 1) Parameters related to connection draining. (see [below for nested schema](#nestedblock--connection_draining))
-   `labels` (List of String) Labels enable you to categorize your repository.
-   `mongodb_settings` (Block Set, Max: 1) Parameters related to MongoDB repositories. (see [below for nested schema](#nestedblock--mongodb_settings))
-   `redshift_settings` (Block Set, Max: 1) Parameters related to Redshift repositories. (see [below for nested schema](#nestedblock--redshift_settings))

### Read-Only

-   `id` (String) ID of this resource in Cyral environment.

<a id="nestedblock--repo_node"></a>

### Nested Schema for `repo_node`

Optional:

-   `dynamic` (Boolean) _Only supported for MongoDB in cluster configurations._
    Indicates if the node is dynamically discovered, meaning that the sidecar will query the cluster to get the topology information and discover the addresses of the dynamic nodes. If set to `true`, `host` and `port` must be empty. A node with value of this field as false considered `static`.
    The following conditions apply:
    -   The total number of declared `repo_node` blocks must match the actual number of nodes in the cluster.
    -   If there are static nodes in the configuration, they must be declared before all dynamic nodes.
    -   See the MongoDB-specific configuration in the [mongodb_settings](#nested-schema-for-mongodb_settings).
-   `host` (String) Repo node host (ex: `somerepo.cyral.com`). Can be empty if node is dynamic.
-   `name` (String) Name of the repo node.
-   `port` (Number) Repository access port (ex: `3306`). Can be empty if node is dynamic.

<a id="nestedblock--connection_draining"></a>

### Nested Schema for `connection_draining`

Optional:

-   `auto` (Boolean) Whether connections should be drained automatically after a listener dies.
-   `wait_time` (Number) Seconds to wait to let connections drain before starting to kill all the connections, if auto is set to true.

<a id="nestedblock--mongodb_settings"></a>

### Nested Schema for `mongodb_settings`

Required:

-   `server_type` (String) Type of the MongoDB server. Allowed values:

    -   `replicaset`
    -   `standalone`
    -   `sharded`

    The following conditions apply:

    -   If `sharded` and `srv_record_name` _not_ provided, then all `repo_node` blocks must be static (see [`dynamic`](#dynamic)).
    -   If `sharded` and `srv_record_name` provided, then all `repo_node` blocks must be dynamic (see [`dynamic`](#dynamic)).
    -   If `standalone`, then only one `repo_node` block can be declared and it must be static (see [`dynamic`](#dynamic)). The `srv_record_name` is not supported in this configuration.
    -   If `replicaset` and `srv_record_name` _not_ provided, then `repo_node` blocks may mix dynamic and static nodes (see [`dynamic`](#dynamic)).
    -   If `replicaset` and `srv_record_name` provided, then `repo_node` blocks must be dynamic (see [`dynamic`](#dynamic)).

Optional:

-   `flavor` (String) The flavor of the MongoDB deployment. Allowed values:

    -   `mongodb`
    -   `documentdb`

    The following conditions apply:

    -   The `documentdb` flavor cannot be combined with the MongoDB Server type `sharded`.

-   `replica_set_name` (String) Name of the replica set, if applicable.
-   `srv_record_name` (String) Name of a DNS SRV record which contains cluster topology details. If specified, then all `repo_node` blocks must be declared dynamic (see [`dynamic`](#dynamic)). Only supported for `server_type="sharded"` or `server_type="replicaset".

<a id="nestedblock--redshift_settings"></a>

### Nested Schema for `redshift_settings`

Optional:

-   `aws_region` (String) Code of the AWS region where the Redshift instance is deployed.
-   `cluster_identifier` (String) Name of the provisioned cluster.
-   `workgroup_name` (String) Workgroup name for serverless cluster.
