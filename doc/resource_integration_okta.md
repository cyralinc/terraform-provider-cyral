# Okta Integration

CRUD operations for Okta integration.

## Usage

```hcl
resource "cyral_integration_okta" "my-okta" {
    name = ""
    certificate = ""
    email_domains = [""]
    signin_url = ""
    signout_url = ""
}
```

## Variables

|  Name           |  Default  |  Description                                                          | Required |
|:----------------|:---------:|:----------------------------------------------------------------------|:--------:|
| `name`          |           | Integration name that will be used internally in Control Plane.       | Yes      |
| `certificate`   |           | Okta Certificate.                                                     | Yes      |
| `email_domains` |           | List of allowed signin domains.                                       | No       |
| `signin_url`    |           | Okta Signin URL. Make sure to include a valid URL starting with https://  | Yes      |
| `signout_url`   |           | Okta Signout URL. Make sure to include a valid URL starting with https:// | No       |

## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `id`         | Unique ID of the resource in the Control Plane.                     |
