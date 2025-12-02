variable "resource_group_name" {
  type        = string
  description = "Resource Group name"
  default     = "rg-api-news"
}

variable "location" {
  type    = string
  default = "eastus"
}

variable "acr_name" {
  type    = string
  default = "myaCRname12345" # replace
}

variable "acr_sku" {
  type    = string
  default = "Basic"
}

variable "cluster_name" {
  type    = string
  default = "aks-api-news"
}

variable "node_count" {
  type    = number
  default = 1
}
