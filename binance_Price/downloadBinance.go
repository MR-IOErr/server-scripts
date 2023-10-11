package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func getBinanceData(url string) map[string]string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	jsonDataFromHttp, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var binanceJsonData []Records

	err = json.Unmarshal([]byte(jsonDataFromHttp), &binanceJsonData)
	if err != nil {
		panic(err)
	}

	convertBinanceDataToMap := make(map[string]string)
	for _, i := range binanceJsonData {
		convertBinanceDataToMap[i.Symbol] = i.Price
	}

	return convertBinanceDataToMap
}
