---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cyral_sidecar_health Data Source - terraform-provider-cyral"
subcategory: ""
description: |-
    Retrieve aggregated information about the sidecar's health https://cyral.com/docs/sidecars/manage/#check-sidecar-cluster-status, considering all instances of the sidecar.
---

# cyral_sidecar_health (Data Source)

Retrieve aggregated information about the [sidecar's health](https://cyral.com/docs/sidecars/manage/#check-sidecar-cluster-status), considering all instances of the sidecar.

<!-- schema generated by tfplugindocs -->

## Schema

### Required

-   `sidecar_id` (String) ID of the Sidecar that will be used to retrieve health information.

### Read-Only

-   `id` (String) Data source identifier.
-   `status` (String) Sidecar health status. Possible values are: `HEALTHY`, `DEGRADED`, `UNHEALTHY` and `UNKNOWN`. For more information, see [Sidecar Status](https://cyral.com/docs/sidecars/manage/#check-sidecar-cluster-status).
