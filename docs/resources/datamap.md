# Datamap Resource

Provides a resource to handle [Data map](https://cyral.com/docs/policy#data-map).

## Example Usage

```hcl
resource "cyral_datamap" "some_resource_name" {
    mapping {
        label = ""
        data_location {
            repo = ""
            attributes = [""]
        }
        data_location {
            repo = ""
            attributes = [""]
        }
    }
    mapping {
        label = ""
        data_location {
            repo = ""
            attributes = [""]
        }
        data_location {
            repo = ""
            attributes = [""]
        }
    }
}
```

## Argument Reference

* `mapping` = (Required) Block that supports mapping attributes in repos to a given label.
* `label` = (Required) Label given to the data specified in the corresponding list (ex: `your_label_id`)
* `data_location` = (Required) Block to inform a data location set: repository name and attributes specification.
* `repo` = (Required) Name of the repository containing the data as specified through the Cyral management console (ex: `your_repo_name`).
* `attributes` = (Required) List containing the specific locations of the data within the repo, following the pattern `{SCHEMA}.{TABLE}.{ATTRIBUTE}` (ex: `[your_schema_name.your_table_name.your_attr_name]`).
  > Note: When referencing data in Dremio repository, please include the complete location in `attributes`, separating spaces by dots. For example, an attribute `my_attr` from table `my_tbl` within space `inner_space` within space `outer_space` would be referenced as `outer_space.inner_space.my_tbl.my_attr`. For more information, please see the [Policy Guide](https://cyral.com/docs/reference/policy/).

## Attribute Reference

* `last_updated` - Timestamp of the latest update performed on the datamap.