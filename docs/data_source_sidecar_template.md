# Sidecar Template

Returns the deployment template for a given sidecar.

## Usage

```hcl
data "cyral_sidecar_template" "SOME_DATA_SOURCE_NAME" {
    sidecar_id = ""
}
```

## Variables

|  Name         |  Default  |  Description                                               | Required |
|:--------------|:---------:|:-----------------------------------------------------------|:--------:|
| `sidecar_id`  |           | The sidecar id you want the template for.                  | Yes      |


## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `template`   |  The output variable that will contain the template.                | 