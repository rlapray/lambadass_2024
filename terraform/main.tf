locals {
  env = { for tuple in regexall("(.*?)=(.*)", file("../.env")) : tuple[0] => tuple[1] }
}

module "gateway" {
  source  = "./modules/api_gateway_http"
  project = var.project
}

output "gateway" {
  value = module.gateway
}
