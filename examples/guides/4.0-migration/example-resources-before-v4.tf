resource "cyral_repository" "mongo_repo" {
  name = "mongo-repo"
  host = "mongodb.cyral.com"
  port = 27017  # This is the port in the database host
  type = "mongodb"
}

resource "cyral_repository_binding" "mongo_binding" {
  enabled                       = true
  repository_id                 = cyral_repository.mongo_repo.id
  sidecar_id                    = cyral_sidecar.sidecar.id
  listener_port                 = 27020  # This is the port the sidecar will expose
                                          # for users/applications connecting to the
                                          # database.
  sidecar_as_idp_access_gateway = true
}
