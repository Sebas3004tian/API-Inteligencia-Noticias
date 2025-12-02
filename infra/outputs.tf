output "resource_group_name" {
  value = module.rg.name
}

output "acr_login_server" {
  value = module.acr.login_server
}

output "aks_kube_admin_config" {
  value     = module.aks.kube_admin_config
  sensitive = true
}

output "aks_cluster_name" {
  value = module.aks.cluster_name
}
