package utils

import (
	"bytes"
	"video/pkg/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadVideoToS3(videoData []byte,s3Path string) error {
	c,err:=config.LoadConfig()
	if err != nil{
		return err
	}

	sess,err:=session.NewSession(&aws.Config{
		Region: &c.Region,
		Credentials: credentials.NewStaticCredentials(
			c.AWS_ACCESS_KEY_ID,
			c.AWS_SECRET_ACCESS_KEY,
			"",
		),
	})
	if err != nil{
		return err
	}
	uploader:=s3manager.NewUploader(sess)

	_,err=uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("streamerz"),
		Key: aws.String(s3Path),
		Body: bytes.NewReader(videoData),
	})
	if err != nil{
		return err
	}
	return nil
}