package main

import (
	"fmt"
	"net/http"

	"github.com/richardsric/microservices/bittrexmicro/bittrex"
	"github.com/richardsric/microservices/bittrexmicro/bittrex/public"
)

func main() {

	http.HandleFunc("/returnOrderBuy", bittrex.BuyLimit)
	http.HandleFunc("/returnOrderSell", bittrex.SellLimit)
	http.HandleFunc("/returnOrderCancel", bittrex.CancelOrder)
	http.HandleFunc("/returnOrderInfo", bittrex.GetOrderInfo)
	http.HandleFunc("/returnBalances", bittrex.GetBalances)
	http.HandleFunc("/returnNonZeroBalances", bittrex.GetNonZeroBalances)
	http.HandleFunc("/returnBalance", bittrex.GetBalance)
	http.HandleFunc("/stat", public.Statz)
	http.HandleFunc("/pair/price", public.BittrexSinglePair)
	http.HandleFunc("/bittrex_ticker", public.BittrexMarketData1)

	http.ListenAndServe(":5030", nil)
}

func init() {
	var name = "iTradeCoin Bitt-MicroService"
	var version = "0.001 DEVEL"
	var developer = "iYOCHU Nig LTD"

	fmt.Println("App Name: ", name)
	fmt.Println("App Version: ", version)
	fmt.Println("Developer Name: ", developer)

	fmt.Println("iTradeCoin Bittrex Microservice: Running on Port 'https://localhost:5030/'")
	go public.MarketDataService()
	go public.ClearStatData()
}
