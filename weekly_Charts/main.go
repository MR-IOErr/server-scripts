package main

import (
	"fmt"
	"time"
)

type Market struct {
	Name                 string `json:"name"`
	ID                   int    `json:"id"`
	Symbol               string `json:"symbol"`
	Price                string `json:"price"`
	Coinmarketcap_Id     int    `json:"coinmarketcap_id"`
	Coinmarketcap_Symbol string `json:"coinmarketcap_symbol"`
}

type NobitexMartekCapResponse map[string]Market

type ChartParams struct {
	Name string
	URL  string
}

const (
	NOBITEX_API_URL   = "https://api.nobitex.ir/coinmarketcap/v1/ids"
	COINMARKETCAP_API = "https://s3.coinmarketcap.com/generated/sparklines/web/7d/2781/"
)

func main() {
	startTime := time.Now()

	weeklyChart := getDataFromNobitexAPI(NOBITEX_API_URL)

	download(weeklyChart)
	upload(weeklyChart)
	deleteFiles()

	fmt.Println(time.Since(startTime))

}
