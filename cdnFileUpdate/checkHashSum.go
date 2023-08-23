package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-redis/redis"
)

func checkHash() {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis.localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var fileMap = make(map[string]string)

	walker := make(fileWalk)

	go func() {
		// Gather the files to upload by walking the path recursively
		if err := filepath.Walk(localPATH, walker.Walk); err != nil {
			log.Fatalln("Walk failed:", err)
		}
		close(walker)
	}()

	for filePath := range walker {
		fileName := strings.TrimPrefix(filePath, localPATH)
		file, _ := os.ReadFile(filePath)
		hashSUM := sha256.Sum256([]byte(file))
		fileMap[fileName] = hex.EncodeToString(hashSUM[:])
		hash := hex.EncodeToString(hashSUM[:])

		exist, _ := client.Exists("CDNFile:::" + fileName).Result()

		if exist == 1 {
			val, _ := client.Get("CDNFile:::" + fileName).Result()
			if val == hash {
				fmt.Println(fileName, "Exists on redis")
				continue
			} else {
				fmt.Println(fileName, "Updated on redis")
				client.Set("CDNFile:::"+fileName, hash, 0).Err()
			}
		} else {
			fmt.Println(fileName, "Added to redis")
			client.Set("CDNFile:::"+fileName, hash, 0).Err()
		}
		f, err := os.Open(filePath)
		if err != nil {
			log.Println("Failed to open file: ", filePath, err)
			continue
		}

		uploadToS3(f, fileName)

	}

}
