package main

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer/packer"
)

// PostProcessor is the core of this library
// Packer performs `PostProcess()` method of this processor
type PostProcessor struct {
	testMode bool
	cleaner  AbstractCloner
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
	if p.config.Identifier == "" {
		return errors.New("empty `identifier` is not allowed. Please make sure that it is set correctly")
	}
	return nil
}

// PostProcess clones a vsphere template to a content library
func (p *PostProcessor) PostProcess(ctx context.Context, ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, bool, error) {
	log.Println("Running the post-processor")
	return nil, false, false, nil
}

func (p *PostProcessor) uiMessage(message string) string {
	if p.config.DryRun {
		return "[DryRun] " + message
	}
	return message
}
