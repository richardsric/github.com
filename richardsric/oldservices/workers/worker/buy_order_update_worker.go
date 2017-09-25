package worker

import (
	"encoding/json"
	"fmt"

	"github.com/richardsric/workers/worker/helpers"
)

// BuyOrderUpdateWorker is func to update buy order
func BuyOrderUpdateWorker() {

	fmt.Println("Entered Our buyOrderUpdateWorker func To Get The Requestd Data From DB")
	con, err := helpers.OpenConnection()
	if err != nil {
		//return err
		fmt.Println(err)
	}
	defer con.Close()

	row, err := con.Db.Query("SELECT buy_orders.order_number,buy_orders.account_id,buy_orders.exchange_id,apks.key,apks.secret FROM buy_orders INNER JOIN apks ON apks.account_id = buy_orders.account_id WHERE buy_orders.work_status = 0")
	if err != nil {
		fmt.Println("Select Failed Due To: ", err)
	}
	defer row.Close()

	for row.Next() {
		fmt.Println("Entered row dot Next")
		var orderID, accountID, exchangeID, apiKey, secret string
		err = row.Scan(&orderID, &accountID, &exchangeID, &apiKey, &secret)
		if err != nil {
			fmt.Println("Row Scan Failed Due To: ", err)
		}
		//localhost:5000/getOrderInfo?apiKey=110982d6fd72480d9968cbca3473a868&secret=c14d8e9f65ac44d48ea484320c07230c&uuid=34a42ddc-22b5-493d-a42b-4ddf88ef9ed8&eid=1&aid=1
		// call the end point with the gotten values.
		body, err := helpers.GetHTTPRequest("http://localhost:5000/getOrderInfo?apiKey=" + apiKey + "&secret=" + secret + "&uuid=" + orderID + "&eid=" + exchangeID + "&aid=" + accountID + "")
		fmt.Println(string(body))
		if err != nil {
			fmt.Println("Error On Bittrex GetTicker Func", err)
			return
		}
		// unmarshal the json response.
		var m interface{}
		err = json.Unmarshal(body, &m)
		if err != nil {
			//panic(err)
			fmt.Println(err)
		}
		t := m.(map[string]interface{})
		for key, val := range t {
			fmt.Println("Got Key1 As:", key, "||", "Got Values1 As:", val)

			if key == "result" && val == "error" {
				//panic(err)
				fmt.Println("Got Sucess As False:", val)
			}

			if key == "order_number" {
				//OrderNumber := val
			}

			if key == "order_details" {
				fmt.Println("Enter Order details")
				fmt.Println(val.(map[string]interface{})["market"])
				actualQty := val.(map[string]interface{})["actual_quantity"]
				actualRate := val.(map[string]interface{})["actual_rate"]
				orderStatus := val.(map[string]interface{})["order_status"]
				fee := val.(map[string]interface{})["fee"]
				//orderDate := val2.(map[string]interface{})["order_date"]

				_, err := con.Db.Exec("UPDATE buy_orders SET actual_rate = $1, actual_quantity= $2,order_status =$3,txn_fee=$4 WHERE order_number = $5", actualRate, actualQty, orderStatus, fee.(float64), orderID)
				if err != nil {
					fmt.Println("Execute Insert Failed Due To: ", err)
				}

			}
		}
	}
}
