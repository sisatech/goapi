package goapi

import (
	"net/http"

	"github.com/sisatech/goapi/pkg/objects"

	"github.com/sisatech/goapi/pkg/file"
)

// ListNodes ...
func (c *Client) ListNodes() ([]objects.Node, error) {
	req := c.NewRequest(`query{
		listNodes{
			host
			name
			type
		}
	}`)

	type responseContainer struct {
		ListNodes []objects.Node `json:"listNodes"`
	}

	nodesWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &nodesWrapper); err != nil {
		return nil, err
	}
	return nodesWrapper.ListNodes, nil
}

// RemoveRepository ...
func (c *Client) RemoveRepository(name string) error {

	req := c.NewRequest(`mutation($name: String!){
		removeNode(name: $name)
	}`)

	req.Var("name", name)

	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// Push ...
func (c *Client) Push(app, tag, germ, node, bucket string, compressionLevel int, injects []string) (objects.GerminateOperation, error) {

	req := c.NewRequest(`mutation($app: String!, $tag: String, $germ: GermString!, $node: String!, $bucket: String!, $compressionLevel: Int, $injects: [String]){
		push(app: $app, tag: $tag, germ: $germ, node: $node, bucket: $bucket, compressionLevel: $compressionLevel, injections: $injects){
			job{
				description
				id
				logFilePath
				logPlainFilePath
				name
				progress {
					error
					finished
					progress
					started
					status
					total
					units
				}
			}
			uri
		}
	}`)

	req.Var("app", app)
	req.Var("tag", tag)
	req.Var("germ", germ)
	req.Var("node", node)
	req.Var("bucket", bucket)
	req.Var("compressionLevel", compressionLevel)
	req.Var("injects", injects)

	type responseContainer struct {
		Push objects.GerminateOperation `json:"push"`
	}

	pushWrapper := new(responseContainer)
	if err := c.client.Run(c.ctx, req, &pushWrapper); err != nil {
		return objects.GerminateOperation{}, err
	}

	return pushWrapper.Push, nil
}

// NewRepository ...
func (c *Client) NewRepository(insecure bool, name, host, credentials string) error {

	req := c.NewRequest(`mutation($insecure: Boolean, $name: String!, $host: String!, $credentials: String!){
		newNode(name: $name, host: $host, credentials: $credentials, insecureSkipVerify: $insecure){
			name
		}
	}`)

	req.Var("insecure", insecure)
	req.Var("name", name)
	req.Var("host", host)
	req.Var("credentials", credentials)

	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}
	return nil
}

// UploadApp ...
func (c *Client) UploadApp(uri string, f file.File) error {

	defer f.Close()

	r, err := http.NewRequest(http.MethodPost, c.cfg.Address+uri, f)
	if err != nil {
		return err
	}

	resp, err := c.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	data, _ := ioutil.ReadAll(resp.Body)
	// 	return fmt.Errorf("status [%v] %s", resp.StatusCode, data)
	// }

	return nil
}

// NewBucket ...
func (c *Client) NewBucket(bucket string) error {
	req := c.NewRequest(`mutation($bucket: String!){
		newBucket(name: $bucket){
			name
		}
	}`)

	req.Var("bucket", bucket)

	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// NewAppVersion ...
func (c *Client) NewAppVersion(app, tag, bucket string) (string, error) {
	req := c.NewRequest(`mutation($bucket: String!, $app: String!, $tag: String!){
		newAppVersion(app: $app, tag: $tag, bucket: $bucket)
	}`)

	req.Var("bucket", bucket)
	req.Var("app", app)
	req.Var("tag", tag)

	type responseContainer struct {
		NewApp string `json:"newApp"`
	}

	newAppWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &newAppWrapper); err != nil {
		return "", err
	}

	return newAppWrapper.NewApp, nil
}

// NewApp return url to post the application to and an error if not successfull.
func (c *Client) NewApp(app, tag, bucket string) (string, error) {
	req := c.NewRequest(`mutation($bucket: String!, $app: String!, $tag: String!){
		newApp(app: $app, tag: $tag, bucket: $bucket)
	}`)

	req.Var("bucket", bucket)
	req.Var("app", app)
	req.Var("tag", tag)

	type responseContainer struct {
		NewApp string `json:"newApp"`
	}

	newAppWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &newAppWrapper); err != nil {
		return "", err
	}

	return newAppWrapper.NewApp, nil
}

// DeleteBucket ... deletes a bucket from a repository
func (c *Client) DeleteBucket(bucket string) error {
	req := c.NewRequest(`mutation($bucket: String!){
		deleteBucket(name: $bucket)
	}`)

	req.Var("bucket", bucket)

	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}

	return nil

}

// DeleteApp ... deletes an application from a repository.
func (c *Client) DeleteApp(bucket, app string) error {
	req := c.NewRequest(`mutation($bucket: String!, $app: String!){
		deleteApp(bucketName: $bucket, appName: $app)
	}`)

	req.Var("bucket", bucket)
	req.Var("app", app)

	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// DeleteAppVersion ... deletes an application from repository parsing a version tag as well.
func (c *Client) DeleteAppVersion(bucket, app, vers string) error {
	req := c.NewRequest(`mutation($bucket: String!, $app: String!, $vers: String!){
		deleteAppVersion(bucketName: $bucket, reference: $vers, appName: $app){
			name
		}
	}`)

	req.Var("bucket", bucket)
	req.Var("app", app)
	req.Var("vers", vers)

	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}

	return nil
}
