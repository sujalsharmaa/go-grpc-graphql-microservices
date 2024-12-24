locals {
  env         = "dev"
  region      = "us-east-1"
  zone1       = "us-east-1a"
  zone2       = "us-east-1b"
  eks_name    = "url_shornter"
  eks_version = "1.31"
  db_password = "postgres"
}

variable "ingress_ports_monitoring_server" {
  type    = list(number)
  default = [22, 9090, 6379, 3000, 587]
}