#!/bin/bash

# Function to display usage and exit
usage() {
  echo "Usage: $0 -h <host> -P <port> -u <username> -p <password> -d <database> -b <s3_bucket> -k <encryption_password>"
  exit 1
}

# Parse input arguments
while getopts ":h:P:u:p:d:b:k:" opt; do
  case $opt in
    h) DB_HOST="$OPTARG"
    ;;
    P) DB_PORT="$OPTARG"
    ;;
    u) DB_USER="$OPTARG"
    ;;
    p) DB_PASS="$OPTARG"
    ;;
    d) DB_NAME="$OPTARG"
    ;;
    b) S3_BUCKET="$OPTARG"
    ;;
    k) ENCRYPTION_PASSWORD="$OPTARG"
    ;;
    *) usage
    ;;
  esac
done

# Check if all required variables are set
if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_USER" ] || [ -z "$DB_PASS" ] || [ -z "$DB_NAME" ] || [ -z "$S3_BUCKET" ] || [ -z "$ENCRYPTION_PASSWORD" ]; then
    echo "Missing required arguments."
    usage
fi


DATE=$(date +"%Y%m%d")
BACKUP_FILE="/tmp/${DB_NAME}_backup_$DATE.sql"

# Perform the MySQL backup using mysqldump
mysqldump -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" > "$BACKUP_FILE"

# Check if the backup was successful
if [ $? -ne 0 ]; then
  echo "MySQL dump failed. Exiting."
  exit 1
fi

# Ensure the backup file is fully written to disk
sync "$BACKUP_FILE"

# Compress the backup file using gzip
gzip -f "$BACKUP_FILE"

BACKUP_FILE_GZ="$BACKUP_FILE.gz"

# Encrypt the backup file using OpenSSL with the provided encryption password
openssl enc -aes-256-cbc -salt -in "$BACKUP_FILE_GZ" -out "$BACKUP_FILE_GZ.enc" -k "$ENCRYPTION_PASSWORD"

# Upload the encrypted backup file to S3
aws s3 cp --endpoint-url https://api-s3.nxbo.ir "$BACKUP_FILE_GZ.enc" "s3://$S3_BUCKET/${DB_NAME}_backup_$DATE.sql.gz.enc"

# Check if the upload was successful
if [ $? -eq 0 ]; then
  echo "Encrypted backup uploaded to S3 successfully: ${DB_NAME}_backup_$DATE.sql.gz.enc"
else
  echo "Upload to S3 failed!"
fi

# Clean up local compressed and encrypted backup files
rm "$BACKUP_FILE_GZ" "$BACKUP_FILE_GZ.enc"
