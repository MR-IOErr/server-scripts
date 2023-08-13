package main

type Records struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func main() {
	binanceAPI := "https://api.binance.com/api/v3/ticker/price"

	getData(binanceAPI)
	uploadFileToArvanS3()
	deleteDownloadedFile(RECORD)
}
