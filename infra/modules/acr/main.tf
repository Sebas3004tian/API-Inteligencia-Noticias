resource "azurerm_container_registry" "acr" {
  name                = "${var.project_name}acr"
  resource_group_name = var.resource_group
  location            = var.location
  sku                 = "Basic"
  admin_enabled       = true
}

output "acr_id" {
  value = azurerm_container_registry.acr.id
}
