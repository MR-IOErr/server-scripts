import json
import os
import time

import boto3


def upload_to_s3(data, bucket, filename, content_type=None, acl=None):
    if isinstance(data, str):
        data = data.encode('utf8')
        content_type = content_type or 'text/plain'
    elif isinstance(data, (dict, list)):
        data = json.dumps(data).encode('utf8')
        content_type = content_type or 'application/json'
    else:
        content_type = content_type or 'application/octet-stream'
    s3_resource = boto3.resource(
        's3',
        endpoint_url=os.environ.get('S3_ENDPOINT') or 'https://s3.ir-thr-at1.arvanstorage.ir',
        aws_access_key_id=os.environ.get('S3_ID'),
        aws_secret_access_key=os.environ.get('S3_SECRET'),
    )
    bucket = s3_resource.Bucket(bucket)
    return bucket.put_object(
        ACL=acl or 'public-read',
        Body=data,
        Key=filename,
        ContentType=content_type,
    )


if __name__ == '__main__':
    status = {"status": "ok", "lastUpdate": int(time.time())}
    upload_to_s3(status, 'nobitex-cdn', 'data/health.txt')
    print(status)
