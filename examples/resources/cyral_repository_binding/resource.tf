terraform {
  required_providers {
    cyral = {
      source = "local/terraform/cyral"
    }
  }
}

provider "cyral" {
    # Follow the instructions in the Cyral Terraform Provider page to set up the
    # credentials: https://registry.terraform.io/providers/cyralinc/cyral/latest/docs
    client_id     = "sa/default/b06e0569-25b2-4ea1-a3a1-13e3d4a17d44"
    client_secret = "_V12lHPvvElcUchgXQaBrTjlC2R1DeTyKdbpnS_MvnX-jiAa"

    control_plane = "port-multiplex-v2.apdev.cyral.com"

}

resource "cyral_repository" "repo" {
  name = "tf-account-repo"
  type          = "mongodb"
  repo_node {
        name = "single-node-mongo"
        host = "mongodb.cyral.com"
        port = 27017
  }
}

resource "cyral_sidecar" "sidecar" {
  name = "tf-account-sidecar"
  deployment_method = "docker"
}

resource "cyral_sidecar_listener" "listener" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["mongodb"]
  network_address {
    host          = "mongodb.cyral.com"
    port          = 27017
  }
}

resource "cyral_repository_binding" "binding" {
  sidecar_id = cyral_sidecar.sidecar.id
  repository_id = cyral_repository.repo.id
  enabled = true
  listener_binding {
    listener_id = cyral_sidecar_listener.listener.listener_id
    node_index = 0
  }
}
