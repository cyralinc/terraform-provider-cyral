resource "cyral_repository" "myrepo" {
    type = "mongodb"
    name = "myrepo"

    repo_node {
        name = "node-1"
        host = "mongodb.cyral.com"
        port = 27017
    }

    mongodb_settings {
      server_type = "standalone"
    }
}

resource "cyral_policy_wizard_v1" "repo_lockdown_example" {
  wizardId    = "repo_lockdown"
  name        = "block by default"
  description = "Block all access to this repository by default, but allow queries not parsed by Cyral"
  enabled     = true
  tags        = ["default block", "fail open"]
  scope {
    repo_ids = [cyral_repository.myrepo.id]
  }
  wizardParameters    = jsonencode({
    denyByDefault = true
    })
  enabled    = true
}
