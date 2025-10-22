package req

import "mad-scanner.com/scriner/common"

// parseStringPairs converts slices like [][]string{{price, amount}, ...} into []common.Order
// It silently skips malformed entries.
func parseStringPairs(p [][]string) []common.Order {
	var res []common.Order
	for _, pair := range p {
		if len(pair) < 2 {
			continue
		}
		price, err := ConvertStringToFloat64(pair[0])
		if err != nil {
			continue
		}
		amount, err := ConvertStringToFloat64(pair[1])
		if err != nil {
			continue
		}
		res = append(res, common.Order{Price: price, Amount: amount})
	}
	return res
}

// parseFlatStringPairs converts a flat []string like [price1,amount1,price2,amount2,...]
// into []common.Order. It stops on first parse error.
func parseFlatStringPairs(p []string) []common.Order {
	var res []common.Order
	for i := 0; i+1 < len(p); i += 2 {
		price, err := ConvertStringToFloat64(p[i])
		if err != nil {
			break
		}
		amount, err := ConvertStringToFloat64(p[i+1])
		if err != nil {
			break
		}
		res = append(res, common.Order{Price: price, Amount: amount})
	}
	return res
}

// parseFloat64Pairs converts slices like [][]float64{{price, amount}, ...} into []common.Order
func parseFloat64Pairs(p [][]float64) []common.Order {
	var res []common.Order
	for _, pair := range p {
		if len(pair) < 2 {
			continue
		}
		res = append(res, common.Order{Price: pair[0], Amount: pair[1]})
	}
	return res
}
