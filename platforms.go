package goapi

import (
	"errors"
	"log"

	"github.com/sisatech/goapi/pkg/objects"
)

// GenericPlatformArgs are the general arguments every platform that gets created takes.
type GenericPlatformArgs struct {
	Name    string
	Verbose bool
}

// GeneralHypervisorArgs are the variables required to create a new KVM or VirtualBox platform.
type GeneralHypervisorArgs struct {
	GenericPlatformArgs
	Bridge   string
	Headless bool
}

// VMwareHypervisorArgs are the variables required to create a workstation, fusion or player platform using vmware.
type VMwareHypervisorArgs struct {
	GeneralHypervisorArgs
	Type string
}

// GCPCloudArgs are the variables requried to create a GCP platform.
type GCPCloudArgs struct {
	GenericPlatformArgs
	Bucket      string
	Zone        string
	Network     string
	KeyFilePath string
	Verbose     bool
}

// AWSCloudArgs are the variables required to create a AWS platform.
type AWSCloudArgs struct {
	GenericPlatformArgs
	Bucket string
	Zone   string
	KeyID  string
	Secret string
}

// VcenterCloudArgs are the variables required to create a Vcenter platform
type VcenterCloudArgs struct {
	GenericPlatformArgs
	Username   string
	Password   string
	Address    string
	Datacenter string
	Network    string
	Datastore  string
	Cluster    string
}

type AzureCloudArgs struct {
	GenericPlatformArgs
	StorageAccount    string `json:"storageAccount"`
	StorageAccountKey string `json:"storageAccountKey"`
	VirtualNetwork    string `json:"vnet"`
	AuthFilePath      string `json:"authpath"`
	ResourceGroup     string `json:"resourceGroup"`
	Container         string `json:"container"`
	Location          string `json:"location"`
	SubID             string `json:"subID"`
	SubNetwork        string `json:"subnet"`
}

// ListPlatforms gathers a list of platforms currently up on vorteil.
func (c *Client) ListPlatforms() ([]objects.Platform, error) {

	req := c.NewRequest(`query{
		listPlatforms{
			name
			type
		}
	}`)

	type responseContainer struct {
		ListPlatforms []objects.Platform `json:"listPlatforms"`
	}
	listPlatformsWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &listPlatformsWrapper); err != nil {
		return nil, err
	}

	return listPlatformsWrapper.ListPlatforms, nil
}

// NewPlatformAzure is a function used to create a new platform on vorteil running on the azure cloud.
func (c *Client) NewPlatformAzure(args AzureCloudArgs) (*objects.Platform, error) {
	req := c.NewRequest(`mutation($name: String!, $verbose: Boolean!, $storageAcc: String!, $storageAccKey: String!, $resourceGroup: String!, $authpath: String!, $container: String!, $location: String!, $subid: String!, $subnet: String!, $vnet: String!){
		addPlatformAzure(name: $name, verbose: $verbose, storageAccount: $storageAcc, storageAccountKey: $storageAccKey, resourceGroup: $resourceGroup, authpath: $authpath, container: $container, location: $location, subID: $subid, subnet: $subnet, vnet: $vnet){
			name
			type
		}
	}`)

	req.Var("name", args.Name)
	req.Var("verbose", args.Verbose)
	req.Var("storageAcc", args.StorageAccount)
	req.Var("storageAccKey", args.StorageAccountKey)
	req.Var("resourceGroup", args.ResourceGroup)
	req.Var("authpath", args.AuthFilePath)
	req.Var("container", args.Container)
	req.Var("location", args.Location)
	req.Var("subid", args.SubID)
	req.Var("subnet", args.SubNetwork)
	req.Var("vnet", args.VirtualNetwork)

	type responseContainer struct {
		AddPlatformAzure objects.Platform `json:"addPlatformAzure"`
	}

	azureWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &azureWrapper); err != nil {
		return nil, err
	}

	return &azureWrapper.AddPlatformAzure, nil
}

