package public

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/richardsric/gateway/helper"
)

type OrderResponse struct {
	Result      string `json:"result"`
	Message     string `json:"message"`
	OrderNumber string `json:"order_number"`
}

//BuyLimit is a gateway route to be called by apps when placing buy orders to exchanges
// eg localhost:5000/sellOrder?market=BTC-VTC&quantity=300&rate=0.00019802&eid=1&apiKey=110982d6fd72480d9968cbca3473a868&aid=1
func BuyLimit(w http.ResponseWriter, req *http.Request) {
	fmt.Println("******************************************************************************************************************")
	fmt.Println("GATEWAY: entered Buy Limit order function")
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")

	apiKey := req.FormValue("apiKey")
	market := req.FormValue("market")
	quantityStr := req.FormValue("quantity")
	rateStr := req.FormValue("rate")
	eidStr := req.FormValue("eid")
	accountIdStr := req.FormValue("aid") //required from the worker calling the gateway

	switch {
	case apiKey == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'apiKey' parameter not found.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	case market == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'market' parameter not found.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	case quantityStr == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'quantity' parameter not found.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	case rateStr == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'rate' parameter not found.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	case eidStr == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'eid' parameter not found.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	case accountIdStr == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'aid' parameter not found.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	}
	quantity, _ := strconv.ParseFloat(quantityStr, 64)
	rate, _ := strconv.ParseFloat(rateStr, 64)
	eid, _ := strconv.Atoi(eidStr)
	accountId, _ := strconv.Atoi(accountIdStr)

	apiSecret, err := helper.GetSecret(eid, accountId, apiKey) //req.FormValue("secret")
	if err != nil {
		d := OrderResponse{
			Result:  "error",
			Message: "'apiKey' parameter not correct.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	}
	fmt.Printf("apiKey=%v\napiSecret=%v\nmarket=%v\nquantity=%v\nrate=%v\nexchange_id=%v\naccount_id=%v\n", apiKey, apiSecret, market, quantity, rate, eid, accountId)

	//	Use switch statement to switch through the exchange_id and call the appropriate object required
	switch eid {
	case 1:
		queryString := fmt.Sprintf("http://localhost:5030/returnOrderBuy?market=%s&quantity=%v&rate=%v&apiKey=%s&secret=%s",
			market, quantity, rate, apiKey, apiSecret)
		fmt.Println("buyOrder query string is ", queryString)
		resp, err := http.Get(queryString)
		if err != nil {
			d := OrderResponse{
				Result:  "error",
				Message: "No connection could be made because the target machine actively refused it.",
			}
			bs, _ := json.Marshal(d)
			fmt.Println(string(bs))
			fmt.Fprintln(w, string(bs))
			return
		}
		defer resp.Body.Close()
		jsonResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			d := OrderResponse{
				Result:  "error",
				Message: err.Error(),
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		fmt.Fprintln(w, string(jsonResp))
		fmt.Println(string(jsonResp))
		fmt.Println("##############################")
		valid, orderId := helper.IsRequestValid(jsonResp)
		if !valid {
			fmt.Println("Buy limit request submitted is not valid")
			return
		}
		fmt.Println("Buy limit request submitted is valid")
		//Selecting the txnFee from the db of that particular exchange
		txnFee, err := helper.GetTxnFee(eid)
		if err != nil {
			fmt.Println("couldn't fetch the txnFee due to ", err)
		}
		fmt.Println("The txnFee of the exchange fetched is ", txnFee)
		fmt.Printf("inserting: accountId-%v,exchangeId-%v,quantity-%v,rate-%v,market-%v,orderId-%v,txnFee-%v\n",
			accountId, eid, quantity, rate, market, orderId, txnFee)
		err = helper.DoInsert(accountId, eid, quantity, rate, txnFee, market, "BUY", orderId)
		if err != nil {
			fmt.Println(err)
		}
	default:
		d := OrderResponse{
			Result:  "error",
			Message: "Sorry!, Other exchanges are not yet supported.",
		}
		bs, _ := json.Marshal(d)
		fmt.Println(string(bs))
		fmt.Fprintln(w, string(bs))
	}

}

//SellLimit is a gateway route to be called by apps when placing selling orders to exchanges
// eg localhost:5000/buyOrder?market=BTC-VTC&quantity=300&rate=0.00019802&eid=1&apiKey=110982d6fd72480d9968cbca3473a868&aid=1
func SellLimit(w http.ResponseWriter, req *http.Request) {
	fmt.Println("******************************************************************************************************************")
	fmt.Println("GATEWAY: entered Sell Limit order function")
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")
	apiKey := req.FormValue("apiKey")
	accountIdStr := req.FormValue("aid") //required from the worker calling the gateway
	market := req.FormValue("market")
	quantityStr := req.FormValue("quantity")
	rateStr := req.FormValue("rate")
	eidStr := req.FormValue("eid")

	switch {
	case apiKey == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'apiKey' parameter not found.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	case market == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'market' parameter not found.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	case quantityStr == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'quantity' parameter not found.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	case rateStr == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'rate' parameter not found.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	case eidStr == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'eid' parameter not found.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	case accountIdStr == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'aid' parameter not found.",
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
	eid, _ := strconv.Atoi(eidStr)
	accountId, _ := strconv.Atoi(accountIdStr)

	apiSecret, err := helper.GetSecret(eid, accountId, apiKey) //req.FormValue("secret")
	if err != nil {
		d := OrderResponse{
			Result:  "error",
			Message: "'apiKey' parameter not correct.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	}
	fmt.Printf("apiKey=%v\napiSecret=%v\nmarket=%v\nquantity=%v\nrate=%v\nexchange_id=%v\naccount_id=%v\n", apiKey, apiSecret, market, quantity, rateStr, eid, accountId)

	//	Use switch statement to switch the exchange_id and call the appropriate object required
	//	instanciats the particulat exchange to send the request to
	switch eid {
	case 1:
		queryString := fmt.Sprintf("http://localhost:5030/returnOrderSell?market=%s&quantity=%v&rate=%v&apiKey=%s&secret=%s",
			market, quantity, rate, apiKey, apiSecret)
		fmt.Println("sellOrder query string is ", queryString)
		resp, err := http.Get(queryString)
		if err != nil {
			d := OrderResponse{
				Result:  "error",
				Message: "No connection could be made because the target machine actively refused it.",
			}
			bs, _ := json.Marshal(d)
			fmt.Println(string(bs))
			fmt.Fprintln(w, string(bs))
		}
		defer resp.Body.Close()
		jsonResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			d := OrderResponse{
				Result:  "error",
				Message: err.Error(),
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		fmt.Fprintln(w, string(jsonResp))
		fmt.Println(string(jsonResp))
		fmt.Println("##############################")
		valid, orderId := helper.IsRequestValid(jsonResp)
		if !valid {
			fmt.Println("Sell limit request submitted is not valid")
			return
		}
		fmt.Println("Sell limit request submitted is valid")
		//Selecting the txnFee from the db of that particular exchange
		txnFee, err := helper.GetTxnFee(eid)
		if err != nil {
			fmt.Println("couldn't fetch the txnFee due to ", err)
		}
		fmt.Println("The txnFee of the exchange fetched is ", txnFee)
		fmt.Printf("inserting: accountId-%v,exchangeId-%v,quantity-%v,rate-%v,market-%v,orderId-%v,txnFee-%v\n",
			accountId, eid, quantity, rate, market, orderId, txnFee)
		err = helper.DoInsert(accountId, eid, quantity, rate, txnFee, market, "SELL", orderId)
		if err != nil {
			fmt.Println(err)
		}
	default:
		d := OrderResponse{
			Result:  "error",
			Message: "Sorry!, Other exchanges are not yet supported.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
	}
	fmt.Println("******************************************************************************************************************")

}

//CancelOrder is a gateway route to be called by apps when canceling orders placed in exchanges
//eg localhost:5000/cancelOrder?uuid=9b68ff76-a447-4aa2-8032-af3bb6fb7046&eid=1&aid=1&apiKey=110982d6fd72480d9968cbca3473a868
func CancelOrder(w http.ResponseWriter, req *http.Request) {
	fmt.Println("******************************************************************************************************************")
	fmt.Println("GATEWAY: entered cancel order function")
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")
	apiKey := req.FormValue("apiKey")
	accountIdStr := req.FormValue("aid") //required from the worker calling the gateway

	orderNumber := req.FormValue("uuid")
	eidStr := req.FormValue("eid")

	switch {
	case apiKey == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'apiKey' parameter not found.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	case orderNumber == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'uuid' parameter not found.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	case eidStr == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'eid' parameter not found.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	case accountIdStr == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'aid' parameter not found.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	}

	eid, _ := strconv.Atoi(eidStr)
	accountId, _ := strconv.Atoi(accountIdStr)

	apiSecret, err := helper.GetSecret(eid, accountId, apiKey) //req.FormValue("secret")
	if err != nil {
		d := OrderResponse{
			Result:  "error",
			Message: "'apiKey' parameter not correct.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	}
	fmt.Printf("apiKey=%v\napiSecret=%v\norder number=%v\nexchange_id=%v\naccount_id=%v\n", apiKey, apiSecret, orderNumber, eid, accountId)
	switch eid {
	case 1:
		queryString := fmt.Sprintf("http://localhost:5030/returnOrderCancel?uuid=%s&apiKey=%s&secret=%s",
			orderNumber, apiKey, apiSecret)
		fmt.Println("cancelOrder query string is ", queryString)
		resp, err := http.Get(queryString)
		if err != nil {
			d := OrderResponse{
				Result:  "error",
				Message: "No connection could be made because the target machine actively refused it.",
			}
			bs, _ := json.Marshal(d)
			fmt.Println(string(bs))
			fmt.Fprintln(w, string(bs))
		}
		defer resp.Body.Close()
		jsonResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			d := OrderResponse{
				Result:  "error",
				Message: err.Error(),
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		fmt.Fprintln(w, string(jsonResp))
		fmt.Println(string(jsonResp))
	default:
		d := OrderResponse{
			Result:  "error",
			Message: "Sorry!, Other exchanges are not yet supported.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
	}
	fmt.Println("******************************************************************************************************************")

}

//GetOrderInfo is a gateway route to be called by apps to get information of an order (SELL OR BUY)
//eg localhost:5000/getOrderInfo?apiKey=110982d6fd72480d9968cbca3473a868&uuid=34a42ddc-22b5-493d-a42b-4ddf88ef9ed8&eid=1&aid=1
func GetOrderInfo(w http.ResponseWriter, req *http.Request) {
	fmt.Println("******************************************************************************************************************")
	fmt.Println("GATEWAY: entered get order info function")
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")
	apiKey := req.FormValue("apiKey")
	accountIdStr := req.FormValue("aid") //required from the worker calling the gateway

	orderNumber := req.FormValue("uuid")
	eidStr := req.FormValue("eid")

	switch {
	case apiKey == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'apiKey' parameter not found.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	case eidStr == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'eid' parameter not found.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	case accountIdStr == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'aid' parameter not found.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	case orderNumber == "":
		d := OrderResponse{
			Result:  "error",
			Message: "'uuid' parameter not found.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	}
	eid, _ := strconv.Atoi(eidStr)
	accountId, _ := strconv.Atoi(accountIdStr)

	//this line gets the secret using apiKey,aid,eid from the database
	apiSecret, err := helper.GetSecret(eid, accountId, apiKey)
	if err != nil {
		d := OrderResponse{
			Result:  "error",
			Message: "'apiKey' parameter not correct.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
		return
	}
	fmt.Printf("apiKey=%v\napiSecret=%v\nuuid=%v\nexchange_id=%v\naccount_id=%v\n", apiKey, apiSecret, orderNumber, eid, accountId)
	switch eid {
	case 1:
		queryString := fmt.Sprintf("http://localhost:5030/returnOrderInfo?uuid=%s&apiKey=%s&secret=%s",
			orderNumber, apiKey, apiSecret)
		fmt.Println("getOrderInfo: query string is, ", queryString)
		resp, err := http.Get(queryString)
		if err != nil {
			d := OrderResponse{
				Result:  "error",
				Message: "No connection could be made because the target machine actively refused it.",
			}
			bs, _ := json.Marshal(d)
			fmt.Println(string(bs))
			fmt.Fprintln(w, string(bs))
			return
		}
		defer resp.Body.Close()
		jsonResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			d := OrderResponse{
				Result:  "error",
				Message: err.Error(),
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		fmt.Fprintln(w, string(jsonResp))
		fmt.Println(string(jsonResp))
	default:
		d := OrderResponse{
			Result:  "error",
			Message: "Sorry!, Other exchanges are not yet supported.",
		}
		bs, _ := json.Marshal(d)
		fmt.Fprintln(w, string(bs))
	}
	fmt.Println("******************************************************************************************************************")

}
