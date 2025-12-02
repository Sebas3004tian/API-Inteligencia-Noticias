output "kube_admin_config" {
  value     = azurerm_kubernetes_cluster.aks.kube_admin_config[0]
  sensitive = true
}

output "cluster_name" {
  value = azurerm_kubernetes_cluster.aks.name
}
