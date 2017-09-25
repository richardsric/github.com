package public

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AskBidPair is use to return market data for Bittrex Exchange and will be used as bittrex.BittrexMarketData by other packages;.
func AskBidPair(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//	fmt.Println("Entered AskBidPair Func")

	pair := r.FormValue("pair")

	url := "https://bittrex.com/api/v1.1/public/getmarketsummary?market=" + pair + ""

	body, err := GetTicker(url)

	if err != nil {
		//panic(err)
		fmt.Println("ask and bid:", err)
	}

	if len(body) == 0 {
		fmt.Println("Nil Response Gotten From The Request", url)
		fmt.Println("Kindly Check Your Internet Connection")
		return
	}

	var m interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		//panic(err)
		fmt.Println("The error Unmarshal Json:", err)
	}
	t := m.(map[string]interface{})

	for key, val := range t {
		//fmt.Println("Got Key1 As:", key, "||", "Got Values1 As:", val)
		if key == "success" && val == false {
			//panic(err)
			//		fmt.Println("Got Sucess As False:", val)
		}
		if key == "result" {
			for _, val2 := range val.([]interface{}) {
				//fmt.Println("Got Key2 As:", key2, "||", "Got Values2 As:", val2)
				//fmt.Println("Got Sucess As False:", MarketName)

				MarketName := val2.(map[string]interface{})["MarketName"]
				ask := val2.(map[string]interface{})["Ask"]
				bid := val2.(map[string]interface{})["Bid"]
				//last := val2.(map[string]interface{})["Last"]
				high24hr := val2.(map[string]interface{})["High"]
				low24hr := val2.(map[string]interface{})["Low"]
				//vol := val2.(map[string]interface{})["Volume"]
				vol := val2.(map[string]interface{})["BaseVolume"]
				//exchangeID := 2

				result := AskBid{
					Success: `true`,
					Message: "",
					Market:  MarketName.(string),
					Ask:     ask.(float64),
					Bid:     bid.(float64),
					High:    high24hr.(float64),
					Low:     low24hr.(float64),
					Volume:  vol.(float64),
				}
				res, _ := json.Marshal(result)
				//fmt.Println(string(res))
				fmt.Fprint(w, string(res))

			}
		}
	}
}