// NewPlatformVcenter is a function used to create a new platform on Vorteil talking to VCenter.
func (c *Client) NewPlatformVcenter(args VcenterCloudArgs) (*objects.Platform, error) {

	req := c.NewRequest(`mutation($name: String!, $username: String!, $password: String!, $address: String!, $cluster: String!, $datacenter: String!, $datastore: String!, $network: String!, $verbose: Boolean!){
		addPlatformVcenter(name: $name, username: $username, password: $password, address: $address, cluster: $cluster, datacenter: $datacenter, datastore: $datastore, network: $network, verbose: $verbose){
			name
			type
		}
	}`)

	req.Var("name", args.Name)
	req.Var("username", args.Username)
	req.Var("verbose", args.Verbose)
	req.Var("password", args.Password)
	req.Var("address", args.Address)
	req.Var("cluster", args.Cluster)
	req.Var("datacenter", args.Datacenter)
	req.Var("datastore", args.Datastore)
	req.Var("network", args.Network)

	type responseContainer struct {
		AddPlatformVcenter objects.Platform `json:"addPlatformVcenter"`
	}

	vcenterWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &vcenterWrapper); err != nil {
		return nil, err
	}

	return &vcenterWrapper.AddPlatformVcenter, nil
}

// NewPlatformAWS is a function used to create a new platform on Vorteil talking to AWS.
func (c *Client) NewPlatformAWS(args AWSCloudArgs) (*objects.Platform, error) {
	req := c.NewRequest(`mutation($bucket: String!, $zone: String!, $name: String!, $keyid: String!, $secret: String!, $verbose: Boolean!){
		addPlatformAWS(name: $name, bucket: $bucket, zone: $zone, keyid: $keyid, verbose: $verbose, accesskey: $secret){
			name
			type
		}
	}`)

	req.Var("bucket", args.Bucket)
	req.Var("name", args.Name)
	req.Var("keyid", args.KeyID)
	req.Var("accesskey", args.Secret)
	req.Var("secret", args.Secret)
	req.Var("verbose", args.Verbose)
	req.Var("zone", args.Zone)

	type responseContainer struct {
		AddPlatformAWS objects.Platform `json:"addPlatformAWS"`
	}

	awsWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &awsWrapper); err != nil {
		return nil, err
	}

	return &awsWrapper.AddPlatformAWS, nil
}

// NewPlatformGCP is a function used to create a new platform on Vorteil using GCP.
func (c *Client) NewPlatformGCP(args GCPCloudArgs) (*objects.Platform, error) {
	network := args.Network
	if network == "" {
		network = "default"
	}

	req := c.NewRequest(`mutation($verbose: Boolean!, $name: String!, $bucket: String!, $zone: String!, $network: String!, $key: String!){
		addPlatformGCP(name: $name, bucket: $bucket, zone: $zone, network: $network, key: $key, verbose: $verbose){
			name
			type
		}
	}`)

	req.Var("verbose", args.Verbose)
	req.Var("name", args.Name)
	req.Var("bucket", args.Bucket)
	req.Var("zone", args.Zone)
	req.Var("network", args.Network)
	req.Var("key", args.KeyFilePath)

	type responseContainer struct {
		AddPlatformGCP objects.Platform `json:"addPlatformGCP"`
	}

	gcpWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &gcpWrapper); err != nil {
		return nil, err
	}

	return &gcpWrapper.AddPlatformGCP, nil
}

