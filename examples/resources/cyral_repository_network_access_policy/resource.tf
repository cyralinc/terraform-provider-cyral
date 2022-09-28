# Repository the policy refers to
resource "cyral_repository" "my_sqlserver_repo" {
    name = "my-sqlserver-repo"
    type = "sqlserver"
    host = "sqlserver.cyral.com"
    port = 1433
}

# Allow access from IPs 1.2.3.4 and 4.3.2.1 for Admin database
# account, and from any IP address for accounts Engineer and
# Analyst.
resource "cyral_repository_network_access_policy" "my_sqlserver_repo_policy" {
    repository_id = cyral_repository.my_sqlserver_repo.id
    network_access_rule {
        name = "rule1"
        db_accounts = ["Admin"]
        source_ips = ["1.2.3.4", "4.3.2.1"]
    }
    network_access_rule {
        name = "rule2"
        db_accounts = ["Engineer", "Analyst"]
    }
}
