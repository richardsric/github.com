package public

import (
	"encoding/json"
	"fmt"

	"github.com/richardsric/microservices/bittrexmicro/bittrex/public/helpers"

	"net/http"
)

// BittrexSinglePair is the function that will return the ask and bid or error.
func BittrexSinglePair(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Request Received Here On Bittrex")
	pair := r.FormValue("pair")
	eID := r.FormValue("eid")
	if pair == "" || eID == "" {
		result := AskBid{
			Success: `false`,
			Message: "Empty Field Selected",
		}
		res, _ := json.Marshal(result)
		fmt.Println(string(res))
		fmt.Fprint(w, string(res))
		return
	}
	if eID == "1" {

		//body, err := GetTicker("http://localhost:5052/bittrex_ticker")
		if len(BittreMarketTicker) == 0 {
			fmt.Fprint(w, string(BittreMarketTicker))
			//if err != nil {
			//return err
			fmt.Println("Error Getting Response From Bittrex. Going to Our DB For The Rquest Data:")
			res := GetAskBidDB(pair, eID)
			/* if err == nil {
				result := AskBid{
					Success: `false`,
					Message: "DB Service Not Responding",
				}
				res, _ := json.Marshal(result)
				//fmt.Println(string(res))
				fmt.Fprint(w, string(res))
				return
			} */
			fmt.Fprint(w, string(res))
			return
		}

		var m interface{}
		err := json.Unmarshal(BittreMarketTicker, &m)
		if err != nil {
			//panic(err)
			fmt.Println("The error on itradecoin ask and bid:", err)
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
					MarketName := val2.(map[string]interface{})["MarketName"]
					fmt.Println("Got Sucess As False:", MarketName)

					if pair == MarketName.(string) {

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
						return
					}

				}
			}
		}
		result := AskBid{
			Success: `false`,
			Message: "Invalid Pair For The Selected Market",
		}
		res, _ := json.Marshal(result)
		//fmt.Println(string(res))
		fmt.Fprint(w, string(res))
		return
	}

	result := AskBid{
		Success: `false`,
		Message: "Exchange Not Yet Supported",
	}
	res, _ := json.Marshal(result)
	//fmt.Println(string(res))
	fmt.Fprint(w, string(res))
	return

}

// GetAskBidDB is use to get pair ask and bid price when exchange is nit responding.
func GetAskBidDB(pair, eID string) []byte {
	fmt.Println("Entered Our GetAskBidDB func To Get The Requestd Data From DB Since Bittrex Is Not Responding")
	con, err := helpers.OpenConnection()
	if err != nil {
		//return err
		fmt.Println(err)
	}
	defer con.Close()

	row, err := con.Db.Query("SELECT pair,ask,bid,high24hr,low24hr,volume FROM market_data WHERE pair =$1 AND exchange_id  = $2 ORDER BY date_time DESC LIMIT 1", pair, eID)
	//row, err := con.Db.Query("SELECT pair,ask,bid,high24hr,low24hr,volume FROM market_data WHERE pair ='BTC-BCC' AND exchange_id  = 2 ORDER BY date_time DESC LIMIT 1")

	if err != nil {
		fmt.Println("Select Failed Due To: ", err)
	}
	defer row.Close()

	for row.Next() {
		fmt.Println("Entered row dot Next")
		var pairs string
		var ask, bid, high24hr, low24hr, volume float64

		err = row.Scan(&pairs, &ask, &bid, &high24hr, &low24hr, &volume)
		if err != nil {
			// handle this error
			//panic("Row Scan From Staff Deduction Failed Due To: ", err)
			fmt.Println("Row Scan From Staff Deduction Failed Due To: ", err)
		}
		fmt.Println("Gotten DB value: ", pairs, ask, bid, high24hr, low24hr, volume)
		result := AskBid{
			Success: `true`,
			Message: "",
			Market:  pairs,
			Ask:     ask,
			Bid:     bid,
			High:    high24hr,
			Low:     low24hr,
			Volume:  volume,
		}
		res, _ := json.Marshal(result)
		fmt.Println(string(res))
		//fmt.Fprint(w, string(res))
		return res
	}

	result := AskBid{
		Success: `false`,
		Message: "No Data Return Currently. Check Back Later",
	}
	res, _ := json.Marshal(result)
	fmt.Println(string(res))
	//fmt.Fprint(w, string(res))
	return res

}

// GetAskBid is the function that will return the ask and bid or error.
// func GetAskBid(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	pair := r.FormValue("pair")
// 	eID := r.FormValue("eid")
// 	if pair == "" || eID == "" {
// 		result := AskBid{
// 			Success: `false`,
// 			Message: "Empty Field Selected",
// 		}
// 		res, _ := json.Marshal(result)
// 		fmt.Println(string(res))
// 		fmt.Fprint(w, string(res))
// 		return
// 	}
// 	if eID == "2" {

