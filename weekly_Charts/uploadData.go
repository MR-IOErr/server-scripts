package main

import (
	"flag"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func upload(weeklyChart []ChartParams) {
	var wg sync.WaitGroup
	var timeout time.Duration

	flag.DurationVar(&timeout, "d", 0, "30")

	flag.Parse()

	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(AWS_Access_Key, AWS_Secret_Key, ""),
	}))

	svc := s3.New(sess, &aws.Config{
		Region:   aws.String(Region),
		Endpoint: aws.String(AWS_S3_URL),
	})

	for _, chart := range weeklyChart {
		wg.Add(1)
		go uploadToBucket(chart, svc, &wg)
	}
	wg.Wait()

}
