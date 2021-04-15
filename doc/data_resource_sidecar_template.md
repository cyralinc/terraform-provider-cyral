# Repository

Get a sidecar template for a given sidecar.

## Usage

```hcl
data "cyral_sidecar_template" "tf_test_cyral_sidecar_template" {
  sidecar_id = ""
}
```

## Variables

|  Name         |  Default  |  Description                                    | Required |:--------------|:---------:|:----------------------------------------------------------------------------------------|:--------:|
| `sidecar_id`  |           | The sidecar id you want the template for.                  | Yes      |


## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `template`   |  The output variable that will contain the template.                | 