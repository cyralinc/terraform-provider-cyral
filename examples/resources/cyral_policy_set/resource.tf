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

resource "cyral_policy_set" "repo_lockdown_example" {
  wizard_id    = "repo-lockdown"
  name        = "default block with failopen"
  description = "This default policy will block by default all queries for myrepo except the ones not parsed by Cyral"
  enabled = true
  tags = ["default block", "fail open"]
    scope {
    repo_ids = [cyral_repository.myrepo.id]
  }
  wizard_parameters    = jsonencode({
    denyByDefault = true
    failClosed = false
    })
}
