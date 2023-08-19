package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func getDataFromNobitexAPI(URL string) []ChartParams {
	var weeklyChart []ChartParams
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	jsonDataFromHttp, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var nobitexMartekCapResponse NobitexMartekCapResponse

	err = json.Unmarshal([]byte(jsonDataFromHttp), &nobitexMartekCapResponse)
	if err != nil {
		log.Fatal(err)
	}

	for _, coinData := range nobitexMartekCapResponse {

		var chartParams ChartParams
		chartParams.Name = strings.ToLower(coinData.Symbol + ".svg")
		chartParams.URL = COINMARKETCAP_API + strconv.Itoa(coinData.Coinmarketcap_Id) + ".svg"

		weeklyChart = append(weeklyChart, chartParams)

	}
	return weeklyChart

}
