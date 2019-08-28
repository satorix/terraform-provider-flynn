resource "flynn_discovery_token" "cluster" {}

output "token" {
  value = flynn_discovery_token.cluster.token
}