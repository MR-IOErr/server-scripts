package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func uploadFileToArvanS3() {
	var bucket, key string
	var timeout time.Duration

	flag.StringVar(&bucket, "b", "", AWS_S3_URL)
	flag.StringVar(&key, "k", "", RECORD)
	flag.DurationVar(&timeout, "d", 0, "30")
	flag.Parse()

	file, errFile := os.Open(RECORD)
	if errFile != nil {
		log.Fatal(errFile)
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(AWS_Access_Key, AWS_Secret_key, ""),
	}))

	svc := s3.New(sess, &aws.Config{
		Region:   aws.String("default"),
		Endpoint: aws.String("s3.ir-thr-at1.arvanstorage.ir"),
	})

	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}

	if cancelFn != nil {
		defer cancelFn()
	}
	_, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(AWS_Bucket),
		Key:    aws.String(DST),
		Body:   file,
		ACL:    aws.String(PER),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {

			fmt.Fprintf(os.Stderr, "upload canceled due to timeout, %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "failed to upload object, %v\n", err)
		}
		os.Exit(1)
	}
}
