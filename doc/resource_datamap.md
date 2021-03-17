# Datamap

This resource provides CRUD operations in Cyral datamaps, allowing users to Create, Read, Update and Delete datamaps.

## Usage

```hcl
resource "cyral_datamap" "SOME_RESOURCE_NAME" {
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

## Variables

|  Name           |  Default  |  Description                                                                         | Required |
|:----------------|:---------:|:-------------------------------------------------------------------------------------|:--------:|
| `mapping`       |           | Block that supports mapping attributes in repos to a given label.                    | Yes      |
| `label`         |           | Label given to the data specified in the corresponding list (ex: `your_label_id`)    | Yes      |
| `data_location` |           | Block to inform a data location set: repository name and attributes specification.   | Yes      |
| `repo`          |           | Name of the repository containing the data as specified through the Cyral management console (ex: `your_repo_name`). | Yes      |
| `attributes`    |           | List containing the specific locations of the data within the repo, following the pattern `{SCHEMA}.{TABLE}.{ATTRIBUTE}` (ex: `[your_schema_name.your_table_name.your_attr_name]`). | Yes      |


> Note: When referencing data in a Dremio repository, please include the complete location, with each nested Dremio space separated by a period. For example, an attribute `my_attr` contained by table `my_tbl` within space `inner_space` within space `outer_space` would be referenced as `outer_space.inner_space.my_tbl.my_attr`.
