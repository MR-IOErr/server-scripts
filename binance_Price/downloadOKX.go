package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func getOKXData(url string) map[string]string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	jsonDataFromHttp, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var okxJsonData map[string]interface{}

	err = json.Unmarshal([]byte(jsonDataFromHttp), &okxJsonData)
	if err != nil {
		panic(err)
	}

	var okxData []Records

	okxRecords := okxJsonData["data"].([]interface{})

	for _, j := range okxRecords {
		k := j.(map[string]interface{})
		okxData = append(okxData, Records{Symbol: strings.Replace(k["instId"].(string), "-", "", -1), Price: k["last"].(string)})
	}

	convertOkxDataToMap := make(map[string]string)
	for _, i := range okxData {
		convertOkxDataToMap[i.Symbol] = i.Price
	}
	return convertOkxDataToMap

}