// NewPlatformVMware is a function used to create any VMware hypervisor platform on Vorteil.
func (c *Client) NewPlatformVMware(args VMwareHypervisorArgs) (*objects.Platform, error) {
	switch args.Type {
	case "":
		args.Type = "player"
	case "player":
	case "workstation":
	case "fusion":
	default:
		return nil, errors.New("unknown type, must be 'player', 'workstation', or 'fusion'")
	}

	req := c.NewRequest(`mutation($name: String!, $headless: Boolean!, $verbose: Boolean!, $bridgedNetwork: String!, $type: VMwareTypes!){
		addPlatformVMware(name: $name, headless: $headless, verbose: $verbose, bridgedNetwork: $bridgedNetwork, type: $type){
			name
			type
		}
	}`)

	req.Var("name", args.Name)
	req.Var("headless", args.Headless)
	req.Var("verbose", args.Verbose)
	req.Var("bridgedNetwork", args.Bridge)
	req.Var("type", args.Type)

	type responseContainer struct {
		AddPlatformVMware objects.Platform `json:"addPlatformVMware"`
	}

	vmwareWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &vmwareWrapper); err != nil {
		return nil, err
	}

	return &vmwareWrapper.AddPlatformVMware, nil
}

// NewPlatformVirtualBox takes NewPlatformVirtualBoxArgs to create a platform on vorteil. Returns information about the platform or an error if unsuccessfull.
func (c *Client) NewPlatformVirtualBox(args GeneralHypervisorArgs) (*objects.Platform, error) {
	req := c.NewRequest(`mutation($name: String!, $headless: Boolean!, $verbose: Boolean!, $bridgedNetwork: String!){
		addPlatformVirtualBox(name: $name, headless: $headless, verbose: $verbose, bridgedNetwork: $bridgedNetwork ){
			name
			type
		}
	}`)

	req.Var("name", args.Name)
	req.Var("headless", args.Headless)
	req.Var("verbose", args.Verbose)
	req.Var("bridgedNetwork", args.Bridge)

	type responseContainer struct {
		AddPlatformVirtualBox objects.Platform `json:"addPlatformVirtualBox"`
	}

	vboxWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &vboxWrapper); err != nil {
		return nil, err
	}

	return &vboxWrapper.AddPlatformVirtualBox, nil
}

// NewPlatformKVM takes NewPlatformKVMArgs to create a platform on vorteil. Returns information about the platform or an error if unsuccessfull.
func (c *Client) NewPlatformKVM(args GeneralHypervisorArgs) (*objects.Platform, error) {
	req := c.NewRequest(`mutation($name: String!, $headless: Boolean!, $verbose: Boolean!, $bridgedNetwork: String!){
		addPlatformKVM(name: $name, headless: $headless, verbose: $verbose, bridgedNetwork: $bridgedNetwork ){
			name
			type
		}
	}`)

	req.Var("name", args.Name)
	req.Var("headless", args.Headless)
	req.Var("verbose", args.Verbose)
	req.Var("bridgedNetwork", args.Bridge)

	type responseContainer struct {
		AddPlatformKVM objects.Platform `json:"addPlatformKVM"`
	}

	kvmWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &kvmWrapper); err != nil {
		return nil, err
	}

	return &kvmWrapper.AddPlatformKVM, nil
}

// DetailPlatform takes a platform and describes the details about it.
func (c *Client) DetailPlatform(platform string) (*objects.DetailPlatform, error) {
	req := c.NewRequest(`query($name: String!){
		detailPlatform(name: $name){
			dbug
			more
		}
	}`)

	req.Var("name", platform)

	type responseContainer struct {
		DetailPlatform objects.DetailPlatform `json:"detailPlatform"`
	}

	detailPlatWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &detailPlatWrapper); err != nil {
		return nil, err
	}

	return &detailPlatWrapper.DetailPlatform, nil
}

// RemovePlatform takes a platform and deletes it from vorteil. Returns an error if unsuccessfull.
func (c *Client) RemovePlatform(platform string) error {
	req := c.NewRequest(`mutation($name: String!){
		removePlatform(name: $name)
	}`)

	req.Var("name", platform)

	if err := c.client.Run(c.ctx, req, nil); err != nil {
		log.Fatal(err)
	}

	return nil
}
