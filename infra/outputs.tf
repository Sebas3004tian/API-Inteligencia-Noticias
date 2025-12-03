output "kubeconfig" {
  value     = module.aks.kube_config
  sensitive = true
}
