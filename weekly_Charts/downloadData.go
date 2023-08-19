package main

import (
	"os"
	"sync"
)

func download(weeklyChart []ChartParams) {
	var wg sync.WaitGroup

	os.Mkdir(localPATH, 0744)

	for _, file := range weeklyChart {
		wg.Add(1)
		go getDataFromCoinMarketAPI(file, &wg)
	}

	wg.Wait()

}
