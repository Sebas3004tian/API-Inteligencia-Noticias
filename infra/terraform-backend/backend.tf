provider "azurerm" {
  features {}
}

resource "azurerm_storage_account" "tfstate" {
  name                     = "newsplatetfstate"
  resource_group_name      = "newsPlatform-rg"
  location                 = "eastus"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "tfstate" {
  name                  = "tfstate"
  storage_account_id    = azurerm_storage_account.tfstate.id
  container_access_type = "private"
}

