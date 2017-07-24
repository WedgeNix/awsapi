package awsapi

import (
	"errors"
	"os"

	"encoding/json"

	"strings"

	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mrmiguu/print"
	"github.com/mrmiguu/un"
	"github.com/wedgenix/awsapi/dir"
	"github.com/wedgenix/awsapi/file"
)

// AwsController method struct for StartAWS
type AwsController struct {
	c3svc  *s3.S3
	bucket string
}

// New starts a AWS method
func New() *AwsController {
	bucket := os.Getenv("AWS_BUCKET")
	return &AwsController{
		bucket: bucket,
		c3svc: s3.New(session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewEnvCredentials(),
		}))),
	}
}

// Get gets JSON from AWS S3 populates the custom struct of the file
// key is the dir + "/" + filename
// returns false if err reads NoSuchKey meaning does not exist. Can read true
// if another error happens, so must determain how to handle error
func (ac *AwsController) Get(key string, f file.Any) (bool, error) {

	print.Debug("preparing input")

	input := &s3.GetObjectInput{
		Bucket: aws.String(ac.bucket),
		Key:    aws.String(key),
	}

	print.Debug("getting result from AWS object")

	result, err := ac.c3svc.GetObject(input)

	if err == nil {

		print.Debug("no error; populate our version of the file")

		json.NewDecoder(result.Body).Decode(f)

		print.Debug("successfully populated")

		return true, nil
	}
	if strings.Contains(err.Error(), "NoSuchKey") {

		print.Msg("no such key for '", key, "'")

		return false, nil
	}

	print.Debug("gotten!")

	return true, err
}

// Put sends a file to AWS S3 bucket, uses name of file.
// This will Put the file in the main bucket directory.
func (ac *AwsController) Put(filename string, f file.Any) (*string, error) {
	r := bytes.NewReader(un.Bytes(json.Marshal(f)))

	input := &s3.PutObjectInput{
		Body:                 aws.ReadSeekCloser(r),
		Bucket:               aws.String(ac.bucket),
		Key:                  aws.String(filename),
		ServerSideEncryption: aws.String("AES256"),
	}

	result, err := ac.c3svc.PutObject(input)

	return result.VersionId, err
}

// PutDir writes the given directory to AWS at the specified path.
func (ac *AwsController) PutDir(path string, d dir.Any) error {
	switch d := d.(type) {
	case dir.Monitor:
		for filename, f := range d {
			_, err := ac.Put(filename, f)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return errors.New("unknown type; possibly unimplemented")
	}
}

// GetDir gets a list of files in the bucket
func (ac *AwsController) GetDir(path string, d dir.Any) error {

	print.Debug("prepare prefix and input")

	prefix := path + `/`
	input := &s3.ListObjectsInput{
		Bucket: aws.String(ac.bucket),
		Prefix: aws.String(prefix),
	}

	print.Debug("grab list of objects")

	loo, err := ac.c3svc.ListObjects(input)
	if err != nil {
		return err
	}

	// if err != nil {
	// 	if aerr, ok := err.(awserr.Error); ok {
	// 		switch aerr.Code() {
	// 		case s3.ErrCodeNoSuchBucket:
	// 			fmt.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
	// 		default:
	// 			fmt.Println(aerr.Error())
	// 		}
	// 	} else {
	// 		// Print the error, cast err to awserr.Error to get the Code and
	// 		// Message from an error.
	// 		fmt.Println(err.Error())
	// 	}
	// 	// return nil
	// 	return nil
	// }

	// print.Debug("attempt to convert to monitor")

	// dm, ok := d.(*dir.Monitor)
	// if !ok {
	// 	return errors.New("unknown type; possibly unimplemented")
	// }

	print.Debug("run through all contents")

	print.Msg(`len(loo.Contents) = `, len(loo.Contents))
	for _, obj := range loo.Contents {
		k := *obj.Key
		if prefix == k {

			// print.Debug("the key and prefix are not equal")

			continue
		}
		print.Msg(`key='`, k, `'`)

		// print.Debug("getting the actual file, populating the monitor")

		// _, err := ac.Get(k, (*dm)[k])
		// if err != nil {
		// 	print.Msg(`could not get '`, k, `'`)
		// 	return err
		// }

		// print.Msg("len=", len(*dm))
		// for k, v := range *dm {
		// 	print.Msg(`"`, k, `": `, &v)
		// }
	}

	return nil
}
