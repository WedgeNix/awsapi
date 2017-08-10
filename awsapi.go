package awsapi

import (
	"errors"
	"io/ioutil"
	"os"

	"encoding/json"

	"strings"

	"bytes"

	"github.com/WedgeNix/awsapi/dir"
	"github.com/WedgeNix/awsapi/file"
	"github.com/WedgeNix/util"
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
func New(test ...bool) (*Controller, error) {
	var c Controller

	c.bucket = "wedgenix-app-storage"
	if len(test) > 0 && test[0] {
		c.bucket = "wedgenixtestbucket"
	}

	sess, err := session.NewSession(&aws.Config{Credentials: credentials.NewEnvCredentials()})
	if err != nil {
		return &c, err
	}

	c.c3svc = s3.New(sess)

	c.verIDs = map[string]*string{}

	return &c, nil
}

// GetVerIDs grabs a view of all captured version IDs.
func (c Controller) GetVerIDs() map[string]string {
	var bin map[string]string
	for filename, vID := range c.verIDs {
		bin[filename] = *vID
	}
	return bin
}

// OpenFile opens a generic file on AWS.
func (c *Controller) OpenFile(name string) (*os.File, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(name),
	}

	resp, err := c.c3svc.GetObject(input)

	f, err := os.Create(name)
	if err != nil {
		return f, err
	}

	if err == nil {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return f, err
		}

		_, err = f.Write(b)
		if err != nil {
			return f, err
		}

		return f, nil
	}

	if strings.Contains(err.Error(), "NoSuchKey") {
		return f, nil
	}

	return f, err
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

// SaveFile saves a generic file on AWS.
func (c *Controller) SaveFile(f *os.File) error {
	input := &s3.PutObjectInput{
		Body:                 aws.ReadSeekCloser(f),
		Bucket:               aws.String(c.bucket),
		Key:                  aws.String(f.Name()),
		ServerSideEncryption: aws.String("AES256"),
	}

	result, err := c.c3svc.PutObject(input)
	if err != nil {
		return err
	}

	c.verIDs[f.Name()] = result.VersionId

	return nil
}

// Save sends a file to AWS S3 bucket, uses name of file.
// This will Put the file in the main bucket directory.
func (c *Controller) Save(name string, f file.Any) error {
	b, err := json.Marshal(f)
	if err != nil {
		return err
	}
	util.Log(string(b))

	r := bytes.NewReader(b)

	input := &s3.PutObjectInput{
		Body:                 aws.ReadSeekCloser(r),
		Bucket:               aws.String(c.bucket),
		Key:                  aws.String(name),
		ServerSideEncryption: aws.String("AES256"),
	}

	switch f.(type) {
	case file.BananasMon:
		input.Tagging = aws.String("App Name=hit-the-bananas")
	case file.D2sVendorDays:
		input.Tagging = aws.String("App Name=drive-2-sku")
	}

	result, err := c.c3svc.PutObject(input)
	if err != nil {
		return err
	}

	c.verIDs[name] = result.VersionId

	return nil
}

// SaveDir writes the given directory to AWS at the specified path.
func (c *Controller) SaveDir(path dir.Path, d dir.Any) error {
	parts := strings.Split(string(path), "*")
	if len(parts) < 1 {
		return errors.New("no extension provided")
	}
	folder := parts[0]
	if len(parts) < 2 {
		return errors.New("no extension provided")
	}
	ext := parts[1]

	switch d := d.(type) {
	case dir.BananasMon:
		util.Log("files in dir: ", len(d))
		for name, f := range d {
			util.Log("Saving '", name, "'")
			err := c.Save(folder+name+ext, f)
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
func (c *Controller) OpenDir(path dir.Path, d dir.Any) error {
	parts := strings.Split(string(path), "/*")
	if len(parts) < 1 {
		return errors.New("no extension provided")
	}
	folder := parts[0]
	if len(parts) < 2 {
		return errors.New("no extension provided")
	}
	ext := parts[1]

	input := &s3.ListObjectsInput{
		Bucket: aws.String(c.bucket),
		Prefix: aws.String(folder),
	}
	output, err := c.c3svc.ListObjects(input)
	if err != nil {
		return err
	}

	switch d := d.(type) {
	case dir.BananasMon:
		for _, obj := range output.Contents {
			fname := *obj.Key
			if fname == folder {
				continue
			}

			if !strings.Contains(fname, ext) {
				continue
			}

			var f file.BananasMon
			exists, err := c.Open(fname, &f)
			if err != nil {
				return err
			}
			if !exists {
				return errors.New("unabled to open '" + fname + "'")
			}
			withoutDir := strings.Replace(fname, folder+`/`, "", -1)
			withoutDirExt := strings.Replace(withoutDir, ext, "", -1)
			d[withoutDirExt] = f
		}

		return nil

	default:
		return errors.New("unknown type; possibly unimplemented")
	}
}
