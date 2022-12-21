# Create a repository
resource "cyral_repository" "repo" {
  name = "tf-account-repo"
  type          = "mongodb"
  repo_node {
        host = "mongodb.cyral.com"
        port = 27017
  }
}

# Create a sidecar
resource "cyral_sidecar" "sidecar" {
  name = "tf-account-sidecar"
  deployment_method = "docker"
}

# Create a listener for the sidecar
resource "cyral_sidecar_listener" "listener" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["mongodb"]
  network_address {
    host          = "mongodb.cyral.com"
    port          = 27017
  }
}

# Bind the sidecar listener to the repository
resource "cyral_repository_binding" "binding" {
  sidecar_id = cyral_sidecar.sidecar.id
  repository_id = cyral_repository.repo.id
  enabled = true
  listener_binding {
    listener_id = cyral_sidecar_listener.listener.listener_id
    node_index = 0
  }
}

# Set the sidecar and binding as the access gateway
# for the repository.
resource "cyral_repository_access_gateway" "access_gateway" {
		repository_id  = cyral_repository.repo.id
		sidecar_id  = cyral_sidecar.sidecar.id
		binding_id = cyral_repository_binding.binding.binding_id
}
