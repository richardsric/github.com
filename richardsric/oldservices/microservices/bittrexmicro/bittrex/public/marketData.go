package public

import (
	"encoding/json"
	"fmt"
	"time"
)

func MarketDataService1() {
	fmt.Println("Bittrex Market Data Service Started.... The Service Starts In The Next 5 Sec")
	for {

		time.Sleep(2 * time.Second)

		MarketData()

		fmt.Println("Bittrex Market Service Run Successfully.... The Service Will Run Again In The Next 5 Sec")

	}
}

// MarketData is use to return market data for Bittrex Exchange and will be used as bittrex.BittrexMarketData by other packages;.
func MarketData() (result []byte, err error) {
	fmt.Println("Entered MarketData Func")

	body, err := GetTicker("https://bittrex.com/api/v1.1/public/getmarketsummaries")
	if err != nil {
		//panic(err)
		fmt.Println("ask and bid:", err)
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
			fmt.Println("Got Sucess As False:", val)
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

				res := AskBid{
					Success: `true`,
					Message: "",
					Market:  MarketName.(string),
					Ask:     ask.(float64),
					Bid:     bid.(float64),
					High:    high24hr.(float64),
					Low:     low24hr.(float64),
					Volume:  vol.(float64),
				}
				result, _ = json.Marshal(res)
				//fmt.Println(string(result))

			}
		}
	}
	return result, nil
}
