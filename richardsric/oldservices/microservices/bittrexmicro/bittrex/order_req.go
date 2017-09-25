package bittrex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/richardsric/microservices/bittrexmicro/helper"
)

func BuyLimit(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Microservice: entered Buy Limit order function")
	if req.Method == "GET" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add("Accept", "application/json")

		apiKey := req.FormValue("apiKey")
		apiSecret := req.FormValue("secret")

		market := req.FormValue("market")
		quantityStr := req.FormValue("quantity")
		rateStr := req.FormValue("rate")
		switch {
		case apiKey == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'apiKey' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case apiSecret == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'secret' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case market == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'market' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case quantityStr == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'quantity' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case rateStr == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'rate' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		quantity, _ := strconv.ParseFloat(quantityStr, 64)
		rate, _ := strconv.ParseFloat(rateStr, 64)
		fmt.Printf("apiKey=%v\napiSecret=%v\nmarket=%v\nquantity=%v\nrate=%v\n", apiKey, apiSecret, market, quantity, rate)

		b := New(apiKey, apiSecret)
		jsonResp := b.BuyLimit(market, quantity, rate)
		fmt.Fprintln(w, jsonResp)
		fmt.Println(jsonResp)
	} else {
		fmt.Println("The method shouldn't be POST method")
	}
}

func SellLimit(w http.ResponseWriter, req *http.Request) {
	fmt.Println("******************************************************************************************************************")
	fmt.Println("Microservice: entered Sell Limit order function")
	if req.Method == "GET" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add("Accept", "application/json")
		apiKey := req.FormValue("apiKey")
		apiSecret := req.FormValue("secret")

		market := req.FormValue("market")
		quantityStr := req.FormValue("quantity")
		rateStr := req.FormValue("rate")

		switch {
		case apiKey == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'apiKey' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case apiSecret == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'secret' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case market == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'market' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case quantityStr == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'quantity' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case rateStr == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'rate' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		quantity, err := strconv.ParseFloat(quantityStr, 64)
		if err != nil {
			fmt.Println(err)
		}
		rate, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("apiKey=%v\napiSecret=%v\nmarket=%v\nquantity=%v\nrate=%v\n", apiKey, apiSecret, market, quantity, rateStr)

		b := New(apiKey, apiSecret)
		jsonResp := b.SellLimit(market, quantity, rate)
		fmt.Println(jsonResp)
		fmt.Fprintln(w, jsonResp)
		fmt.Println("******************************************************************************************************************")
	} else {
		fmt.Println("The method shouldn't be POST method")
	}
}

func CancelOrder(w http.ResponseWriter, req *http.Request) {
	fmt.Println("******************************************************************************************************************")
	fmt.Println("Microservice: entered cancel order function")
	if req.Method == "GET" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add("Accept", "application/json")
		apiKey := req.FormValue("apiKey")
		apiSecret := req.FormValue("secret")

		orderId := req.FormValue("uuid")
		switch {
		case apiKey == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'apiKey' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case apiSecret == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'secret' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		fmt.Printf("apiKey=%v\napiSecret=%v\norder number=%v\n", apiKey, apiSecret, orderId)
		b := New(apiKey, apiSecret)
		res := b.CancelOrder(orderId)
		fmt.Println(res)
		fmt.Fprintln(w, res)
		fmt.Println("******************************************************************************************************************")
	} else {
		fmt.Println("The method shouldn't be POST method")
	}
}

