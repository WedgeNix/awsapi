package awsapi

import (
	"os"
	"time"

	"encoding/json"

	"fmt"

	"strings"

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

// Monitor holds SKU information needed for just-in-time calculations.
type Monitor struct {
	Sold    int
	Days    int
	Then    time.Time
	Pending bool
}

// Object represents any package-level type.
type Object interface {
	__()
}

// ObjectMonitors maps SKUs to their respective just-in-time data.
type ObjectMonitors map[string]Monitor

func (_ ObjectMonitors) __() {}

// GetObject gets JSON from AWS S3 populates the custom struct of the file
// key is the dir + "/" + filename
// returns false if err reads NoSuchKey meaning does not exist. Can read true
// if another error happens, so must determain how to handle error
func (ac *AwsController) GetObject(key string, o Object) (bool, error) {
	result, err := ac.c3svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(ac.bucket),
		Key:    aws.String(key),
	})
	if err == nil {
		json.NewDecoder(result.Body).Decode(o)
		return true, nil
	}

	return !strings.Contains(err.Error(), "NoSuchKey"), err
}

// PutObject sends a file to AWS S3 bucket, uses name of file.
// This will Put the file in the main bucket directory.
func (ac *AwsController) PutObject(name string) *s3.PutObjectOutput {
	f, err := os.Open(name)
	key := splitName(name)
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

// PutObjectDir sends a file to AWS S3 bucket, uses name of file and directory to save file.
// This will Put the file in the directory specified.
func (ac *AwsController) PutObjectDir(name string, dir string) *s3.PutObjectOutput {
	f, err := os.Open(name)
	keyName := splitName(name)
	key := dir + "/" + keyName
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

func splitName(name string) string {
	sp := strings.Split(name, "/")
	return sp[len(sp)-1]
}
