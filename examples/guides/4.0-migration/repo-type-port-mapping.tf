locals {
  type_port_map = merge([
    for key, repo in local.repos : {
      for port in repo.ports :
      "${repo.type}_${port}" => {
          type = repo.type
          port = port
      }
    }
  ]...)
}
