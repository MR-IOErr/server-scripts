package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func getDataFromCoinMarketAPI(file ChartParams, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(file.URL)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	// data, err := os.Create(localPATH + file.Name)
	// io.Copy(data, resp.Body)
	// if err != nil {
	// 	fmt.Printf("Error reading response body: %s\n", err)
	// }

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
	}
	err = os.WriteFile(localPATH+file.Name, data, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %s\n", err)
	}

}
