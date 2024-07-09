#!/bin/bash

# Function to display usage and exit
usage() {
  echo "Usage: $0 -h <mysql-host> -P <mysql-port> -u <mysql-username> -p <mysql-password> -d <mysql-database> -b1 <s3_bucket_mysql> -b2 <s3_bucket_nginx> -b3 <s3_bucket_www> -k <encryption_password>"
  exit 1
}

# Parse input arguments
while getopts ":h:P:u:p:d:b1:b2:b3:k:" opt; do
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
    b1) S3_BUCKET_MYSQL="$OPTARG"
    ;;
    b2) S3_BUCKET_NGINX="$OPTARG"
    ;;
    b3) S3_BUCKET_WWW="$OPTARG"
    ;;
    k) ENCRYPTION_PASSWORD="$OPTARG"
    ;;
    *) usage
    ;;
  esac
done

# Check if all required variables are set
if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_USER" ] || [ -z "$DB_PASS" ] || [ -z "$DB_NAME" ] || [ -z "$S3_BUCKET_MYSQL" ] || [ -z "$S3_BUCKET_NGINX" ] || [ -z "$S3_BUCKET_WWW" ] || [ -z "$ENCRYPTION_PASSWORD" ]; then
    echo "Missing required arguments."
    usage
fi

# Perform MySQL backup
DATE=$(date +"%Y%m%d")
MYSQL_BACKUP_FILE="/tmp/${DB_NAME}_mysql_backup_$DATE.sql"

# Perform the MySQL backup using mysqldump
mysqldump -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" > "$MYSQL_BACKUP_FILE"

# Check if the backup was successful
if [ $? -ne 0 ]; then
  echo "MySQL dump failed. Exiting."
  exit 1
fi


sync "$MYSQL_BACKUP_FILE"


# Tar MySQL backup
tar czf "${MYSQL_BACKUP_FILE}.tar.gz" "$MYSQL_BACKUP_FILE"

# Encrypt MySQL backup using OpenSSL with the provided encryption password
openssl enc -aes-256-cbc -salt -in "${MYSQL_BACKUP_FILE}.tar.gz" -out "${MYSQL_BACKUP_FILE}.tar.gz.enc" -k "$ENCRYPTION_PASSWORD"

# Upload MySQL encrypted backup file to S3 bucket for MySQL
aws s3 cp --endpoint-url https://api-s3.nxbo.ir "${MYSQL_BACKUP_FILE}.tar.gz.enc" "s3://$S3_BUCKET_MYSQL/${DB_NAME}_mysql_backup_$DATE.sql.tar.gz.enc"

# Check if the upload was successful
if [ $? -eq 0 ]; then
  echo "Encrypted MySQL backup uploaded to S3 successfully: ${DB_NAME}_mysql_backup_$DATE.sql.tar.gz.enc"
else
  echo "Upload of MySQL backup to S3 failed!"
fi

# Backup /var/www/ directory
WWW_BACKUP_DIR="/tmp/www_backup"
mkdir -p "$WWW_BACKUP_DIR"
tar czf "${WWW_BACKUP_DIR}/${HOSTNAME}_www_backup_$DATE.tar.gz" -C /var/www .

# Encrypt /var/www/ backup using OpenSSL with the provided encryption password
openssl enc -aes-256-cbc -salt -in "${WWW_BACKUP_DIR}/${HOSTNAME}_www_backup_$DATE.tar.gz" -out "${WWW_BACKUP_DIR}/${HOSTNAME}_www_backup_$DATE.tar.gz.enc" -k "$ENCRYPTION_PASSWORD"

# Upload /var/www/ encrypted backup file to S3 bucket for WWW
aws s3 cp --endpoint-url https://api-s3.nxbo.ir "${WWW_BACKUP_DIR}/${HOSTNAME}_www_backup_$DATE.tar.gz.enc" "s3://$S3_BUCKET_WWW/${HOSTNAME}_www_backup_$DATE.tar.gz.enc"

# Check if the upload was successful
if [ $? -eq 0 ]; then
  echo "Encrypted /var/www/ backup uploaded to S3 successfully: ${HOSTNAME}_www_backup_$DATE.tar.gz.enc"
else
  echo "Upload of /var/www/ backup to S3 failed!"
fi

# Backup /etc/nginx/ directory
NGINX_BACKUP_DIR="/tmp/nginx_backup"
mkdir -p "$NGINX_BACKUP_DIR"
tar czf "${NGINX_BACKUP_DIR}/${HOSTNAME}_nginx_backup_$DATE.tar.gz" -C /etc/nginx .

# Encrypt /etc/nginx/ backup using OpenSSL with the provided encryption password
openssl enc -aes-256-cbc -salt -in "${NGINX_BACKUP_DIR}/${HOSTNAME}_nginx_backup_$DATE.tar.gz" -out "${NGINX_BACKUP_DIR}/${HOSTNAME}_nginx_backup_$DATE.tar.gz.enc" -k "$ENCRYPTION_PASSWORD"

# Upload /etc/nginx/ encrypted backup file to S3 bucket for NGINX
aws s3 cp --endpoint-url https://api-s3.nxbo.ir "${NGINX_BACKUP_DIR}/${HOSTNAME}_nginx_backup_$DATE.tar.gz.enc" "s3://$S3_BUCKET_NGINX/${HOSTNAME}_nginx_backup_$DATE.tar.gz.enc"

# Check if the upload was successful
if [ $? -eq 0 ]; then
  echo "Encrypted /etc/nginx/ backup uploaded to S3 successfully: ${HOSTNAME}_nginx_backup_$DATE.tar.gz.enc"
else
  echo "Upload of /etc/nginx/ backup to S3 failed!"
fi

# Clean up temporary backup files
rm -rf "$WWW_BACKUP_DIR" "$NGINX_BACKUP_DIR" "${MYSQL_BACKUP_FILE}.tar.gz" "${MYSQL_BACKUP_FILE}.tar.gz.enc" "${WWW_BACKUP_DIR}/${HOSTNAME}_www_backup_$DATE.tar.gz" "${WWW_BACKUP_DIR}/${HOSTNAME}_www_backup_$DATE.tar.gz.enc" "${NGINX_BACKUP_DIR}/${HOSTNAME}_nginx_backup_$DATE.tar.gz" "${NGINX_BACKUP_DIR}/${HOSTNAME}_nginx_backup_$DATE.tar.gz.enc"

echo "Temporary backup files cleaned up."
