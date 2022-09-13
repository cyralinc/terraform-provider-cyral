resource "cyral_repository_access_rules" "some_resource_name" {
  repository_id = ""
  user_account_id = ""

  # First rule to be considered
  rule {
    identity {
      type = "username|email|group"
      name = ""
    }
    config {
      policy_ids = ["policy1", "policy2"]
    }
    valid_from = "2022-09-19T23:50:48.606Z"
    valid_until = "2022-10-10T00:00:00Z"
  }

  # Second rule to be considered
  rule {
    identity {
      type = "username|email|group"
      name = ""
    }
  }

  # ... and so on
  rule {
    identity {
      type = "username|email|group"
      name = ""
    }
    config {
      policy_ids = ["some_mfa_policy"]
    }
  }
}
