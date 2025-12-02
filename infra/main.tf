module "rg" {
  source   = "./modules/resource_group"
  name     = var.resource_group_name
  location = var.location
}

module "network" {
  source              = "./modules/network"
  location            = var.location
  resource_group_name = module.rg.name
}

module "acr" {
  source              = "./modules/acr"
  resource_group_name = module.rg.name
  location            = var.location
  acr_name            = var.acr_name
  sku                 = "Basic"
}

module "aks" {
  source              = "./modules/aks"
  resource_group_name = module.rg.name
  location            = var.location
  cluster_name        = var.cluster_name
  node_count          = var.node_count
  acr_id              = module.acr.id
}