/*
func GetOrderHistory(w http.ResponseWriter, req *http.Request) {
	fmt.Println("******************************************************************************************************************")
  fmt.Println("entered get orders function")
  if req.Method == "POST" {
    req.Header.Add("Content-Type", "application/json;charset=utf-8")
    req.Header.Add("Accept", "application/json")
		apiKey := req.FormValue("apiKey")
		apiSecret := req.FormValue("secret")
		accountIdStr := req.FormValue("aid") //required from the worker calling the gateway

    market := req.FormValue("market")
		eidStr := req.FormValue("eid")
		eid, _ := strconv.Atoi(eidStr)
		accountId, _ := strconv.Atoi(accountIdStr)
		fmt.Printf("apiKey=%v\napiSecret=%v\nmarket=%v\nexchange_id=%v\naccount_id=%v",apiKey,apiSecret,market,eid,accountId)
		switch {
		case apiKey == "":
			d := bittrex.OrderResponse{
				Result:"error",
				Message: "'apiKey' parameter not found.",
			}
			bs,_:= json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case apiSecret == "":
			d := bittrex.OrderResponse{
				Result:"error",
				Message: "'secret' parameter not found.",
			}
			bs,_:= json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case eidStr == "":
			d := bittrex.OrderResponse{
				Result:"error",
				Message: "'eid' parameter not found.",
			}
			bs,_:= json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case market == "":
			d := bittrex.OrderResponse{
				Result:"error",
				Message: "'market' parameter can not be empty (input either 'all' or a market).",
			}
			bs,_:= json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case accountIdStr == "":
			d := bittrex.OrderResponse{
				Result:"error",
				Message: "'aid' parameter not found.",
			}
			bs,_:= json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		switch eid {
		case 1:
		//	instanciats the particulat exchange to send the request to
			b := bittrex.New(apiKey, apiSecret)
			res,err := b.GetOrderHistory(market)
			if err != "" {
				fmt.Fprintln(w,res)
			}
			fmt.Println(res,err)
			fmt.Fprintln(w,res)
		default:
			d := bittrex.OrderResponse{
				Result:"error",
				Message: "Sorry!, Other exchanges are not yet supported.",
			}
			bs,_:= json.Marshal(d)
			fmt.Fprintln(w, string(bs))
		}
		fmt.Println("******************************************************************************************************************")
  }
}

func GetOpenOrders(w http.ResponseWriter, req *http.Request) {
	fmt.Println("******************************************************************************************************************")
	fmt.Println("entered getOpenOrders function")
  if req.Method == "POST" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
    req.Header.Add("Accept", "application/json")
		apiKey := req.FormValue("apiKey")
		apiSecret := req.FormValue("secret")
		market := req.FormValue("market")
		accountIdStr := req.FormValue("aid") //required from the worker calling the gateway
		eidStr := req.FormValue("eid")
		eid, _ := strconv.Atoi(eidStr)
		accountId, _ := strconv.Atoi(accountIdStr)
		if apiKey == "" {
			d := bittrex.OrderResponse{
				Result:"error",
				Message:  "'apiKey' parameter not found.",
			}
			bs,_:= json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		if apiSecret == "" {
			d := bittrex.OrderResponse {
				Result:"error",
				Message:  "'secret' parameter not found.",
			}
			bs,_:= json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		if market == "" {
			d := bittrex.OrderResponse {
				Result:"error",
				Message:  "'market' parameter can not be empty (input either 'all' or a market).",
			}
			bs,_:= json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		if accountIdStr == "" {
			d := bittrex.OrderResponse {
				Result:"error",
				Message:  "'aid' parameter can not be empty.",
			}
			bs,_:= json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		if eidStr == "" {
			d := bittrex.OrderResponse {
				Result:"error",
				Message:  "'eid' parameter can not be empty.",
			}
			bs,_:= json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		fmt.Printf("apiKey=%v\napiSecret=%v\nmarket=%v\nexchange_id=%v\naccount_id=%v\neid=%v",apiKey,apiSecret,market,eid,accountId,eid)
		switch eid {
		case 1:
			//	instanciats the particulat exchange to send the request to
			b := bittrex.New(apiKey, apiSecret)
			res,errStr := b.GetOpenOrders(market)
			if errStr != "" {
				fmt.Fprintln(w,errStr)
			}
			fmt.Fprintln(w,res)
			//fmt.Fprintln(w,res)
		default:
			d := bittrex.OrderResponse {
				Result:"error",
				Message: "Sorry!, Other exchanges are not yet supported.",
			}
			bs,_:= json.Marshal(d)
			fmt.Fprintln(w, string(bs))
		}
		fmt.Println("******************************************************************************************************************")
	}
}
*/
func GetOrderInfo(w http.ResponseWriter, req *http.Request) {
	fmt.Println("******************************************************************************************************************")
	fmt.Println("Microservice: entered get order info function")
	if req.Method == "GET" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add("Accept", "application/json")
		apiKey := req.FormValue("apiKey")
		apiSecret := req.FormValue("secret")
		uuid := req.FormValue("uuid")
		switch {
		case apiKey == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'apiKey' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case apiSecret == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'secret' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		fmt.Printf("apiKey=%v\napiSecret=%v\nuuid=%v\n", apiKey, apiSecret, uuid)
		b := New(apiKey, apiSecret)
		jsonResp := b.GetOrderInfo(uuid)
		fmt.Println(jsonResp)
		fmt.Fprintln(w, jsonResp)
		fmt.Println("******************************************************************************************************************")
	} else {
		fmt.Println("The method shouldn't be POST method")
	}
}
