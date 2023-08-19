package main

import (
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	AWS_Access_Key = ""
	AWS_Secret_Key = ""
	AWS_Bucket     = ""
	AWS_S3_URL     = ""
	Permission     = "public-read"
	Region         = "default"
	contentType    = "image/svg+xml"
	localPATH      = "/root/scripts/.charts/"
	destPATH       = "/charts/"
)

func uploadToBucket(chart ChartParams, svc *s3.S3, wg *sync.WaitGroup) {
	defer wg.Done()

	file, errFile := os.Open(localPATH + chart.Name)
	if errFile != nil {
		log.Fatal(errFile)
	}

	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(AWS_Bucket),
		Key:    aws.String(destPATH + chart.Name),
		// Body:          file,
		ACL:         aws.String(Permission),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Fatal(err)
	}

}
