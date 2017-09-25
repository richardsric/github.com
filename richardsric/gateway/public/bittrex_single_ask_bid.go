package public

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/richardsric/gateway/public/helpers"
)

// AskBid is a struct use to return ask and bid of request pair.
type AskBid struct {
	Success string  `json:"success"`
	Message string  `json:"message"`
	Market  string  `json:"market"`
	Ask     float64 `json:"ask"`
	Bid     float64 `json:"bid"`
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Volume  float64 `json:"volume"`
}

// MainAskBid1 this is use to get single request
type MainAskBid1 struct {
	Values []AskBid
}

// AskBidPair is the function that will return the ask and bid or error.
func AskBidPair(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pair := r.FormValue("pair")
	eID := r.FormValue("eid")
	if pair == "" || eID == "" {
		result := AskBid{
			Success: `false`,
			Message: "Empty Field Selected",
		}
		res, _ := json.Marshal(result)
		//fmt.Println(string(res))
		fmt.Fprint(w, string(res))
		return
	}
	if eID == "1" {

		body, err := GetTicker("http://localhost:5030/pair/price?pair=" + pair + "")
		//fmt.Fprint(w, string(body))
		if err != nil {
			//return err
			fmt.Println("Error Getting Response From Bittrex. Going to Our DB For The Rquest Data:", err)
			res := GetAskBidDB(pair, eID)
			fmt.Fprint(w, string(res))
			return
		}

		//res, _ := json.Marshal(body)
		//fmt.Println(string(res))
		fmt.Fprint(w, string(body))
		return
	}

	//// Exchange not supported

	result := AskBid{
		Success: `false`,
		Message: "Exchange Not Yet Supported",
	}
	res, _ := json.Marshal(result)
	//fmt.Println(string(res))
	fmt.Fprint(w, string(res))
	return

}

// AskBidDB is use to get pair ask and bid price when exchange is nit responding.
func AskBidDB(pair, eID string) []byte {
	//	fmt.Println("Entered Our GetAskBidDB func To Get The Requestd Data From DB Since Bittrex Is Not Responding")
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
		//	fmt.Println("Entered row dot Next")
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
