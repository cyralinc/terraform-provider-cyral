# Datamap

This resource provides CRUD operations in Cyral datamaps, allowing users to Create, Read, Update and Delete datamaps.

## Usage

```hcl
resource "cyral_datamap" "SOME_SOURCE_NAME" {
    labels {
        label_id = ""
        label_info {
            repo = ""
            attributes = [""]
        }
        label_info {
            repo = ""
            attributes = [""]
        }
    }
    labels {
        label_id = ""
        label_info {
            repo = ""
            attributes = [""]
        }
        label_info {
            repo = ""
            attributes = [""]
        }
    }
}
```

## Variables

|  Name         |  Default  |  Description                                                                         | Required |
|:--------------|:---------:|:-------------------------------------------------------------------------------------|:--------:|
| `label_id`    |           | Label given to the data specified in the corresponding list (ex: `your_label_id`)    | Yes      |
| `repo`        |           | Name of the repository containing the data as specified through the Cyral management console (ex: `your_repo_name`) | Yes      |
| `attributes`  |           | List containing the specific locations of the data within the repo, following the pattern `{SCHEMA}.{TABLE}.{ATTRIBUTE}`. (ex: `[your_schema_name.your_table_name.your_attr_name]`) | Yes      |


> Note: When referencing data in a Dremio repository, please include the complete location, with each nested Dremio space separated by a period. For example, an attribute `my_attr` contained by table `my_tbl` within space `inner_space` within space `outer_space` would be referenced as `outer_space.inner_space.my_tbl.my_attr`.


