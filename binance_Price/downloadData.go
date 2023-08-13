package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func getData(url string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	jsonDataFromHttp, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var jsonData []Records

	err = json.Unmarshal([]byte(jsonDataFromHttp), &jsonData)
	if err != nil {
		panic(err)
	}

	os.WriteFile(RECORD, []byte(jsonDataFromHttp), 0644)
}
