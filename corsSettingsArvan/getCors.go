package main

// snippet-start:[s3.go.set_cors.imports]
import (
  "fmt"
  "os"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
)

// Lists CORS configuration
//
// Usage:
//    go run s3_get_bucket_cors.go BUCKET_NAME
func main() {

  if len(os.Args) != 2 {
      exitErrorf("Bucket name required\nUsage: go run", os.Args[0], "BUCKET")
  }

  bucket := os.Args[1]

  sess, err := session.NewSession(&aws.Config{
      Credentials: credentials.NewStaticCredentials("AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", ""),
  })
  svc := s3.New(sess, &aws.Config{
      Region:   aws.String("default"),
      Endpoint: aws.String("https://s3.ir-thr-at1.arvanstorage.ir"),
  })
  input := &s3.GetBucketCorsInput{
      Bucket: aws.String(bucket),
  }

  result, err := svc.GetBucketCors(input)
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
      return
  }

  fmt.Println(result)
}

// snippet-start:[s3.go.set_cors.exit]
func exitErrorf(msg string, args ...interface{}) {
  fmt.Fprintf(os.Stderr, msg+"\n", args...)
  os.Exit(1)
}


