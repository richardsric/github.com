package helper
/*
import (
  "encoding/json"
	"strings"
)
*/
type BalancesResponse struct {
  	Result string            `json:"result"`
  	Message string          `json:"message"`
  	Details []Balance       `json:"details"`
}
type BalanceResponse struct {
	Result string            `json:"result"`
	Message string          `json:"message"`
	Details Balance       `json:"details"`
}
type Address struct {
	Currency string `json:"Currency"`
	Address  string `json:"Address"`
}
type AllBalance struct  {
  Values []Balance
}
type Balance struct  {
  Currency      string  `json:"Currency"`
  Balance       float64 `json:"Balance"`
  Available     float64 `json:"Available"`
  Pending       float64 `json:"Pending"`
  CryptoAddress string  `json:"CryptoAddress"`
}
