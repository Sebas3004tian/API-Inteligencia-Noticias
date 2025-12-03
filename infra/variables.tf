variable "project_name" {
  type        = string
  description = "Project prefix"
}

variable "location" {
  type    = string
  default = "eastus"
}

variable "aks_node_count" {
  type    = number
  default = 1
}

variable "aks_vm_size" {
  type    = string
  default = "Standard_B2s"
}

variable "rg_name"{
  type    = string
  default = "newsPlatform-rg"
}
