# Repository Binding

This resource provides CRUD operations in policy rules.


## Usage

```hcl
resource "cyral_policy_rule" "SOME_RESOURCE_NAME" {
    policy_id = ""
    identities {
        db_roles = [""]
        groups = [""]
        services = [""]
        users = [""]
    }
    deletes {
        additional_checks = ""
        data = [""]
        dataset_rewrites {
            dataset = ""
            repo = ""
            substitution = ""
            parameters = [""]
        }
    }
    reads {
        additional_checks = ""
        data = [""]
        dataset_rewrites {
            dataset = ""
            repo = ""
            substitution = ""
            parameters = [""]
        }
    }
    updates {
        additional_checks = ""
        data = [""]
        dataset_rewrites {
            dataset = ""
            repo = ""
            substitution = ""
            parameters = [""]
        }
    }
}
```

## Variables

|  Name           |  Default    |  Description                                                                         | Required |
|:----------------|:-----------:|:-------------------------------------------------------------------------------------|:--------:|


