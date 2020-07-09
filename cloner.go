package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/packer/packer"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/vapi/library"
	libfinder "github.com/vmware/govmomi/vapi/library/finder"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vapi/vcenter"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
)

// AbstractCloner is an interface of Cloner
type AbstractCloner interface {
	CloneToContentLibrary(ctx context.Context, artifact packer.Artifact, contentLibrary string) error
	//DeleteTemplate()
}

// Cloner is a wrapper of the rest Client.
type Cloner struct {
	restClient *rest.Client
	config     Config
	finder     *find.Finder
	now        time.Time
}

// NewCloner creates a Cloner instance from config.
func NewCloner(config Config) (*Cloner, error) {
	ctx := context.TODO()

	vcenterUrl, err := url.Parse(fmt.Sprintf("https://%v/sdk", config.ConnectConfig.VCenterServer))
	if err != nil {
		return nil, err
	}
	log.Println(fmt.Sprintf("User: %s, Password %s", config.VCenterUsername, config.VCenterPassword))
	credentials := url.UserPassword(config.VCenterUsername, config.VCenterPassword)
	vcenterUrl.User = credentials

	//Forcing insecure for now, I'm not unmarshalling config.ConnectConfig properly.
	soapClient := soap.NewClient(vcenterUrl, true)
	vimClient, err := vim25.NewClient(ctx, soapClient)
	if err != nil {
		return nil, err
	}

	vimClient.RoundTripper = session.KeepAlive(vimClient.RoundTripper, 10*time.Minute)
	client := &govmomi.Client{
		Client:         vimClient,
		SessionManager: session.NewManager(vimClient),
	}

	err = client.SessionManager.Login(ctx, credentials)
	if err != nil {
		return nil, err
	}

	finder := find.NewFinder(client.Client, false)
	datacenter, err := finder.DatacenterOrDefault(ctx, config.Datacenter)
	if err != nil {
		return nil, err
	}
	finder.SetDatacenter(datacenter)

	restClient := rest.NewClient(vimClient)
	err = restClient.Login(ctx, credentials)

	cloner := Cloner{
		restClient: restClient,
		finder:     finder,
	}
	return &cloner, nil
}

// CloneToContentLibrary clones a template artifact to a Content library.
func (c *Cloner) CloneToContentLibrary(ctx context.Context, artifact packer.Artifact, contentLibrary string) error {
	log.Println("Running Clone Operation......")
	libManager := library.NewManager(c.restClient)
	libraryFinder := libfinder.NewFinder(libManager)
	foundObj, err := libraryFinder.Find(ctx, contentLibrary)
	if err != nil {
		return fmt.Errorf("Problem finding Library\n%s", err)
	}

	// Error if nothing is found
	if len(foundObj) != 1 {
		return fmt.Errorf("`ContentLibrary` %s not found! Please make sure that it is set correctly", contentLibrary)
	}
	// Errors when finder gets something but it's the wrong type.
	library, ok := foundObj[0].GetResult().(library.Library)
	if !ok {
		return fmt.Errorf("%q is a %T", contentLibrary, foundObj[0].GetResult())
	}

	vm, err := c.finder.VirtualMachine(ctx, artifact.String())
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("Found Machine...... %s", vm.Name))

	ovf := vcenter.OVF{
		Spec: vcenter.CreateSpec{
			Name: artifact.String(),
		},
		Source: vcenter.ResourceID{
			Value: vm.Reference().Value,
		},
		Target: vcenter.LibraryTarget{
			LibraryID: library.ID,
		},
	}
	id, err := vcenter.NewManager(c.restClient).CreateOVF(ctx, ovf)
	log.Println(fmt.Sprintf("Copy with ID...... %s", id))
	if err != nil {
		return err
	}
	return nil
}
