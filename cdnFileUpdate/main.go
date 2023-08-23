package main

const (
	AWS_Access_Key = ""
	AWS_Secret_Key = ""
	AWS_Bucket     = ""
	AWS_S3_URL     = "s3.ir-thr-at1.arvanstorage.ir"
	Permission     = "public-read"
	Region         = "default"
	localPATH      = "/var/www/cdn/"
	//destPATH must be empty for our scenario,
	//but if you need to upload the files in a specific path, put the address in destPATH. for e.x: destPAHT = "/example/dir/"
	destPATH = ""
)

func main() {
	checkHash()
}
