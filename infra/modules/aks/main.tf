resource "azurerm_kubernetes_cluster" "aks" {
  name                = var.cluster_name
  location            = var.location
  resource_group_name = var.resource_group_name
  dns_prefix          = var.cluster_name

  default_node_pool {
    name       = "system"
    node_count = var.node_count
    vm_size    = "Standard_B2s"
  }

  identity {
    type = "SystemAssigned"
  }

  role_based_access_control_enabled = true

  network_profile {
    network_plugin = "azure"
  }
}

# Allow AKS to pull images from ACR
resource "azurerm_role_assignment" "acr_pull" {
  principal_id         = azurerm_kubernetes_cluster.aks.kubelet_identity[0].object_id
  role_definition_name = "AcrPull"
  scope                = var.acr_id
}

output "kube_admin_config" {
  value     = azurerm_kubernetes_cluster.aks.kube_admin_config[0]
  sensitive = true
}

output "cluster_name" {
  value = azurerm_kubernetes_cluster.aks.name
}
