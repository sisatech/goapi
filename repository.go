package goapi

import (
	"fmt"
	"net/http"

	"github.com/sisatech/vcli/pkg/util/file"
)

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

	fmt.Println("2")
	defer f.Close()

	r, err := http.NewRequest(http.MethodPost, c.cfg.Address+uri, f)
	if err != nil {
		return err
	}
	fmt.Println("ffdsf")
	resp, err := c.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("ffdsdssf")
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
