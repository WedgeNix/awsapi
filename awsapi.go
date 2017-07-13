package awsapi

import (
	"os"

	"encoding/json"

	"fmt"

	"github.com/WedgeNix/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// AwsController method struct for StartAWS
type AwsController struct {
	c3svc  *s3.S3
	bucket string
}

// StartAWS starts a AWS method
func StartAWS() *AwsController {
	bucket := os.Getenv("AWS_BUCKET")
	return &AwsController{
		bucket: bucket,
		c3svc: s3.New(session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewEnvCredentials(),
		}))),
	}
}

// GetObject gets JSON from AWS S3 populates the custom struct of the file
// key is the dir + "/" + filename
func (ac *AwsController) GetObject(key string, inter interface{}) {
	result, err := ac.c3svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(ac.bucket),
		Key:    aws.String(key),
	})
	util.E(err)
	json.NewDecoder(result.Body).Decode(&inter)
}

// PutObject sends a file to AWS S3 bucket, uses name of file on system, and AWS directory if there is one
func (ac *AwsController) PutObject(name string, dir ...string) *s3.PutObjectOutput {
	f, err := os.Open(name)
	var key string
	if len(dir) > 0 {
		key = dir[0] + "/" + name
	} else {
		key = name
	}
	util.E(err)
	imput := &s3.PutObjectInput{
		Body:                 aws.ReadSeekCloser(f),
		Bucket:               aws.String(ac.bucket),
		Key:                  aws.String(key),
		ServerSideEncryption: aws.String("AES256"),
	}
	result, err := ac.c3svc.PutObject(imput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}
	return result
}

// GetList gets a list of files in the bucket
func (ac *AwsController) GetList() *s3.ListObjectsOutput {
	input := &s3.ListObjectsInput{
		Bucket: aws.String(ac.bucket),
	}

	result, err := ac.c3svc.ListObjects(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				fmt.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}
	return result
}
