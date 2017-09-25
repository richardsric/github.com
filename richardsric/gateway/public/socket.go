package public

import (
	"encoding/json"
	"fmt"

	"github.com/richardsric/gateway/public/helpers"

	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

// Response this use to construct json response for sockek data
type Response struct {
	Result []MAskBid `json:"result"`
}

// MAskBid This struct is use to push out socket market data
type MAskBid struct {
	Market string  `json:"market"`
	Ask    float64 `json:"ask"`
	Bid    float64 `json:"bid"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Volume float64 `json:"volume"`
}

// MainAskBid this use to construct ask bid for socket data
type MainAskBid struct {
	Values []MAskBid
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Stat is use to track the number excution time of the programe
var Stat = make([]time.Duration, 0)

//var start time.Time

// RequestNo is the number request that pass through the end point within an hr
var RequestNo = 0

// Websocket this is end point for web socket connections ws://localhost:3000/websocket
func Websocket(w http.ResponseWriter, r *http.Request) {
	//var res []byte
	//var start time.Time
	askBidConstruct := make([]MAskBid, 0)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("Client subscribed")

	for {
		//start = time.Now() // get current time
		timeInterval := helpers.GetTimerInterval("Websocket")
		//time.Sleep(2 * time.Second)
		time.Sleep(timeInterval * time.Second)

		body, err := GetTicker("http://localhost:5030/bittrex_ticker")
		//fmt.Fprint(w, string(body))

		if err != nil {
			//return err
			fmt.Println("Error Getting Response From Bittrex. Going to Our DB For The Rquest Data:", err)
			//res := GetAskBidDB(pair, eID)
			//fmt.Fprint(w, string(res))
			return
		}

		var m interface{}
		err = json.Unmarshal(body, &m)
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
					ask := val2.(map[string]interface{})["Ask"]
					bid := val2.(map[string]interface{})["Bid"]
					//last := val2.(map[string]interface{})["Last"]
					high24hr := val2.(map[string]interface{})["High"]
					low24hr := val2.(map[string]interface{})["Low"]
					vol := val2.(map[string]interface{})["Volume"]
					//baseVol := val2.(map[string]interface{})["BaseVolume"]
					//exchangeID := 2

					result := MAskBid{
						Market: MarketName.(string),
						Ask:    ask.(float64),
						Bid:    bid.(float64),
						High:   high24hr.(float64),
						Low:    low24hr.(float64),
						Volume: vol.(float64),
					}
					askBidConstruct = append(askBidConstruct, result)
					//res, _ = json.Marshal(result)
					//fmt.Println(string(res))
					//fmt.Fprint(w, string(res))
					//elapsed := time.Since(start)
					//Something = append(Something, elapsed)
				}
			}
		}
		data := Response{
			Result: askBidConstruct,
		}
		response, _ := json.Marshal(data)
		err = conn.WriteMessage(websocket.TextMessage, response)
		if err != nil {
			fmt.Println(err)
			//	fmt.Println("Client unsubscribed")
			conn.Close()
			break
		}

	}
	//RequestNo = RequestNo + 1
	//elapsed := time.Since(start)
	//Stat = append(Stat, elapsed)

	//fmt.Printf("Timer Logger: %s\n", Stat)
}

// Test this is use to test webscoket conncetion for data
func Test(w http.ResponseWriter, r *http.Request) {

	indexFile, err := os.Open("index.html")
	if err != nil {
		fmt.Println(err)
	}
	index, err := ioutil.ReadAll(indexFile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(index))
}
