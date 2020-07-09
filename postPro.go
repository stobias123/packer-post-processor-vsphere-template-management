package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
)

// PostProcessor is the core of this library
// Packer performs `PostProcess()` method of this processor
type PostProcessor struct {
	testMode bool
	cloner   AbstractCloner
	config   Config
}

// ConfigSpec returns HCL object spec
func (p *PostProcessor) ConfigSpec() hcldec.ObjectSpec {
	return p.config.FlatMapstructure().HCL2Spec()
}

// Configure generates post-processor's configuration
// This should set up our restClient.
func (p *PostProcessor) Configure(raws ...interface{}) error {
	//p.config.ConnectConfig

	err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter:  &interpolate.RenderFilter{},
	}, raws...)

	if err != nil {
		return err
	}
	if p.config.ContentLibrary == "" {
		return errors.New("empty `content_library` is not allowed. Please make sure that it is set correctly")
	}
	if p.config.VCenterServer == "" {
		return errors.New("empty `vcenter_server` is not allowed. Please make sure that it is set correctly")
	}
	if p.config.VCenterUsername == "" {
		return errors.New("empty `vcenter_username` is not allowed. Please make sure that it is set correctly")
	}
	if p.config.VCenterPassword == "" {
		return errors.New("empty `vcenter_password` is not allowed. Please make sure that it is set correctly")
	}
	return nil
}

// PostProcess clones a vsphere template to a content library
func (p *PostProcessor) PostProcess(ctx context.Context, ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, bool, error) {
	log.Println("Running the post-processor")
	//log.Println(artifact.String())
	ui.Message(p.uiMessage(fmt.Sprintf("Uploading to content library %s", p.config.ContentLibrary)))
	var err error
	p.cloner, err = NewCloner(p.config)
	if err != nil {
		return artifact, true, false, err
	}
	err = p.cloner.CloneToContentLibrary(ctx, artifact, p.config.ContentLibrary)
	if err != nil {
		return artifact, true, false, err
	}

	return artifact, true, false, nil
}

func (p *PostProcessor) uiMessage(message string) string {
	if p.config.DryRun {
		return "[DryRun] " + message
	}
	return message
}
