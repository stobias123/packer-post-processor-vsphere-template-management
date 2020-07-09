//go:generate mapstructure-to-hcl2 -type Config

package main

import (
	vspherecommon "github.com/hashicorp/packer/builder/vsphere/common"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/template/interpolate"
)

// Config is a post-processor's configuration
// PostProcessor generates it using Packer's configuration in `Configure()` method
type Config struct {
	common.PackerConfig          `mapstructure:",squash"`
	vspherecommon.LocationConfig `mapstructure:",squash"`
	vspherecommon.ConnectConfig  `mapstructure:",squash"`

	Identifier    string `mapstructure:"identifier"`
	VCenterServer string `mapstructure:"vcenter_server"`
	// TODO: Get the ConnectConfig to map out properly.
	VCenterUsername string `mapstructure:"vcenter_username"`
	VCenterPassword string `mapstructure:"vcenter_password"`
	ContentLibrary  string `mapstructure:"content_library"`
	KeepReleases    int    `mapstructure:"keep_releases"`
	KeepDays        int    `mapstructure:"keep_days"`
	DryRun          bool   `mapstructure:"dry_run"`

	ctx interpolate.Context
}
