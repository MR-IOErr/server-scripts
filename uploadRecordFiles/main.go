package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	srcDir   = "/home/arian/recordings/"
	dstDir   = "/home/arian/record/"
	videoLog = "/home/arian/record/logfile.log"
)

func calculateHashSum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil

}

func checkMetadataExist(dir string) bool {
	metadataPath := filepath.Join(dir, "metadata.json")
	_, err := os.Stat(metadataPath)
	return !os.IsNotExist(err)
}

func copyFiles(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())

}

func main() {
	err := filepath.Walk(srcDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				metaPath := filepath.Dir(path)
				if checkMetadataExist(metaPath) {
					fileName := info.Name()
					fileParts := strings.Split(fileName, "_")
					if len(fileParts) > 1 {
						newFileName := fileParts[0] + filepath.Ext(fileName)
						newPath := filepath.Join(dstDir, newFileName)
						if _, err := os.Stat(newPath); err == nil {
							srcHash, err := calculateHashSum(path)
							if err != nil {
								return err
							}

							dstHash, err := calculateHashSum(newPath)
							if err != nil {
								return err

							}
							if srcHash == dstHash {

								os.RemoveAll(metaPath)
								// fmt.Printf("File already Exist and has been removed, skipping: %s\n", metaPath)
								return nil
							}
						} else {
							err := copyFiles(path, newPath)
							if err != nil {
								return err
							}

							logFile, err := os.OpenFile(videoLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
							if err != nil {
								log.Fatal(err)
							}

							defer logFile.Close()

							log.SetOutput(logFile)

							log.Printf("Copied: %s -> %s\n", path, newPath)

						}
					}
				}

			}
			return err
		})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

}
