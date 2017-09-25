package public

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/richardsric/microservices/bittrexmicro/bittrex/public/helpers"
)

// Stat is use to track the number excution time of the programe
var Stat = make([]time.Duration, 0)

// RequestNo is the number request that pass through the end point within an hr
var RequestNo = 0
var numberOfFailed = 0

// BittreMarketTicker this hold the bittrex market data and can be accessed to get the data
var BittreMarketTicker []byte

// MarketDataService this is timer for bittrex market data.
func MarketDataService() {
	//	fmt.Println("Bittrex Market Data Service Started.... The Service Starts In The Next 5 Sec")
	for {

		timeInterval := helpers.GetTimerInterval("MarketDataService")

		time.Sleep(timeInterval * time.Second)
		//time.Sleep(2 * time.Second)

		BittrexMarketDataService()

		//	fmt.Println("Bittrex Market Service Run Successfully.... The Service Will Run Again In The Next 5 Sec")

	}
}

// BittrexMarketData1 is use to return market data for Bittrex Exchange and will be used as bittrex.BittrexMarketData by other packages;.
func BittrexMarketData1(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Entered BittrexMarketData Func")

	//resp, err := BittrexMarketDataService()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	fmt.Fprint(w, string(BittreMarketTicker))

}

// BittrexMarketDataService is use to return market data for Bittrex Exchange and will be used as bittrex.BittrexMarketData by other packages;.
func BittrexMarketDataService() {
	start := time.Now() // get current time
	con, err := helpers.OpenConnection()
	if err != nil {
		//return err
		fmt.Println(err)
	}
	defer con.Close()
	url := "https://bittrex.com/api/v1.1/public/getmarketsummaries"
	//BittreMarketTicker, err = GetTicker("https://bittrex.com/api/v1.1/public/getmarketsummaries")

	BittreMarketTicker, err = GetTicker(url)
	//fmt.Println(string(body))
	if err != nil {
		fmt.Println("Error On Bittrex GetTicker Func", err)
		return
	}

	if len(BittreMarketTicker) == 0 {
		fmt.Println("Nil Response Gotten From The Request", url)
		fmt.Println("Kindly Check Your Internet Connection")
		return
	}
	var m interface{}
	err = json.Unmarshal(BittreMarketTicker, &m)
	if err != nil {
		//panic(err)
		fmt.Println(err)
	}
	if m != nil {
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
					pair := val2.(map[string]interface{})["MarketName"]
					ask := val2.(map[string]interface{})["Ask"]
					bid := val2.(map[string]interface{})["Bid"]
					last := val2.(map[string]interface{})["Last"]
					high24hr := val2.(map[string]interface{})["High"]
					low24hr := val2.(map[string]interface{})["Low"]
					vol := val2.(map[string]interface{})["Volume"]
					baseVol := val2.(map[string]interface{})["BaseVolume"]
					exchangeID := 1

					_, err := con.Db.Exec("INSERT INTO market_data (pair,ask,bid,last,high24hr,low24hr,volume,base_volume,exchange_id)VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)", pair, ask, bid, last, high24hr, low24hr, vol, baseVol, exchangeID)
					if err != nil {
						fmt.Println("Execute Insert Failed Due To: ", err)
					}

				}
			}
		}
	} else {
		numberOfFailed = numberOfFailed + 1

	}

	RequestNo = RequestNo + 1
	elapsed := time.Since(start)
	Stat = append(Stat, elapsed)

	//return body, nil

}

// Statz use t0 show statz
func Statz(w http.ResponseWriter, r *http.Request) {
	fmt.Println("...........................................Entered Stat Function............................................................................")
	var n, smallest, biggest time.Duration
	x := Stat

	for _, v := range x {
		if v > n {
			fmt.Println(v, ">", n)
			n = v
			biggest = n
		} else {
			fmt.Println(v, "<", n)
		}
	}

	fmt.Println("The biggest number is ", biggest)
	for _, v := range x {
		if v > n {
			fmt.Println(v, ">", n)
		} else {
			fmt.Println(v, "<", n)
			n = v
			smallest = n
		}
	}
	fmt.Println("The smallest number is ", smallest)

	fmt.Fprint(w, "iTradeCoin Bittrex Market Update Service Running", "\n\n")
	fmt.Fprint(w, "Number Of Request Recevied Within 1Hr: ", RequestNo, "\n")
	fmt.Fprint(w, "Minimum Execution Time Within 1Hr: ", smallest, "\n")
	fmt.Fprint(w, "Maximum Execution Time Within 1Hr: ", biggest, "\n")
	fmt.Fprint(w, "Number Of Failed Requests Within 1Hr: ", numberOfFailed, "\n")
}

// ClearStatData this clear stat after an 1hr
func ClearStatData() {
	//	fmt.Println("Enter clearStatData.... The Service Starts In The Next 1hr 0 Sec")
	for {

		time.Sleep(3600 * time.Second)
		RequestNo = 0
		Stat = nil

		//		fmt.Println("clearStatData Run Successfully.... The Service Will Run Again In The Next 1hr Sec")

	}
}
