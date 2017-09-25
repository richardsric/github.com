package public

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/richardsric/gateway/helper"

	"strconv"

	"io/ioutil"
)

type BalancesResponse struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

func GetBalances(w http.ResponseWriter, req *http.Request) {
	fmt.Println("******************************************************************************************************************")
	fmt.Println("GATEWAY: entered GetBalances function")
	if req.Method == "GET" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add("Accept", "application/json")

		apiKey := req.FormValue("apiKey")
		eidStr := req.FormValue("eid")
		accountIdStr := req.FormValue("aid")
		switch {
		case apiKey == "":
			d := BalancesResponse{
				Result:  "error",
				Message: "'apiKey' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case eidStr == "":
			d := BalancesResponse{
				Result:  "error",
				Message: "'eid' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case accountIdStr == "":
			d := BalancesResponse{
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
			d := BalancesResponse{
				Result:  "error",
				Message: "Secret not found, 'apiKey' parameter not correct.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		fmt.Printf("apiKey=%v\napiSecret=%v\nexchange_id=%v\naccount_id=%v\n", apiKey, apiSecret, eid, accountId)
		//Use switch statement to switch the exchange_id and make request to the appropriate exchange
		switch eid {
		case 1:
			queryString := fmt.Sprintf("http://localhost:5030/returnBalances?apiKey=%s&secret=%s", apiKey, apiSecret)
			fmt.Println("buyOrder query string is ", queryString)
			resp, err := http.Get(queryString)
			if err != nil {
				d := BalancesResponse{
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
				d := BalancesResponse{
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
			d := BalancesResponse{
				Result:  "error",
				Message: "Sorry!, Other exchanges are not yet supported.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
		}
	}
}

func GetBalance(w http.ResponseWriter, req *http.Request) {
	fmt.Println("******************************************************************************************************************")
	fmt.Println("GATEWAY: entered GetBalance function")
	if req.Method == "GET" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add("Accept", "application/json")

		apiKey := req.FormValue("apiKey")
		eidStr := req.FormValue("eid")
		accountIdStr := req.FormValue("aid")
		currency := req.FormValue("currency")
		switch {
		case apiKey == "":
			d := BalancesResponse{
				Result:  "error",
				Message: "'apiKey' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case currency == "":
			d := BalancesResponse{
				Result:  "error",
				Message: "'currency' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case eidStr == "":
			d := BalancesResponse{
				Result:  "error",
				Message: "'eid' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case accountIdStr == "":
			d := BalancesResponse{
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
			d := BalancesResponse{
				Result:  "error",
				Message: "Secret not found, 'apiKey' parameter not correct.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		fmt.Printf("apiKey=%v\napiSecret=%v\nexchange_id=%v\naccount_id=%v\n", apiKey, apiSecret, eid, accountId)
		//Use switch statement to switch the exchange_id and make request to the appropriate exchange
		switch eid {
		case 1:
			queryString := fmt.Sprintf("http://localhost:5030/returnBalance?apiKey=%s&secret=%s&currency=%s", apiKey, apiSecret, currency)
			fmt.Println("buyOrder query string is ", queryString)
			resp, err := http.Get(queryString)
			if err != nil {
				d := BalancesResponse{
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
				d := BalancesResponse{
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
			d := BalancesResponse{
				Result:  "error",
				Message: "Sorry!, Other exchanges are not yet supported.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
		}
	}
}

func GetNonZeroBalances(w http.ResponseWriter, req *http.Request) {
	fmt.Println("******************************************************************************************************************")
	fmt.Println("GATEWAY: entered GetNonZeroBalances function")
	if req.Method == "GET" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add("Accept", "application/json")

		apiKey := req.FormValue("apiKey")
		eidStr := req.FormValue("eid")
		accountIdStr := req.FormValue("aid")
		switch {
		case apiKey == "":
			d := BalancesResponse{
				Result:  "error",
				Message: "'apiKey' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case eidStr == "":
			d := BalancesResponse{
				Result:  "error",
				Message: "'eid' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case accountIdStr == "":
			d := BalancesResponse{
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
			d := BalancesResponse{
				Result:  "error",
				Message: "Secret not found, 'apiKey' parameter not correct.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		fmt.Printf("apiKey=%v\napiSecret=%v\nexchange_id=%v\naccount_id=%v\n", apiKey, apiSecret, eid, accountId)
		//Use switch statement to switch the exchange_id and make request to the appropriate exchange
		switch eid {
		case 1:
			queryString := fmt.Sprintf("http://localhost:5030/returnNonZeroBalances?apiKey=%s&secret=%s", apiKey, apiSecret)
			fmt.Println("buyOrder query string is ", queryString)
			resp, err := http.Get(queryString)
			if err != nil {
				d := BalancesResponse{
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
				d := BalancesResponse{
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
			d := BalancesResponse{
				Result:  "error",
				Message: "Sorry!, Other exchanges are not yet supported.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
		}
	}
}
