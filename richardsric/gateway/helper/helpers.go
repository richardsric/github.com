package helper

import "math"

// GetProfitLoss This is use to get profit/Lost
func GetProfitLoss(currentPairWorth, initialPairWorth float64) (float64, float64) {
	pwpl := currentPairWorth - initialPairWorth
	if pwpl > 0 {
		pwpl = RoundDown(pwpl, 8)
		changePwpl := pwpl / initialPairWorth
		percentPwpl := RoundDown(changePwpl*100, 4)
		return pwpl, percentPwpl
	} else {
		pwpl = RoundUp(pwpl, 8)
		changePwpl := pwpl / initialPairWorth
		percentPwpl := RoundUp(changePwpl*100, 4)
		return pwpl, percentPwpl

	}

}

// GetAssetWorthPrimary This is use to get asset worth of a primary currency in a trade
func GetAssetWorthPrimary(pairWorth, BidPrice float64) float64 {
	aWorth := pairWorth * BidPrice
	// Check if the result is negative or postive and round to 8 decmial places
	if aWorth > 0 {
		aWorth = RoundDown(aWorth, 8)
		exchangeFee := (0.25 / 100) * aWorth
		aWorth = aWorth - exchangeFee
		return RoundDown(aWorth, 8)
	} else {
		aWorth = RoundUp(aWorth, 8)
		exchangeFee := (0.25 / 100) * aWorth
		aWorth = aWorth - exchangeFee
		return RoundUp(aWorth, 8)
	}
}

// GetAssetWorthSecondary This is use to get asset worth of a secondary currency in a trade
func GetAssetWorthSecondary(pairWorth, AskPrice float64) float64 {
	aWorth := pairWorth / AskPrice
	// Check if the result is negative or postive and round to 8 decmial places
	if aWorth > 0 {
		aWorth = RoundDown(aWorth, 8)
		exchangeFee := (0.25 / 100) * aWorth
		aWorth = aWorth - exchangeFee
		return RoundDown(aWorth, 8)
	} else {
		aWorth = RoundUp(aWorth, 8)
		exchangeFee := (0.25 / 100) * aWorth
		aWorth = aWorth - exchangeFee
		return RoundUp(aWorth, 8)
	}
}

// GetPairWorthPrimary This is use to get pair worth of a primary currency in a trade
func GetPairWorthPrimary(startCapital, askPrice float64) float64 {
	pWorth := startCapital / askPrice
	// Check if the result is negative or postive and round to 8 decmial places
	if pWorth > 0 {
		pWorth = RoundDown(pWorth, 8)
		exchangeFee := (0.25 / 100) * pWorth
		pWorth = pWorth - exchangeFee
		return RoundDown(pWorth, 8)
	} else {
		pWorth = RoundUp(pWorth, 8)
		exchangeFee := (0.25 / 100) * pWorth
		pWorth = pWorth - exchangeFee
		return RoundUp(pWorth, 8)
	}
}

// GetPairWorthSecondary This is use to get pair worth of a secondary currency in a trade
func GetPairWorthSecondary(startCapital, BidPrice float64) float64 {
	pWorth := startCapital * BidPrice
	// Check if the result is negative or postive and round to 8 decmial places
	if pWorth > 0 {
		pWorth = RoundDown(pWorth, 8)
		exchangeFee := (0.25 / 100) * pWorth
		pWorth = pWorth - exchangeFee
		return RoundDown(pWorth, 8)
	} else {
		pWorth = RoundUp(pWorth, 8)
		exchangeFee := (0.25 / 100) * pWorth
		pWorth = pWorth - exchangeFee
		return RoundUp(pWorth, 8)
	}
}

//var f float64 = -514.89317306
//fmt.Printf("%0.2f \n", f)
//fmt.Println(Round(f)) // round half
//fmt.Println(RoundUp(f, 2)) // up to precision 4
//fmt.Println(RoundDown(f, 2)) // up to precision 4
// sample Output
// var f float64 = 514.89317306
// Round this rounds Output: 514.89317306 to 514
func Round(input float64) float64 {
	if input < 0 {
		return math.Ceil(input - 0.5)
	}
	return math.Floor(input + 0.5)
}

// sample Output in 4 decmial places
// var f float64 = 514.89317306
//RoundUp this rounds Output: 514.89317306 to 514.8932
func RoundUp(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Ceil(digit)
	newVal = round / pow
	return
}

// sample Output in 4 decmial places
// var f float64 = 514.89317306
//RoundUp this round Output: 514.89317306 to 514.8931
func RoundDown(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Floor(digit)
	newVal = round / pow
	return
}
