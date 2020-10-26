# vSphere Template Management

Aims to handle post processing of vmware templates. 

Goals are the following

- [x] Copy to content library
- [ ] Copy to a range of content libraries
- [ ] Clean up old templates in vcenter
- [ ] Clean up old templates in content libraries.

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

## Comparison to standard Packer `Content-Library`
Standard packer allows for use of a `content-library` parameter. This method relies on vsphere's [vm-template api](https://vdc-download.vmware.com/vmwb-repository/dcr-public/423e512d-dda1-496f-9de3-851c28ca0814/0e3f6e0d-8d05-4f0c-887b-3d75d981bae5/VMware-vSphere-Automation-SDK-REST-6.7.0/docs/apidocs/index.html#PKG_com.vmware.vcenter.vm_template) and as such requires VMWare 6.7

This plugin works and has been tested with vsphere 6.5
