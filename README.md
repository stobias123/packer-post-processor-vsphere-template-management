# vSphere Template Management

Aims to handle post processing of vmware templates.

Goals are the following

[x] Copy to content library
[] Copy to a range of content libraries
[] Clean up old templates in vcenter
[] Clean up old templates in content libraries.

## Examples

Check out the `examples/` dir for a complete example. With HCL2 we can fit a basic overview in a snippet.

```hcl
build {
  sources = [
    "source.vsphere-clone.example"
  ]
  # vsphere creds/username/url are included here b/c you could copy to content libraries in vcenters separate from the builder vcenter.
  post-processor "vsphere-template-management" {
      content_library = var.content_library
      vcenter_server = var.vsphere_url
      vcenter_username = var.vsphere_username
      vcenter_password = var.vsphere_password
  }
}
```
