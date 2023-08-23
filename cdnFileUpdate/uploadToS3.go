package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func uploadToS3(s3Body *os.File, s3Key string) {
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(AWS_Access_Key, AWS_Secret_Key, ""),
	}))

	svc := s3.New(sess, &aws.Config{
		Region:   aws.String(Region),
		Endpoint: aws.String(AWS_S3_URL),
	})

	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(AWS_Bucket),
		Key:    aws.String(destPATH + s3Key),
		ACL:    aws.String(Permission),
		Body:   s3Body,
	})
	if err != nil {
		log.Fatalln("Failed to upload file: ", s3Key, err)
	}

}
