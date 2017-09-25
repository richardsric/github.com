package helper

import (
	"encoding/json"
	"fmt"
)

func BalanceResponseHandler(res []byte, err error, reqType string) string {
	var bResponse BittrexJsonResponse
	if err != nil {
		oResp := BalancesResponse{
			Result:  "error",
			Message: "Request to bittrex not successfull due to no connection available", //err.Error(),
		}
		bs, _ := json.Marshal(oResp)
		return string(bs)
	}
	if err = json.Unmarshal(res, &bResponse); err != nil {
		oResp := BalancesResponse{
			Result:  "error",
			Message: err.Error(),
		}
		bs, _ := json.Marshal(oResp)
		return string(bs)
	}
	if err = BittrexErrHandle(bResponse); err != nil {
		oResp := BalancesResponse{
			Result:  "error",
			Message: err.Error(),
		}
		bs, _ := json.Marshal(oResp)
		return string(bs)
	}
	switch reqType {
	case "GetBalances":
		var balances []Balance
		err = json.Unmarshal(bResponse.Result, &balances)
		balanceStr := BalancesResponse{
			Result:  "success",
			Details: balances,
		}
		bs, _ := json.Marshal(balanceStr)
		//		fmt.Println(string(bs))
		return string(bs)
	case "GetNonZeroBalances":
		dat := make([]map[string]interface{}, 0)
		if err := json.Unmarshal(bResponse.Result, &dat); err != nil {
			panic(err)
		}
		bal := make([]Balance, 0)
		for idx := range dat {
			if dat[idx]["Balance"].(float64) > 0 { //Checking if balance is not zero
				//			fmt.Println(dat[idx])
				d := Balance{
					Currency:      dat[idx]["Currency"].(string),
					Balance:       dat[idx]["Balance"].(float64),
					Available:     dat[idx]["Available"].(float64),
					Pending:       dat[idx]["Pending"].(float64),
					CryptoAddress: dat[idx]["CryptoAddress"].(string),
				}
				bal = append(bal, d)
			}
		}
		md := BalancesResponse{
			Result:  "success",
			Details: bal,
		}
		//fmt.Println("Non zero balance details are ",bal)
		bs, _ := json.Marshal(md)
		fmt.Println("Non zero balance details are ", string(bs))
		return string(bs)
	case "GetBalance":
		var balance Balance
		err = json.Unmarshal(bResponse.Result, &balance)
		//		fmt.Println(balance)
		balanceStr := BalanceResponse{
			Result:  "success",
			Details: balance,
		}
		bs, _ := json.Marshal(balanceStr)
		return string(bs)
	default:
		oResp := BalancesResponse{
			Result:  "error",
			Message: "Request not yet handled",
		}
		bs, _ := json.Marshal(oResp)
		return string(bs)
	}
}
