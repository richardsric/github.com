package worker

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/richardsric/workers/worker/helpers"

	"time"
)

// SellWorker is func to update sell worker table
func SellWorker() {
	start := time.Now() // get current time

	fmt.Println("Get Jobs From DB")
	con, err := helpers.OpenConnection()
	if err != nil {
		//return err
		fmt.Println(err)
	}
	defer con.Close()

	row, err := con.Db.Query("SELECT id_sell_worker,market,exchange_id,highest_bid_price,lowest_bid_price,highest_ask_price,lowest_ask_price,highest_volume,lowest_volume,actual_rate,actual_quantity,profit_keep,sell_trigger,last_volume,start_volume,ask_bid,order_type,quantity,profit_lock_start,work_status,account_id,order_date,order_id,threshold,sell_price,volume_diff,vol_percent,pl,percent_profit,profit_lock_start_price,cost,last_proceed,last_bid,exit_price,manual_exit,high_profit,high_profit_perc,actual_locked_profit,actual_locked_perc_profit,locked_proceed,percent_exit_profit,node_id,work_id,profit_locked,akey,stop_loss_active,stop_loss,txn_fee,tpfee,work_age,work_started_on,highest_proceed FROM sell_worker WHERE work_status = 0 ")
	if err != nil {
		fmt.Println("Select Failed Due To: ", err)
	}
	defer row.Close()

	for row.Next() {
		fmt.Println("Entered row dot Next")
		var workerID, market, exchangeID, orderType, aKey, workAge string
		var highBid, lowBid, highAsk, lowAsk, highVol, lowVol, lastVol, startVol, actualRate, actualQty, profitKeep, selTrigger, askBid, Qty float64
		var profitLockStart, thresHold, sellPrice, volumeDiff, volumePercent, PL, percentProfit, profitLockStartPrice, cost, lastProceed, lastBid, exitPrice, highestProceed float64
		var highProfit, highProfitPercent, actualLockProfit, actualLockPerProfit, lockProceed, percentExitProfit, stopLossP, txnFee, tpFee float64
		var workStatus, accountID, orderID, nodeID, workID, manualExit, profitLocked, stopLossActive int
		var orderDate, workStartedOn time.Time
		err = row.Scan(&workerID, &market, &exchangeID, &highBid, &lowBid, &highAsk, &lowAsk, &highVol, &lowVol, &actualRate, &actualQty, &profitKeep,
			&selTrigger, &lastVol, &startVol, &askBid, &orderType, &Qty, &profitLockStart, &workStatus, &accountID, &orderDate, &orderID, &thresHold, &sellPrice,
			&volumeDiff, &volumePercent, &PL, &percentProfit, &profitLockStartPrice, &cost, &lastProceed, &lastBid, &exitPrice, &manualExit, &highProfit, &highProfitPercent,
			&actualLockProfit, &actualLockPerProfit, &lockProceed, &percentExitProfit, &nodeID, &workID, &profitLocked, &aKey, &stopLossActive, &stopLossP, &txnFee, &tpFee,
			&workAge, &workStartedOn, &highestProceed)
		if err != nil {
			fmt.Println("Row Scan Failed Due To: ", err)
		}

		//http://localhost:5000/pair/price?pair=btc-bcc&eid=1
		// call the end point with the gotten values.
		body, err := helpers.GetHTTPRequest("http://localhost:5000/pair/price?pair=" + market + "&eid=" + exchangeID + "")
		fmt.Println(string(body))
		if err != nil {
			fmt.Println("Error On Bittrex GetTicker Func", err)
			return
		}
		// unmarshal the json response.
		var m map[string]interface{}
		err = json.Unmarshal(body, &m)
		if err != nil {
			//panic(err)
			fmt.Println(err)
		}

		//pair := m["market"]
		ask := m["ask"]
		bid := m["bid"]
		//high := m["high"]
		//low := m["low"]
		vol := m["volume"]

		//Check to b sure it retrieves valid data
		if bid.(float64) > 0 && ask.(float64) > 0 {

			/// check Exit price for 1
			if manualExit == 1 {
				fmt.Println("..........Enter Manuel Exit 1...................")
				aID := fmt.Sprintf("%v", accountID)
				aQty := fmt.Sprintf("%v", actualQty)
				Bid := fmt.Sprintf("%v", bid.(float64))
				//key := getKey(aID)

				_, err := SellOrder(exchangeID, aID, Bid, aQty, aKey, market)
				if err != nil {
					fmt.Println("Sell Order Failed ", err)
				} else {
					workStatus = 1
					exitPrice = bid.(float64)
				}
			}

			/// check Exit price for 2

			if manualExit == 2 {
				aID := fmt.Sprintf("%v", accountID)
				eID := fmt.Sprintf("%s", exchangeID)
				aQty := fmt.Sprintf("%v", actualQty)
				exP := fmt.Sprintf("%v", exitPrice)
				// Get api key
				//key := getKey(aID)
				fmt.Println("Quantity =", aQty)
				fmt.Println("Exit Price =", exP)
				_, err := SellOrder(eID, aID, exP, aQty, aKey, market)
				if err != nil {
					fmt.Println("Sell Order Failed ", err)
				} else {
					workStatus = 1
				}

			}

			if lastBid == 0 {
				fmt.Println("Values first time:", lowAsk, lowBid, lowVol)
				highAsk = ask.(float64)
				lowAsk = ask.(float64)
				highBid = bid.(float64)
				lowBid = bid.(float64)
				highVol = vol.(float64)
				lowVol = vol.(float64)
				startVol = vol.(float64)
				lastVol = vol.(float64)
				lastBid = bid.(float64)
				//set work start timestamp here
				workStartedOn = time.Now()
				//workStartedOn = t.Format("20060102150405")
				//fmt.Println("WorkStartedOn: ", workStartedOn)
				//workStartedOn = time.Now()
			}

			/// check if high ask price have changed
			if ask.(float64) > highAsk {
				highAsk = ask.(float64)
			}
			fmt.Println("HighestAsk: ", highAsk)

			/// check if low ask price have changed
			if ask.(float64) < lowAsk {
				lowAsk = ask.(float64)
			}
			fmt.Println("LowestAsk: ", lowAsk)

			/// check if high bid price have changed

			oldHighBid := highBid
			fmt.Println("Old Highest Bid: ", oldHighBid)
			if bid.(float64) > highBid {
				highBid = bid.(float64)
			}
			fmt.Println("New HighestBid: ", highBid)

			/// check if low bid price have changed
			//oldLowBid := lowBid
			if bid.(float64) < lowBid {
				lowBid = bid.(float64)
			}
			fmt.Println("LowestBid: ", lowBid)

			/// check if high vol have changed
			oldHighVol := highVol
			fmt.Println("Old Highest Volume: ", oldHighVol)
			if bid.(float64) > oldHighVol {
				highVol = vol.(float64)
			}
			fmt.Println("HighestVol: ", highVol)

			/// check if low bid price have changed
			if bid.(float64) < lowVol {
				lowVol = vol.(float64)
			}
			fmt.Println("Lowestvol: ", lowVol)

			// Set First time Value insert into the sell worker table.

			// Compute vol_diff
			// volume_diff (last_volume - start_volume)
			volumeDiff = RoundDown(vol.(float64)-startVol, 8)
			fmt.Println("Volume Difference: ", volumeDiff)

			// Compute vol percent
			// vol_percent (volume_diff/start_volume).
			volumePercent = RoundDown((volumeDiff/startVol)*100, 4)
			fmt.Println("Volume Percent: ", volumePercent)

			// profit lock is not disabled
			if profitLockStart != 0 {
				// Compute if we need to set d profit_lock_start_price.
				if profitLockStartPrice == 0 {
					profitLockStartPrice = RoundDown(actualRate+RoundDown(((profitLockStart/100)*actualRate), 8), 8)
				}
				// check if bid is upto price to start d locking.
				if bid.(float64) >= profitLockStartPrice && profitLocked == 0 {
					profitLocked = 1
				}

				if profitLocked == 1 && (thresHold == 0 || (bid.(float64) >= oldHighBid)) {
					//put highbid computations here####
					var highRange, threshval, sellval, highproceed float64
					highRange = highBid - actualRate

					/// compute threshold
					threshval = RoundDown((highRange * (profitKeep + selTrigger) / 100), 8)
					thresHold = RoundDown(actualRate+threshval, 8)
					fmt.Println("threshold: ", thresHold)

					// compute sell_price
					sellval = RoundDown((highRange * (profitKeep + (selTrigger / 2)) / 100), 8)
					sellPrice = RoundDown(actualRate+sellval, 8)
					fmt.Println("SellPrice: ", sellPrice)

					// compute high profit
					highproceed = RoundDown((actualQty * highBid), 8)
					highProfit = RoundDown((highproceed - RoundDown(((txnFee/100)*highproceed), 8) - cost), 8)
					fmt.Println("Highest Profit: ", highProfit)
					highestProceed = highproceed - RoundDown(((txnFee/100)*highproceed), 8)
					fmt.Println("Highest Proceed: ", highproceed)
					// compute high profit percentage
					highProfitPercent = RoundDown(((highBid-actualRate)/actualRate)*100, 4)
					fmt.Println("High Profit Percent: ", highProfitPercent)

					// compute locked_proceed
					lockProceed = RoundDown((actualQty*sellPrice), 8) - RoundDown(((actualQty*sellPrice)*(txnFee/100)), 8)
					fmt.Println("Locked Proceed: ", lockProceed)

					// compute for actual_locked_profit
					actualLockProfit = RoundDown((lockProceed - cost), 8)
					fmt.Println("Actual Locked Profit: ", actualLockProfit)

					// check compute
					actualLockPerProfit = RoundDown(((sellPrice-actualRate)/actualRate)*100, 4)
					fmt.Println("Actual Locked Profit Percent: ", actualLockPerProfit)

				}

			} //profit lock computations done if enabled

			//Calculate normal values

			// Compute last_proceed
			fullproceed := RoundDown(actualQty*bid.(float64), 8)
			lastProceed = RoundDown(fullproceed-RoundDown((txnFee/100)*fullproceed, 8), 8)
			fmt.Println("Last Proceed: ", lastProceed)

			// Compute PL
			PL = RoundDown((lastProceed - cost), 8)
			fmt.Println("PL: ", PL)

			// Compute per percent_profit
			percentProfit = RoundDown(((bid.(float64)-actualRate)/actualRate)*100, 4)
			fmt.Println("PL Percent: ", percentProfit)

			//end normal value computations

			if manualExit == 0 && bid.(float64) <= thresHold {
				//this means dat threshold has bn breached. We issue sell order here.
				aID := fmt.Sprintf("%v", accountID)
				eID := fmt.Sprintf("%s", exchangeID)
				aQty := fmt.Sprintf("%v", actualQty)
				exP := fmt.Sprintf("%v", exitPrice)
				// Get api key
				//key := getKey(aID)

				_, err := SellOrder(eID, aID, exP, aQty, aKey, market)
				if err != nil {
					fmt.Println("Sell Order Failed ", err)
				}

				// compute for exit price
				exitPrice = sellPrice
				workStatus = 1

			}

			if exitPrice > 0 {
				percentExitProfit = RoundDown(((exitPrice-actualRate)/actualRate)*100, 4)
			}

			workID = 1111
			fmt.Println("SELL WORKER ANALYSIS FOR WORKER ID:", workID)

			fmt.Println("...................... ")
			fmt.Println("TICKER INFO")
			fmt.Println("...................... ")
			fmt.Println("Market:", market)
			fmt.Println("Last Bid:", bid)
			fmt.Println("Highest Bid:", highBid)
			fmt.Println("Lowest Bid:", lowBid)
			fmt.Println("Highest Ask:", highAsk)
			fmt.Println("Lowest Ask:", lowAsk)

			fmt.Println("...................... ")
			fmt.Println("MARKET VOLUME INFO")
			fmt.Println("...................... ")
			fmt.Println("Vol @ Initially Work Start:", startVol)
			fmt.Println("Last Vol:", lastVol)
			fmt.Println("Highest Vol:", highVol)
			fmt.Println("Lowest Vol:", lowVol)
			fmt.Println("Volume Difference: ", volumeDiff)
			fmt.Println("Volume Percent: ", volumePercent, "%")

			fmt.Println("......................... ")
			fmt.Println("COSTS/PROFIT/LOSS INFO ")
			fmt.Println("......................... ")
			fmt.Println("Actual Quanity:", actualQty)
			fmt.Println("Actual Rate:", actualRate)
			fmt.Println("Cost:", cost)
			fmt.Println("Last Proceed:", lastProceed)
			fmt.Println("Profit/Loss:", PL)
			fmt.Println("Percent PL:", percentProfit, "%")

			fmt.Println("................................. ")
			fmt.Println("AUTO TRADE PROFIT ADJUSTMENT INFO")
			fmt.Println("................................. ")
			if profitLocked > 0 {
				fmt.Println("Profit Lock/Adjustment: <ON>")
				//show all lock related figure
				fmt.Println("Profit Lock Engaged @:", profitLockStart, "% Profit")
				fmt.Println("Profit Lock Start Price:", profitLockStartPrice)
				fmt.Println("Minimum Profit To Keep:", profitKeep, "% of ", highProfit)
				fmt.Println("Threshold:", thresHold)
				fmt.Println("Auto Sell Price:", sellPrice)
				fmt.Println("Actual Locked Profit:", actualLockProfit, "out of", highProfit)
				fmt.Println("Actual Percent Locked Profit:", actualLockPerProfit, "%")
				fmt.Println("Locked Proceed:", lockProceed)
				fmt.Println("Highest Profit attained:", highProfit)
				fmt.Println("Highest Proceed reached:", highestProceed)
				fmt.Println("Highest Percentage Profit attained:", highProfitPercent, "%")

			} else {
				//hide every lock related figure
				fmt.Println("Profit Lock/Adjustment: <Off>")

			}

			if stopLossActive > 0 || stopLossP > 0 {
				fmt.Println("Stop Loss: <ON>")
				fmt.Println("Stop Loss Set @ :", stopLossP)
			} else {
				fmt.Println("Stop Loss: <OFF>")
			}

			if exitPrice > 0 {
				fmt.Println("Percent Exit Price:", percentExitProfit, "%")
			}
			fmt.Println("Manual Exit:", manualExit, "Current Price: ", bid, " Exit Price: ", exitPrice)
			fmt.Println(".....................................")
			elapsed := time.Since(start)
			fmt.Println(".....................................")
			fmt.Println("TIME TAKEN FOR COMPLETETION")
			fmt.Println("Time Started:", start)
			fmt.Println("Time End:", time.Now())
			fmt.Println("Time Difference/Elapsed:", elapsed)

			_, err = con.Db.Exec("UPDATE sell_worker SET highest_bid_price = $1, lowest_bid_price= $2,highest_ask_price =$3,lowest_ask_price=$4,last_bid=$5,highest_volume=$6,lowest_volume=$7,threshold=$8,sell_price=$9,high_profit=$10,high_profit_perc=$11,actual_locked_profit=$12,actual_locked_perc_profit=$13,locked_proceed=$14,last_proceed=$15,pl=$16,exit_price=$17,last_volume=$18,volume_diff=$19,vol_percent=$20,percent_profit=$21,percent_exit_profit=$22,profit_lock_start_price=$23,start_volume=$24, work_status=$25,stop_loss_active=$26,work_started_on=$27::timestamp,profit_locked=$28,work_age = age(current_timestamp, work_started_on::timestamp), work_count=work_count+1, highest_proceed=$29 WHERE id_sell_worker = $30",
				highBid, lowBid, highAsk, lowAsk, bid, highVol, lowVol, thresHold, sellPrice, highProfit, highProfitPercent, actualLockProfit,
				actualLockPerProfit, lockProceed, lastProceed, PL, exitPrice, vol, volumeDiff, volumePercent, percentProfit, percentExitProfit, profitLockStartPrice, startVol, workStatus, stopLossActive, workStartedOn, profitLocked, highestProceed, workerID)
			if err != nil {
				fmt.Println("Update Failed Due To: ", err)
			}
		} else {
			fmt.Println("Invalid Market Data. ASK/BID Must be greater than 0. Please Check Network connection")

		}
	}

}

// Round var f float64 = 514.89317306
// Round this rounds Output: 514.89317306 to 514
func Round(input float64) float64 {
	if input < 0 {
		return math.Ceil(input - 0.5)
	}
	return math.Floor(input + 0.5)
}

// RoundUp sample output in 4 decmial places
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

// RoundDown Output in 4 decmial places
// var f float64 = 514.89317306
//RoundDown this round Output: 514.89317306 to 514.8931
func RoundDown(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Floor(digit)
	newVal = round / pow
	return
}
