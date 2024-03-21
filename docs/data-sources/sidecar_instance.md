# cyral_sidecar_instance (Data Source)

Retrieve sidecar instances.

## Schema

### Required

-   `sidecar_id` (String) Sidecar identifier.

### Read-Only

-   `id` (String) Data source identifier.
-   `instance_list` (List of Object) List of existing sidecar instances. (see [below for nested schema](#nestedatt--instance_list))

<a id="nestedatt--instance_list"></a>

### Nested Schema for `instance_list`

Read-Only:

-   `id` (String) Instance identifier. Varies according to the computing platform that the sidecar is deployed to.
-   `metadata` (Set of Object) Instance metadata. (see [below for nested schema](#nestedatt--instance_list--metadata))
-   `monitoring` (Set of Object) Instance monitoring information, such as its overall health. (see [below for nested schema](#nestedatt--instance_list--monitoring))

<a id="nestedatt--instance_list--metadata"></a>

### Nested Schema for `instance_list.metadata`

Read-Only:

-   `capabilities` (Set of Object) Set of capabilities that can be enabled or disabled. **Note**: This field is per-instance, not per-sidecar, because not all sidecar instances might be in sync at some point in time. (see [below for nested schema](#nestedatt--instance_list--metadata--capabilities))
-   `dynamic_version` (Boolean) If true, indicates that the instance has dynamic versioning, that means that the version is not fixed at template level and it can be automatically upgraded.
-   `last_registration` (String) The last time the instance reported to the Control Plane.
-   `recycling` (Boolean) Indicates whether the Control Plane has asked the instance to mark itself unhealthy so that it is recycled by the infrastructure.
-   `start_timestamp` (String) The time when the instance started.
-   `version` (String) Sidecar version that the instance is using.

<a id="nestedatt--instance_list--metadata--capabilities"></a>

### Nested Schema for `instance_list.metadata.capabilities`

Read-Only:

-   `recyclable` (Boolean) Indicates if sidecar instance will be recycled (e.g., by an ASG) if it reports itself as unhealthy.

<a id="nestedatt--instance_list--monitoring"></a>

### Nested Schema for `instance_list.monitoring`

Read-Only:

-   `services` (Map of Set of Object) Sidecar instance services monitoring information. (see [below for nested schema](#nestedatt--instance_list--monitoring--services))
-   `status` (String) Aggregated status of all the sidecar services.

<a id="nestedatt--instance_list--monitoring--services"></a>

### Nested Schema for `instance_list.monitoring.services`

Read-Only:

-   `status` (String) Aggregated status of sidecar service.
-   `metrics_port` (Number) Metrics port for service monitoring.
-   `components` (Map of Set of Object) Map of name to monitoring component. A component is a monitored check on the service that has its own status. (see [below for nested schema](#nestedatt--instance_list--monitoring--services--components))
-   `host` (String) Service host on the deployment.

<a id="nestedatt--instance_list--monitoring--services--components"></a>

### Nested Schema for `instance_list.monitoring.services.components`

Read-Only:

-   `status` (String) Component status.
-   `description` (String) Describes what the type of check the component represents.
-   `error` (String) Error that describes what caused the current status.
