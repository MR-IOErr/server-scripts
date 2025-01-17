package main

// snippet-start:[s3.go.set_cors.imports]
import (
  "flag"
  "fmt"
  "os"
  "strings"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
)

// snippet-end:[s3.go.set_cors.imports]

// Configures CORS rules for a bucket by setting the allowed
// HTTP methods.
//
// Requires the bucket name, and can also take a space separated
// list of HTTP methods.
//
// Usage:
//    go run s3_put_bucket_cors.go -b BUCKET_NAME get put
func main() {
  // snippet-start:[s3.go.set_cors.vars]
  bucketPtr := flag.String("b", "", "Bucket to set CORS on, (required)")

  flag.Parse()

  if *bucketPtr == "" {
      exitErrorf("-b <bucket> Bucket name required")
  }

  methods := filterMethods(flag.Args())
  // snippet-end:[s3.go.set_cors.vars]

  // Initialize a session
  // snippet-start:[s3.go.set_cors.session]
  sess, err := session.NewSession(&aws.Config{
      Credentials: credentials.NewStaticCredentials("AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", ""),
  })
  svc := s3.New(sess, &aws.Config{
      Region:   aws.String("default"),
      Endpoint: aws.String("https://s3.ir-thr-at1.arvanstorage.ir"),
  })
  // snippet-end:[s3.go.set_cors.session]

  // Create a CORS rule for the bucket
  // snippet-start:[s3.go.set_cors.rule]
  rule := s3.CORSRule{
      AllowedHeaders: aws.StringSlice([]string{"Authorization"}),
      AllowedOrigins: aws.StringSlice([]string{"*"}),
      MaxAgeSeconds:  aws.Int64(3000),

      // Add HTTP methods CORS request that were specified in the CLI.
      AllowedMethods: aws.StringSlice(methods),
  }
  // snippet-end:[s3.go.set_cors.rule]

  // Create the parameters for the PutBucketCors API call, add add
  // the rule created to it.
  // snippet-start:[s3.go.set_cors.put]
  params := s3.PutBucketCorsInput{
      Bucket: bucketPtr,
      CORSConfiguration: &s3.CORSConfiguration{
          CORSRules: []*s3.CORSRule{&rule},
      },
  }

  _, err = svc.PutBucketCors(&params)
  if err != nil {
      // Print the error message
      exitErrorf("Unable to set Bucket %q's CORS, %v", *bucketPtr, err)
  }

  // Print the updated CORS config for the bucket
  fmt.Printf("Updated bucket %q CORS for %v\n", *bucketPtr, methods)
  // snippet-end:[s3.go.set_cors.put]
}

// snippet-start:[s3.go.set_cors.exit]
func exitErrorf(msg string, args ...interface{}) {
  fmt.Fprintf(os.Stderr, msg+"\n", args...)
  os.Exit(1)
}

// snippet-end:[s3.go.set_cors.exit]

// snippet-start:[s3.go.set_cors.filter]
func filterMethods(methods []string) []string {
  filtered := make([]string, 0, len(methods))
  for _, m := range methods {
      v := strings.ToUpper(m)
      switch v {
      case "POST", "GET", "PUT", "PATCH", "DELETE":
          filtered = append(filtered, v)
      }
  }

  return filtered
}


