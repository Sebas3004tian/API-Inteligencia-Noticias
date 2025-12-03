resource "azurerm_virtual_network" "vnet" {
  name                = "${var.project_name}-vnet"
  location            = var.location
  resource_group_name = var.resource_group
  address_space       = [var.vnet_cidr]
}

resource "azurerm_subnet" "aks" {
  name                 = "aks-subnet"
  resource_group_name  = var.resource_group
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = [var.aks_subnet_cidr]
}

output "aks_subnet_id" {
  value = azurerm_subnet.aks.id
}
