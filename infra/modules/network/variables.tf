variable "project_name" {
  type        = string
  description = "Project name prefix"
}

variable "resource_group" {
  type        = string
  description = "Resource group name"
}

variable "location" {
  type        = string
  description = "Azure region"
}

variable "vnet_cidr" {
  type        = string
  description = "CIDR for VNet"
}

variable "aks_subnet_cidr" {
  type        = string
  description = "CIDR for AKS subnet"
}
