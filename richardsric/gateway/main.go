package main

import (
	"fmt"

	"net/http"

	"github.com/richardsric/gateway/public"
)

func main() {
	http.HandleFunc("/buyOrder", public.BuyLimit)
	http.HandleFunc("/sellOrder", public.SellLimit)
	http.HandleFunc("/cancelOrder", public.CancelOrder)
	http.HandleFunc("/getOrderInfo", public.GetOrderInfo)
	http.HandleFunc("/getNonZeroBalances", public.GetNonZeroBalances)
	http.HandleFunc("/getBalances", public.GetBalances)
	http.HandleFunc("/getBalance", public.GetBalance)
	http.HandleFunc("/pair/price", public.GetWaySinglePair)
	http.HandleFunc("/websocket", public.Websocket)
	http.HandleFunc("/test", public.Test)
	http.HandleFunc("/ticker", public.GetAskBid)

	//start the service report mux on the port specified for it in service_report port: range 6000
	//	fmt.Println("Started Service Reprorting on port 6001")
	//	go http.ListenAndServe(":6001", nil)

	//start the main mux on the port specified in the main service port from settings table
	fmt.Println("Started Main Service on port 5000")
	http.ListenAndServe(":5000", nil)

}

func init() {

	//Load Application Settings here

	// exchange api settings here and export it as read only.
	var name = "iTradeCoin GateWay"
	var version = "0.001 DEVEL"
	var developer = "iYOCHU Nig LTD"

	fmt.Println("App Name: ", name)
	fmt.Println("App Version: ", version)
	fmt.Println("Developer Name: ", developer)

	//	fmt.Println("iTradeCoin Order GateWay: Visit 'http://localhost:5000/' to use")

	//bittrex.ApiSettings//.LoadSettings(API_BASE, API_VERSION, DEFAULT_HTTPCLIENT_TIMEOUT)
	//microservices.Settings()
	go public.TruncateMarketData()

}
