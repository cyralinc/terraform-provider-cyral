---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cyral_integration_idp Data Source - terraform-provider-cyral"
subcategory: ""
description: |-
    ~> DEPRECATED Use resource and data source cyral_integration_idp_saml instead.
---

# cyral_integration_idp (Data Source)

~> **DEPRECATED** Use resource and data source `cyral_integration_idp_saml` instead.

<!-- schema generated by tfplugindocs -->

## Schema

### Optional

-   `display_name` (String) Filter results by the name of an existing IdP integration.
-   `type` (String) Filter results by the IdP integration type.

### Read-Only

-   `id` (String) The ID of this resource.
-   `idp_list` (List of Object) List of existing IdP integrations for the given filter criteria. (see [below for nested schema](#nestedatt--idp_list))

<a id="nestedatt--idp_list"></a>

### Nested Schema for `idp_list`

Read-Only:

-   `alias` (String)
-   `display_name` (String)
-   `enabled` (Boolean)
-   `single_sign_on_service_url` (String)
