resource "cyral_policy_rule" "some_resource_name" {
    policy_id = ""
    hosts = [""]
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
        rows = 1
        severity = "low"
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
        rows = 1
        severity = "low"
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
        rows = 1
        severity = "low"
    }
}
