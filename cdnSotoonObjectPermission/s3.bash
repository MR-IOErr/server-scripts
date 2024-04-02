#!/bin/bash

set -x

aws_access_key_id=
aws_secret_access_key=
aws_default_region=default

bucket_name=""

# Set AWS profile and endpoint URL
aws configure set aws_access_key_id "$aws_access_key_id"
aws configure set aws_secret_access_key "$aws_secret_access_key"
aws configure set default.region "$aws_default_region"

# List all objects in the bucket
objects=$(aws s3api list-objects --bucket "$bucket_name" --query 'Contents[].Key' --endpoint-url https://s3.thr1.sotoon.ir --profile sotoon --output text)

# Iterate through each object and set its ACL to public-read
for object in $objects; do
    # Skip objects that match _nuxt/*
    if [[ $object == *"_nuxt/"* ]]; then
       continue
    else
       aws s3api put-object-acl --bucket "$bucket_name" --key "$object" --acl public-read --endpoint-url https://s3.thr1.sotoon.ir --profile sotoon
    fi

done
