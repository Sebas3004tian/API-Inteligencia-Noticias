resource "azurerm_resource_group" "this" {
  name     = var.name
  location = var.location
}

output "name" {
  value = azurerm_resource_group.this.name
}
