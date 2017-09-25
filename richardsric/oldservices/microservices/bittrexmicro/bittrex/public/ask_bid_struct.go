package public

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
