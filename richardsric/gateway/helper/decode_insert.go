package helper

import (
	"encoding/json"
	"errors"
	"fmt"
)

type OrderResponse struct {
	Result      string `json:"result"`
	OrderNumber string `json:"order_number"`
}

//IsRequestValid is used to check if buy/sell order response gotten from an exchange is successful or not
//if successfull then the order number return will be inserted into the db
func IsRequestValid(response []byte) (status bool, orderNumber string) {
	var dResp OrderResponse
	if err := json.Unmarshal(response, &dResp); err != nil {

		fmt.Println("Error parsing JSON: due to", err)
		return false, ""
	}
	if dResp.Result == "success" {
		return true, dResp.OrderNumber
	}
	return false, ""
}

//DoInsert performs the insertion of the order number if successful to the db
func DoInsert(acct_id, eid int, quantity, limit, txnFee float64, market, tradeType, orderId string) error {
	//	fmt.Println("Starting DoInsert Operation")
	con, err := OpenConnection()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer con.Close()
	//buy_order_table := `buy_order`
	//sell_order_table := `sell_order`
	switch eid {
	case 1:
		fmt.Println("exchange is bittrex")
		fmt.Println("trade type is ", tradeType)
		fmt.Printf("Now inserting info of '%v' as order_number into %s order table \n", orderId, tradeType)
		if tradeType == "BUY" {
			_, err := con.Db.Exec("INSERT INTO buy_orders(market,ask_bid,order_type,order_number,account_id,exchange_id,txn_fee,quantity) VALUES($1,$2,$3,$4,$5,$6,$7,$8)",
				market, limit, tradeType, orderId, acct_id, eid, txnFee, quantity)
			if err != nil {
				fmt.Println("ERROR!!...could not insert into buy_order table due to", err)
				return err
			}
		} else if tradeType == "SELL" {
			_, err := con.Db.Exec("INSERT INTO sell_orders(market,ask_bid,order_type,order_number,account_id,exchange_id,txn_fee,quantity) VALUES($1,$2,$3,$4,$5,$6,$7,$8)",
				market, limit, tradeType, orderId, acct_id, eid, txnFee, quantity)
			if err != nil {
				fmt.Println("ERROR!!...could not insert into sell_order table due to", err)
				return err
			}
		}
		break
	default:
		return errors.New("other exchanges not yet implemented.")
	}
	return nil
}

//GetSecret is used to get the Api Secret stored in the database
//it collects exchange_id,apiKey and account_id as parameter and returns the secret as astring and error
func GetSecret(eid, aid int, apiKey string) (secret string, err error) {
	//fmt.Println("Starting GetSecret Operation")
	con, err := OpenConnection()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer con.Close()
	var apiSecret string
	err = con.Db.QueryRow("SELECT secret FROM apks WHERE account_id = $1 AND exchange_id = $2 and key = $3",
		aid, eid, apiKey).Scan(&apiSecret)
	if err != nil {
		fmt.Println("could not get secret from db due to ", err)
		return "", err
	}
	fmt.Println("successfully selected secret as ", apiSecret)
	return apiSecret, nil
}

//GetTxnFee is used to get the transaction fee of the exchange set in the database,
//it collects the exchange_id as parameter and returns the fee and error
func GetTxnFee(eid int) (txnFee float64, err error) {
	//fmt.Println("Starting GetTxnFee Operation")
	con, err := OpenConnection()
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	defer con.Close()
	err = con.Db.QueryRow("SELECT txn_fee FROM exchange WHERE exchange_id = $1", eid).Scan(&txnFee)
	if err != nil {
		fmt.Println("could not get txnFee from db due to ", err)
		return -1, err
	}
	fmt.Println("successfully selected txnFee as ", txnFee)
	return txnFee, nil
}
