package main

import "os"

func deleteDownloadedFile(file string) {
	os.Remove(file)
}
