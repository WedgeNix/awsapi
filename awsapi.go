package awsapi

import (
	"errors"
	"os"

	"encoding/json"

	"strings"

	"bytes"

	"github.com/WedgeNix/awsapi/dir"
	"github.com/WedgeNix/awsapi/file"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Controller method struct for StartAWS
type Controller struct {
	c3svc  *s3.S3
	bucket string
	verIDs map[string]*string
}

// New starts a AWS method
func New() *Controller {
	bucket := os.Getenv("AWS_BUCKET")
	return &Controller{
		bucket: bucket,
		c3svc: s3.New(session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewEnvCredentials(),
		}))),
	}
}

// GetVerIDs grabs a view of all captured version IDs.
func (c Controller) GetVerIDs() map[string]string {
	var bin map[string]string
	for filename, vID := range c.verIDs {
		bin[filename] = *vID
	}
	return bin
}

// Open gets JSON from AWS S3 populates the custom struct of the file
// key is the dir + "/" + filename
// returns false if err reads NoSuchKey meaning does not exist. Can read true
// if another error happens, so must determain how to handle error
func (c *Controller) Open(name string, f file.Any) (bool, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(name),
	}
	resp, err := c.c3svc.GetObject(input)

	if err == nil {
		json.NewDecoder(resp.Body).Decode(f)
		return true, nil
	}

	if strings.Contains(err.Error(), "NoSuchKey") {
		return false, nil
	}

	return true, err
}

// Save sends a file to AWS S3 bucket, uses name of file.
// This will Put the file in the main bucket directory.
func (c *Controller) Save(name string, f file.Any) error {
	b, err := json.Marshal(f)
	if err != nil {
		return err
	}

	r := bytes.NewReader(b)

	input := &s3.PutObjectInput{
		Body:                 aws.ReadSeekCloser(r),
		Bucket:               aws.String(c.bucket),
		Key:                  aws.String(name),
		ServerSideEncryption: aws.String("AES256"),
	}

	result, err := c.c3svc.PutObject(input)
	if err != nil {
		return err
	}

	c.verIDs[name] = result.VersionId

	return nil
}

// SaveDir writes the given directory to AWS at the specified path.
func (c *Controller) SaveDir(d dir.Any) error {
	switch d := d.(type) {
	case dir.Monitor:
		for name, f := range d {
			err := c.Save(name, f)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return errors.New("unknown type; possibly unimplemented")
	}
}

// OpenDir gets a list of files in the bucket
func (c *Controller) OpenDir(name string, d dir.Any) error {
	if len(name) > 0 && string(name[len(name)-1]) != `/` {
		name += `/`
	}

	input := &s3.ListObjectsInput{
		Bucket: aws.String(c.bucket),
		Prefix: aws.String(name),
	}
	output, err := c.c3svc.ListObjects(input)
	if err != nil {
		return err
	}

	switch d := d.(type) {
	case dir.Monitor:
		for _, obj := range output.Contents {
			fname := *obj.Key
			if fname == name {
				continue
			}

			var f file.Monitor
			_, err := c.Open(fname, &f)
			if err != nil {
				return err
			}
			d[fname] = f
		}

		return nil

	default:
		return errors.New("unknown type; possibly unimplemented")
	}
}
