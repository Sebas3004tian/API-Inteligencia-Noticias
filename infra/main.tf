terraform {
  required_version = ">= 1.4.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.80"
    }
  }

  backend "azurerm" {
    resource_group_name  = "newsPlatform-rg"
    storage_account_name = "newsplatetfstate"
    container_name       = "tfstate"
    key                  = "terraform.tfstate"
  }
}

provider "azurerm" {
  features {}
}

# ----------------------
# RESOURCE GROUP
# ----------------------
resource "azurerm_resource_group" "rg" {
  name     = "${var.project_name}-rg"
  location = var.location
}

# ----------------------
# NETWORK MODULE
# ----------------------
module "network" {
  project_name = var.project_name
  source            = "./modules/network"
  resource_group    = azurerm_resource_group.rg.name
  location          = var.location
  vnet_cidr         = "10.0.0.0/16"
  aks_subnet_cidr   = "10.0.1.0/24"
}

# ----------------------
# ACR MODULE
# ----------------------
module "acr" {
  source         = "./modules/acr"
  resource_group = azurerm_resource_group.rg.name
  location       = var.location
  project_name   = var.project_name
}

# ----------------------
# AKS MODULE
# ----------------------
module "aks" {
  source            = "./modules/aks"
  resource_group    = azurerm_resource_group.rg.name
  location          = var.location
  project_name      = var.project_name
  node_count        = var.aks_node_count
  node_vm_size      = var.aks_vm_size
  subnet_id         = module.network.aks_subnet_id
  acr_id            = module.acr.acr_id
}
