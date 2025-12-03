resource "azurerm_kubernetes_cluster" "aks" {
  name                = "${var.project_name}-aks"
  location            = var.location
  resource_group_name = var.resource_group
  dns_prefix          = "${var.project_name}-dns"


  default_node_pool {
    name           = "systempool"
    node_count     = var.node_count
    vm_size        = var.node_vm_size
    vnet_subnet_id = var.subnet_id
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin      = "azure"
    load_balancer_sku   = "standard"
    
    service_cidr        = "172.16.0.0/16" 
    dns_service_ip      = "172.16.0.10" 
  }

}

output "kube_config" {
  value     = azurerm_kubernetes_cluster.aks.kube_config_raw
  sensitive = true
}

output "kubelet_object_id" {
  value = azurerm_kubernetes_cluster.aks.kubelet_identity[0].object_id
}