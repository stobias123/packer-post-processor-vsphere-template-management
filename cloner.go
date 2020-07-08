package main

import (
	"fmt"
	"context"
	"time"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/session"
	libfinder "github.com/vmware/govmomi/vapi/library/finder"
	"github.com/vmware/govmomi/vapi/library"
	"github.com/vmware/govmomi/vapi/vcenter"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
)

// AbstractCloner is an interface of Cloner
type AbstractCloner interface {
	CloneToContentLibrary() error
	//DeleteTemplate()
}	
// Cloner is a wrapper of the rest Client.
type Cloner struct {
	restClient rest.Client
	config     Config
	finder     *find.Finder
	now        time.Time
}

// NewCloner creates a Cloner instance from config.
func NewCloner(config Config) (*Cloner, error) {
	ctx := context.TODO()

	vcenterUrl, err := url.Parse(fmt.Sprintf("https://%v/sdk", config.VCenterServer))
	if err != nil {
		return nil, err
	}
	credentials := url.UserPassword(config.Username, config.Password)
	vcenterUrl.User = credentials

	soapClient := soap.NewClient(vcenterUrl, config.InsecureConnection)
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
	cloner := Cloner{
		restClient: *restClient,
		finder: finder,
	}
	return &cloner, nil
}

// CloneToContentLibrary clones a template artifact to a Content library.
func(c *Cloner) CloneToContentLibrary(config *Config) error{
	ctx := context.TODO()
	libManager := library.NewManager(&c.restClient)
	libraryFinder := libfinder.NewFinder(libManager)
	foundObj,err  := libraryFinder.Find(ctx, config.ContentLibrary)
	if err != nil {
		return err
	}

	// Error if nothing is found
	if len(foundObj) != 1 {
		return fmt.Errorf("`ContentLibrary` %s not found! Please make sure that it is set correctly", config.ContentLibrary)
	}
	// Errors when finder gets something but it's the wrong type.
	library, ok := foundObj[0].GetResult().(library.Library)
	if !ok {
		return fmt.Errorf("%q is a %T", config.ContentLibrary, foundObj[0].GetResult())
	}

	vm, err := c.finder.VirtualMachine(ctx, config.LocationConfig.VMName)
	if err != nil {
		return err
	}

	ovf := vcenter.OVF{
		Spec: vcenter.CreateSpec{
			Name: config.LocationConfig.VMName,
		},
		Source: vcenter.ResourceID{
			Value: vm.Reference().Value,
		},
		Target: vcenter.LibraryTarget{
			LibraryID: library.ID,
		},
	}
	_, err = vcenter.NewManager(&c.restClient).CreateOVF(ctx, ovf)
	if err != nil {
		return err
	}
	return nil
}