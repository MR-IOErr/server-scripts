package main

type Records struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

const (
	AWS_Access_Key = ""
	AWS_Secret_key = ""
	AWS_Bucket     = "nobitex-cdn"
	AWS_S3_URL     = "s3.ir-thr-at1.arvanstorage.ir"
	RECORD         = "binance.json"
	DST            = "/prices/binance.json"
	PER            = "public-read"
)

func main() {
	binanceAPI := "https://api.binance.com/api/v3/ticker/price"
	okxAPI := "https://www.okx.com/api/v5/market/tickers?instType=SPOT"

	binance := getBinanceData(binanceAPI)
	okx := getOKXData(okxAPI)

	compareSymbols(binance, okx)
	uploadFileToArvanS3()
	deleteDownloadedFile(RECORD)
}
