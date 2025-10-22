package finder

import (
	"fmt"

	"mad-scanner.com/scriner/common"
)

// CalcMaxVolume находит максимальный и минимальный уровень стакана,
// на котором сохраняется положительный профит.
// withdrawFee - комиссия на вывод в монетах!
// return: askTotalVol, bidTotalVol
func CalcMinMaxVolume(asks, bids *[]common.Order, thenMakerFee, withdrawFee float64) (common.LowerLimit, common.UpperLimit) {
	asksVolume := BookSum(asks)

	// results
	var upper common.UpperLimit
	var lower common.LowerLimit
	//
	var profit, profitWithFee float64
	var buyPrice, sellPrice float64
	var firstProfitFunded bool = false

	for srcCoins := 1.0; srcCoins <= asksVolume; srcCoins++ {
		// coins
		dstCoins := (srcCoins - withdrawFee) - (srcCoins-withdrawFee)*thenMakerFee
		//usdt
		buyPrice = CalculatePrice(asks, srcCoins)
		sellPrice = CalculatePrice(bids, dstCoins)
		// usdt
		profit = sellPrice*srcCoins - buyPrice*srcCoins
		profitWithFee = sellPrice*dstCoins - buyPrice*srcCoins

		if profit <= 0 {
			return lower, upper
		}
		if profitWithFee > 0 {
			upper.AskDepth = CalculateDepth(asks, srcCoins)
			upper.BidDepth = CalculateDepth(bids, dstCoins)
			upper.AskVolUsdt = srcCoins * buyPrice
			upper.BidVolUsdt = dstCoins * sellPrice
			upper.BuyPrice = buyPrice
			upper.SellPrice = sellPrice
			upper.Profit = profitWithFee
			upper.Spread = ((sellPrice - buyPrice) / buyPrice) * 100

			if !firstProfitFunded {
				firstProfitFunded = true
				lower.AskDepth = CalculateDepth(asks, srcCoins)
				lower.BidDepth = CalculateDepth(asks, dstCoins)
				lower.AskVolUsdt = srcCoins * buyPrice
				lower.BidVolUsdt = dstCoins * sellPrice
				lower.BuyPrice = buyPrice
				lower.SellPrice = sellPrice
				lower.Profit = profitWithFee
				lower.Spread = ((sellPrice - buyPrice) / buyPrice) * 100
			}
		}
	}
	fmt.Println(upper)
	fmt.Println(lower)
	return lower, upper
}

func CalcBestVolume(asks, bids *[]common.Order, thenMakerFee, withdrawFee float64) common.BestVolume {
	asksVolume := BookSum(asks)

	var result common.BestVolume

	var profit, profitWithFee float64
	var buyPrice, sellPrice float64

	for srcCoins := 1.0; srcCoins <= asksVolume; srcCoins++ {
		// coins
		dstCoins := (srcCoins - withdrawFee) - (srcCoins-withdrawFee)*thenMakerFee
		// usdt
		buyPrice = CalculatePrice(asks, srcCoins)
		sellPrice = CalculatePrice(bids, dstCoins)
		// usdt
		profit = sellPrice*srcCoins - buyPrice*srcCoins
		profitWithFee = sellPrice*dstCoins - buyPrice*srcCoins

		if profit <= 0 {
			return result
		}

		if profitWithFee > 0 && result.Profit < profitWithFee {
			result.AskDepth = CalculateDepth(asks, srcCoins)
			result.BidDepth = CalculateDepth(bids, dstCoins)
			result.AskVolUsdt = srcCoins * buyPrice
			result.BidVolUsdt = dstCoins * sellPrice
			result.BuyPrice = buyPrice
			result.SellPrice = sellPrice
			result.Profit = profitWithFee
			result.Spread = ((sellPrice - buyPrice) / buyPrice) * 100
		}
	}

	return result
}

func CalcVolume(asks, bids *[]common.Order, thenMakerFeem, withdrawFee float64) common.VolumeData {
	lower, upper := CalcMinMaxVolume(asks, bids, thenMakerFeem, withdrawFee)
	best := CalcBestVolume(asks, bids, thenMakerFeem, withdrawFee)
	return common.VolumeData{
		Low:  lower,
		Up:   upper,
		Best: best,
	}
}

// bookSum возвращает суммарный объем ордербука.
func BookSum(orderBook *[]common.Order) float64 {
	var res float64
	for _, order := range *orderBook {
		res += order.Amount
	}
	return res
}

// calculatePrice считает средневзвешенную цену для заданного объема.
func CalculatePrice(book *[]common.Order, coins float64) float64 {
	var sum float64
	var orderBookVol float64
	remainingVol := coins

	for _, order := range *book {
		orderBookVol += order.Amount

		if remainingVol <= order.Amount {
			sum += remainingVol * order.Price
			break
		}

		remainingVol -= order.Amount
		sum += order.Amount * order.Price
	}

	if orderBookVol < coins {
		return -1
	}

	return sum / coins
}

func CalculateDepth(orders *[]common.Order, coins float64) int {
	depth := 0
	remainingCoins := coins

	for _, order := range *orders {
		if remainingCoins <= 0 {
			break
		}

		if order.Amount >= remainingCoins {
			// Этот ордер полностью или частично покрывает оставшиеся монеты
			depth++
			break
		} else {
			// Используем все монеты из этого ордера
			remainingCoins -= order.Amount
			depth++
		}
	}

	return depth
}
