#!/bin/bash

# Load environment variables from .env file
if [ -f /root/script/full-backup/.env ]; then
  export $(cat /root/script/full-backup/.env | grep -v '#' | awk '/=/ {print $1}')
fi

# Function to backup MySQL database
backup_mysql() {
  local db_name=$1
  local db_user=$2
  local db_password=$3
  local bucket_name=$4
  local encryption_password=$5
  local aws_access_key_id=$6
  local aws_secret_access_key=$7
  local timestamp=$(date +%Y%m%d)
  local backup_filename="${db_name}_backup_$timestamp"

  echo "Starting MySQL backup for $db_name..."
  mysqldump -h $MYSQL_HOST -u $db_user -p$db_password $db_name > /tmp/$backup_filename.sql
  tar -cf /tmp/$backup_filename.tar -C /tmp $backup_filename.sql
  xz /tmp/$backup_filename.tar
  openssl enc -aes-256-cbc -salt -in /tmp/$backup_filename.tar.xz -out /tmp/$backup_filename.enc -k $encryption_password

  aws configure set aws_access_key_id $aws_access_key_id
  aws configure set aws_secret_access_key $aws_secret_access_key

  aws s3 cp --endpoint-url https://s3 /tmp/$backup_filename.enc s3://$bucket_name/$backup_filename.enc

  rm /tmp/$backup_filename.sql /tmp/$backup_filename.tar.xz /tmp/$backup_filename.enc
  echo "MySQL backup for $db_name completed."
}

# Function to backup directories
backup_directory() {
  local src_dir=$1
  local bucket_name=$2
  local encryption_password=$3
  local aws_access_key_id=$4
  local aws_secret_access_key=$5
  local prefix=$6
  local dir_name=$(basename $src_dir)
  local timestamp=$(date +%Y%m%d)
  local backup_filename="${prefix}_${dir_name}_backup_$timestamp"

  echo "Starting backup of $src_dir..."
  tar -cf /tmp/$backup_filename.tar -C $(dirname $src_dir) $(basename $src_dir)
  xz /tmp/$backup_filename.tar
  openssl enc -aes-256-cbc -salt -in /tmp/$backup_filename.tar.xz -out /tmp/$backup_filename.enc -k $encryption_password

  aws configure set aws_access_key_id $aws_access_key_id
  aws configure set aws_secret_access_key $aws_secret_access_key

  aws s3 cp --endpoint-url https://s3 /tmp/$backup_filename.enc s3://$bucket_name/$backup_filename.enc

  rm /tmp/$backup_filename.tar.xz /tmp/$backup_filename.enc
  echo "Backup of $src_dir completed."
}

# MySQL backups

backup_mysql $HELP_DB_NAME $HELP_DB_USER $HELP_DB_PASSWORD $HELP_BUCKET $HELP_ENCRYPTION_PASSWORD $HELP_AWS_ACCESS_KEY_ID $HELP_AWS_SECRET_ACCESS_KEY
backup_mysql $MAG_DB_NAME $MAG_DB_USER $MAG_DB_PASSWORD $MAG_BUCKET $MAG_ENCRYPTION_PASSWORD $MAG_AWS_ACCESS_KEY_ID $MAG_AWS_SECRET_ACCESS_KEY
backup_mysql $ACADEMY_DB_NAME $ACADEMY_DB_USER $ACADEMY_DB_PASSWORD $ACADEMY_BUCKET $ACADEMY_ENCRYPTION_PASSWORD $ACADEMY_AWS_ACCESS_KEY_ID $ACADEMY_AWS_SECRET_ACCESS_KEY

# Nginx backups

backup_directory "/etc/nginx/sites-enabled/academy" $NGINX_ACADEMY_BUCKET $NGINX_ENCRYPTION_PASSWORD $NGINX_ACADEMY_AWS_ACCESS_KEY_ID $NGINX_ACADEMY_AWS_SECRET_ACCESS_KEY "nginx"
backup_directory "/etc/nginx/sites-enabled/help" $NGINX_HELP_BUCKET $NGINX_ENCRYPTION_PASSWORD $NGINX_HELP_AWS_ACCESS_KEY_ID $NGINX_HELP_AWS_SECRET_ACCESS_KEY "nginx"

# HTML backups

backup_directory "/var/www/html/help" $HTML_HELP_BUCKET $HTML_ENCRYPTION_PASSWORD $HTML_HELP_AWS_ACCESS_KEY_ID $HTML_HELP_AWS_SECRET_ACCESS_KEY "html"
backup_directory "/var/www/html/mag" $HTML_MAG_BUCKET $HTML_ENCRYPTION_PASSWORD $HTML_MAG_AWS_ACCESS_KEY_ID $HTML_MAG_AWS_SECRET_ACCESS_KEY "html"

echo "All backups completed."
