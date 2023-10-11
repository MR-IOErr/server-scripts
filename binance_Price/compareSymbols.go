package main

import (
	"encoding/json"
	"os"
)

func compareSymbols(binanceData, okxData map[string]string) {
	for okxKey, okxValue := range okxData {
		_, ok := binanceData[okxKey]
		if ok {
			continue
		} else {
			binanceData[okxKey] = okxValue
		}
	}

	finalData := make([]Records, 0)
	for binanceKey, binanceValue := range binanceData {
		finalData = append(finalData, Records{
			Symbol: binanceKey,
			Price:  binanceValue,
		})
	}

	result, _ := json.Marshal(finalData)
	os.WriteFile(RECORD, result, 0644)
}
