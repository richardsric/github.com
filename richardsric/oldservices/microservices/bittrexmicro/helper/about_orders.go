package helper

import (
	"encoding/json"
)

type BittrexOrderInfo struct {
	AccountId                  string
	OrderUuid                  string `json:"OrderUuid"`
	Exchange                   string `json:"Exchange"`
	Type                       string
	Quantity                   float64 `json:"Quantity"`
	QuantityRemaining          float64 `json:"QuantityRemaining"`
	Limit                      float64 `json:"Limit"`
	Reserved                   float64
	ReserveRemaining           float64
	CommissionReserved         float64
	CommissionReserveRemaining float64
	CommissionPaid             float64
	Price                      float64 `json:"Price"`
	PricePerUnit               float64 `json:"PricePerUnit"`
	Opened                     string
	Closed                     string
	IsOpen                     bool
	Sentinel                   string
	CancelInitiated            bool
	ImmediateOrCancel          bool
	IsConditional              bool
	Condition                  string
	ConditionTarget            string
}
type CustomOrderInfo struct {
	Market string 			`json:"market"`
	OrderType string 		`json:"order_type"`
	ActualQuantity float64 	`json:"actual_quantity"`
	ActualRate float64		`json:"actual_rate"`
	OrderStatus string 		`json:"order_status"`
	Fee float64				`json:"fee"`
	OrderDate string 		`json:"order_date"`
}
type BittrexJsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

type BittrexOrderUuid struct {
	OrderNumber string `json:"uuid"`
}

type OrderResponse struct {
	Result string            `json:"result"`
	Message string          `json:"message"`
	OrderNumber string      `json:"order_number"`
}
type OrderInfoResponse struct {
	Result string            `json:"result"`
	Message string          `json:"message"`
	OrderNumber string      `json:"order_number"`
	OrderDetails CustomOrderInfo	`json:"order_details"`
}