// 		body, err := bittrex.BittrexMarketData()
// 		//fmt.Fprint(w, string(body))
// 		if err != nil {
// 			//return err
// 			fmt.Println("Error Getting Response From Bittrex. Going to Our DB For The Rquest Data:", err)
// 			res := GetAskBidDB(pair, eID)
// 			/* if err == nil {
// 				result := AskBid{
// 					Success: `false`,
// 					Message: "DB Service Not Responding",
// 				}
// 				res, _ := json.Marshal(result)
// 				//fmt.Println(string(res))
// 				fmt.Fprint(w, string(res))
// 				return
// 			} */
// 			fmt.Fprint(w, string(res))
// 			return
// 		}

// 		var m interface{}
// 		err = json.Unmarshal(body, &m)
// 		if err != nil {
// 			//panic(err)
// 			fmt.Println("The error on itradecoin ask and bid:", err)
// 		}
// 		t := m.(map[string]interface{})
// 		for key, val := range t {
// 			//fmt.Println("Got Key1 As:", key, "||", "Got Values1 As:", val)
// 			if key == "success" && val == false {
// 				//panic(err)
// 				fmt.Println("Got Sucess As False:", val)
// 			}
// 			if key == "result" {
// 				for _, val2 := range val.([]interface{}) {
// 					//fmt.Println("Got Key2 As:", key2, "||", "Got Values2 As:", val2)
// 					MarketName := val2.(map[string]interface{})["MarketName"]
// 					fmt.Println("Got Sucess As False:", MarketName)

// 					if pair == MarketName.(string) {

// 						ask := val2.(map[string]interface{})["Ask"]
// 						bid := val2.(map[string]interface{})["Bid"]
// 						//last := val2.(map[string]interface{})["Last"]
// 						high24hr := val2.(map[string]interface{})["High"]
// 						low24hr := val2.(map[string]interface{})["Low"]
// 						vol := val2.(map[string]interface{})["Volume"]
// 						//baseVol := val2.(map[string]interface{})["BaseVolume"]
// 						//exchangeID := 2

// 						result := AskBid{
// 							Success: `true`,
// 							Message: "",
// 							Market:  MarketName.(string),
// 							Ask:     ask.(float64),
// 							Bid:     bid.(float64),
// 							High:    high24hr.(float64),
// 							Low:     low24hr.(float64),
// 							Volume:  vol.(float64),
// 						}
// 						res, _ := json.Marshal(result)
// 						//fmt.Println(string(res))
// 						fmt.Fprint(w, string(res))
// 						return
// 					}

// 				}
// 			}
// 		}
// 		result := AskBid{
// 			Success: `false`,
// 			Message: "Invalid Pair For The Selected Market",
// 		}
// 		res, _ := json.Marshal(result)
// 		//fmt.Println(string(res))
// 		fmt.Fprint(w, string(res))
// 		return
// 	}

// 	result := AskBid{
// 		Success: `false`,
// 		Message: "Exchange Not Yet Supported",
// 	}
// 	res, _ := json.Marshal(result)
// 	//fmt.Println(string(res))
// 	fmt.Fprint(w, string(res))
// 	return

// }

// // GetAskBidDB is use to get pair ask and bid price when exchange is nit responding.
// func GetAskBidDB(pair, eID string) []byte {
// 	fmt.Println("Entered Our GetAskBidDB func To Get The Requestd Data From DB Since Bittrex Is Not Responding")
// 	con, err := microservices.OpenConnection()
// 	if err != nil {
// 		//return err
// 		fmt.Println(err)
// 	}
// 	defer con.Close()

// 	row, err := con.Db.Query("SELECT pair,ask,bid,high24hr,low24hr,volume FROM market_data WHERE pair =$1 AND exchange_id  = $2 ORDER BY date_time DESC LIMIT 1", pair, eID)
// 	//row, err := con.Db.Query("SELECT pair,ask,bid,high24hr,low24hr,volume FROM market_data WHERE pair ='BTC-BCC' AND exchange_id  = 2 ORDER BY date_time DESC LIMIT 1")

// 	if err != nil {
// 		fmt.Println("Select Failed Due To: ", err)
// 	}
// 	defer row.Close()

// 	for row.Next() {
// 		fmt.Println("Entered row dot Next")
// 		var pairs string
// 		var ask, bid, high24hr, low24hr, volume float64

// 		err = row.Scan(&pairs, &ask, &bid, &high24hr, &low24hr, &volume)
// 		if err != nil {
// 			// handle this error
// 			//panic("Row Scan From Staff Deduction Failed Due To: ", err)
// 			fmt.Println("Row Scan From Staff Deduction Failed Due To: ", err)
// 		}
// 		fmt.Println("Gotten DB value: ", pairs, ask, bid, high24hr, low24hr, volume)
// 		result := AskBid{
// 			Success: `true`,
// 			Message: "",
// 			Market:  pairs,
// 			Ask:     ask,
// 			Bid:     bid,
// 			High:    high24hr,
// 			Low:     low24hr,
// 			Volume:  volume,
// 		}
// 		res, _ := json.Marshal(result)
// 		fmt.Println(string(res))
// 		//fmt.Fprint(w, string(res))
// 		return res
// 	}

// 	result := AskBid{
// 		Success: `false`,
// 		Message: "No Data Return Currently. Check Back Later",
// 	}
// 	res, _ := json.Marshal(result)
// 	fmt.Println(string(res))
// 	//fmt.Fprint(w, string(res))
// 	return res

// }
