variable "vsphere_url" {}
variable "template_name" {}
variable "clone_from_template" {}
variable "vsphere_username" {}
variable "vsphere_password" {}
variable "vsphere_datastore" {}
variable "vsphere_cluster" {}
variable "ssh_password" {}
variable "content_library" {}

source "vsphere-clone" "example" {
  CPUs = 4
  RAM = 8192
  convert_to_template = true
  vm_name = var.template_name
  template = var.clone_from_template
  insecure_connection = true
  vcenter_server = var.vsphere_url
  cluster = var.vsphere_cluster
  datastore = var.vsphere_datastore

  username = var.vsphere_username
  password = var.vsphere_password

  ssh_username = "root"
  ssh_password = var.ssh_password

}

build {
  sources = [
    "source.vsphere-clone.example"
  ]
  post-processor "vsphere-postprocessor" {
      content_library = var.content_library,
      vcenter_server = var.vsphere_url,
      vcenter_username = var.vsphere_username,
      vcenter_password = var.vsphere_password
  }
}