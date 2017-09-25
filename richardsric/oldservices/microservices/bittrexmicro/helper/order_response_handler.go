package helper

import (
	"errors"
	"encoding/json"
	"fmt"
)

// BittrexErrHandle gets JSON response from Bittrex API and deal with error
func BittrexErrHandle(r BittrexJsonResponse) error {
	if !r.Success {
		return errors.New(r.Message)
	}
	return nil
}

func HandleResponse(res []byte, err error, reqType string) string {
	var bResponse BittrexJsonResponse
	if err != nil {
		oResp := OrderResponse {
			Result: "error",
			Message: "Request to bittrex not successfull due to no connection available",//err.Error(),
		}
		bs, _ := json.Marshal(oResp)
		return string(bs)
	}
	if err = json.Unmarshal(res, &bResponse); err != nil {
		oResp := OrderResponse {
			Result: "error",
			Message: err.Error(),
		}
		bs, _ := json.Marshal(oResp)
		return string(bs)
	}
	if err = BittrexErrHandle(bResponse); err != nil {
		oResp := OrderResponse {
			Result: "error",
			Message: err.Error(),
		}
		bs, _ := json.Marshal(oResp)
		return string(bs)
	}
	
	switch reqType {
		case "SellLimit":
			var u BittrexOrderUuid
			err = json.Unmarshal(bResponse.Result, &u)
			oResp := OrderResponse {
				Result: "success",
				OrderNumber: u.OrderNumber,
			}
			fmt.Println("order number for sell limit is ",u.OrderNumber)
			bs, _ := json.Marshal(oResp)
			return string(bs)
		case "BuyLimit":
			var u BittrexOrderUuid
			err = json.Unmarshal(bResponse.Result, &u)
			oResp := OrderResponse {
				Result: "success",
				OrderNumber: u.OrderNumber,
			}
			fmt.Println("order number for buy limit is ",u.OrderNumber)
			bs, _ := json.Marshal(oResp)
			return string(bs)
		case "CancelOrder":
			oResp := OrderResponse  {
				Result:"success",
			  }
			bs, _ := json.Marshal(oResp)
			return string(bs)
		case "GetOrderInfo":
			var bOrder BittrexOrderInfo
			err = json.Unmarshal(bResponse.Result, &bOrder)
			var oType string
			var oStatus string
			var oDate string
			fmt.Printf("Bittrex response for getOrder is %v\n",bOrder)
			if bOrder.Type == "LIMIT_BUY" {
				oType = "BUY"
			}else{
				oType = "SELL"
			}	
			if bOrder.QuantityRemaining == 0 && bOrder.IsOpen == false {
				oStatus = "COMPLETED"
				oDate = bOrder.Opened
			}else if bOrder.IsOpen == false && bOrder.Quantity == bOrder.QuantityRemaining && bOrder.Quantity > 0 && 
			bOrder.QuantityRemaining > 0 && bOrder.CancelInitiated == true && bOrder.Price == 0 {
				oStatus = "CANCELED"
				oDate = bOrder.Opened
			}else if bOrder.IsOpen == true {
				oStatus = "OPEN"
				oDate = bOrder.Opened
			}
			cStr := CustomOrderInfo {
				Market: bOrder.Exchange,
				OrderType: oType,
				ActualQuantity: bOrder.Quantity,
				ActualRate: bOrder.PricePerUnit,
				OrderStatus: oStatus,
				Fee: bOrder.CommissionPaid,
				OrderDate: oDate,
			}
			oInfoResponse := OrderInfoResponse {
				Result: "success",
				OrderNumber: bOrder.OrderUuid,
				OrderDetails: cStr,
			}
			bs, _ := json.Marshal(oInfoResponse)
			return string(bs)
		case "GetBalances":
			//
			break
		case "GetBalance":
			//
			break
		default:
			oResp := OrderResponse {
				Result: "error",
				Message: "Request not yet handled",
			}
			bs, _ := json.Marshal(oResp)
			return string(bs)
	}
	oResp := OrderResponse {
		Result: "error",
		Message: "Request not known",
	}
	bs, _ := json.Marshal(oResp)
	return string(bs)
} 