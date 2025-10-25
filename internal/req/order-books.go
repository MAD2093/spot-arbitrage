package req

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"

	"mad-scanner.com/scriner/common"
)

type MexcBook struct {
	Bids [][]string
	Asks [][]string
}

type GateBook struct {
	Bids [][]string
	Asks [][]string
}

type BybyitBook struct {
	Result struct {
		Asks [][]string `json:"a"`
		Bids [][]string `json:"b"`
	} `json:"result"`
}

type BingxBook struct {
	Data struct {
		Asks [][]string `json:"asks"`
		Bids [][]string `json:"bids"`
	} `json:"data"`
}

type BitrueBook struct {
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}

func (ob *BitrueBook) UnmarshalJSON(data []byte) error {
	var raw struct {
		LastUpdateId int64           `json:"lastUpdateId"`
		Bids         [][]interface{} `json:"bids"`
		Asks         [][]interface{} `json:"asks"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for _, bid := range raw.Bids {
		if len(bid) >= 2 {
			ob.Bids = append(ob.Bids, []string{
				fmt.Sprintf("%v", bid[0]), // цена
				fmt.Sprintf("%v", bid[1]), // объем
			})
		}
	}

	for _, ask := range raw.Asks {
		if len(ask) >= 2 {
			ob.Asks = append(ob.Asks, []string{
				fmt.Sprintf("%v", ask[0]),
				fmt.Sprintf("%v", ask[1]),
			})
		}
	}

	return nil
}

type BitmartBook struct {
	Data struct {
		Asks [][]string `json:"asks"`
		Bids [][]string `json:"bids"`
	} `json:"data"`
}

type BinanceBook struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
}

type CoinwBook struct {
	Data struct {
		Asks [][]string `json:"asks"`
		Bids [][]string `json:"bids"`
	} `json:"data"`
}

type HtxBook struct {
	Tick struct {
		Bids [][]float64 `json:"bids"`
		Asks [][]float64 `json:"asks"`
	} `json:"tick"`
}

type OkxBook struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Asks [][]string `json:"asks"`
		Bids [][]string `json:"bids"`
	} `json:"data"`
}

type DigifinexBook struct {
	Asks [][]float64 `json:"asks"`
	Bids [][]float64 `json:"bids"`
}

type CoinexBook struct {
	Data struct {
		Depth struct {
			Asks [][]string `json:"asks"`
			Bids [][]string `json:"bids"`
		} `json:"depth"`
	} `json:"data"`
}

type AscendexBook struct {
	Data struct {
		Data struct {
			Asks [][]string `json:"asks"`
			Bids [][]string `json:"bids"`
		} `json:"data"`
	} `json:"data"`
}

type XtBook struct {
	Result struct {
		Asks [][]string `json:"asks"`
		Bids [][]string `json:"bids"`
	} `json:"result"`
}

type PoloniexBook struct {
	Asks []string
	Bids []string
}

func ConvertStringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

var (
	json = jsoniter.ConfigFastest
)

func kucoin_makeRequest(endpoint, symbol string) ([]byte, error) {
	apiKey := "67126b080c4baa0001196c03"
	apiSecret := "fd55d631-5888-4444-bed5-7e1376dd7b5a"
	password := "KarimPetuh13"

	signer := newKcSigner(apiKey, apiSecret, password)
	headers := signer.Headers("GET/api/v3/market/orderbook/level2_100?symbol=" + symbol)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("KUCOIN-makeRequest:: error creating request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("KC-API-KEY", headers[0])
	req.Header.Add("KC-API-SIGN", headers[3])
	req.Header.Add("KC-API-TIMESTAMP", headers[2])
	req.Header.Add("KC-API-PASSPHRASE", headers[1])

	client := GetExchangeClient("KUCOIN")

	return client.MakeRequest(req)
}

func Mexc_get_orderbook(symbol string) common.OrderBook {
	var Orderbook MexcBook
	var result common.OrderBook

	endpoint := "https://api.mexc.com/api/v3/depth?limit=100&symbol=" + strings.ReplaceAll(symbol, "/", "")

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Err(err).Str("exchange", "MEXC").Msg("error creating request")
		return common.OrderBook{}
	}

	client := GetExchangeClient("MEXC")

	body, err := client.MakeRequest(req)
	if err != nil {
		log.Err(err).Str("exchange", "MEXC").Msg("err make request")
		return common.OrderBook{}
	}

	err = json.Unmarshal(body, &Orderbook)
	if err != nil {
		log.Err(err).Str("exchange", "MEXC").Msg("error unmarshaling config data")
	}

	result.Bids = parseStringPairs(Orderbook.Bids)
	result.Asks = parseStringPairs(Orderbook.Asks)

	return result
}

func Gate_get_orderbook(symbol string) common.OrderBook {
	var Orderbook GateBook
	var result common.OrderBook

	endpoint := "https://api.gateio.ws/api/v4/spot/order_book?limit=100&currency_pair=" + strings.ReplaceAll(symbol, "/", "_")

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Err(err).Str("exchange", "GATE").Msg("error creating request")
		return common.OrderBook{}
	}

	client := GetExchangeClient("GATE")

	body, err := client.MakeRequest(req)
	if err != nil {
		log.Err(err).Str("exchange", "GATE").Msg("error making request")
		return common.OrderBook{}
	}

	err = json.Unmarshal(body, &Orderbook)
	if err != nil {
		log.Err(err).Str("exchange", "GATE").Msg("error unmarshaling config data")
	}

	result.Bids = parseStringPairs(Orderbook.Bids)
	result.Asks = parseStringPairs(Orderbook.Asks)

	return result
}

func Bybyit_get_orderbook(symbol string) common.OrderBook {
	var Orderbook BybyitBook
	var result common.OrderBook

	endpoint := "https://api.bybit.com/v5/market/orderbook?limit=100&category=spot&symbol=" + strings.ReplaceAll(symbol, "/", "")

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Err(err).Str("exchange", "BYBIT").Msg("error creating request")
		return common.OrderBook{}
	}

	client := GetExchangeClient("BYBIT")

	body, err := client.MakeRequest(req)
	if err != nil {
		log.Err(err).Str("exchange", "BYBIT").Msg("error making request")
		return common.OrderBook{}
	}

	err = json.Unmarshal(body, &Orderbook)
	if err != nil {
		log.Err(err).Str("exchange", "BYBIT").Msg("error unmarshaling config data")
	}

	result.Bids = parseStringPairs(Orderbook.Result.Bids)
	result.Asks = parseStringPairs(Orderbook.Result.Asks)

	return result
}

func Bingx_get_orderbook(symbol string) common.OrderBook {
	var Orderbook BingxBook
	var result common.OrderBook

	endpoint := "https://open-api.bingx.com/openApi/swap/v2/quote/depth?limit=100&symbol=" + strings.ReplaceAll(symbol, "/", "-")

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Err(err).Str("exchange", "BINGX").Msg("error creating request")
		return common.OrderBook{}
	}

	client := GetExchangeClient("BINGX")

	body, err := client.MakeRequest(req)
	if err != nil {
		log.Err(err).Str("exchange", "BINGX").Msg("error making request")
		return common.OrderBook{}
	}

	err = json.Unmarshal(body, &Orderbook)
	if err != nil {
		log.Err(err).Str("exchange", "BINGX").Msg("error unmarshaling config data")
	}

	result.Bids = parseStringPairs(Orderbook.Data.Bids)
	result.Asks = parseStringPairs(Orderbook.Data.Asks)

	return result
}

func Bitget_get_orderbook(symbol string) common.OrderBook {
	var Orderbook BingxBook
	var result common.OrderBook

	endpoint := "https://api.bitget.com/api/v2/spot/market/orderbook?limit=100&symbol=" + strings.ReplaceAll(symbol, "/", "")

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Err(err).Str("exchange", "BITGET").Msg("error creating request")
		return common.OrderBook{}
	}

	client := GetExchangeClient("BITGET")

	body, err := client.MakeRequest(req)
	if err != nil {
		log.Err(err).Str("exchange", "BITGET").Msg("error making request")
		return common.OrderBook{}
	}

	err = json.Unmarshal(body, &Orderbook)
	if err != nil {
		log.Err(err).Str("exchange", "BITGET").Msg("error unmarshaling config data")
	}

	result.Bids = parseStringPairs(Orderbook.Data.Bids)
	result.Asks = parseStringPairs(Orderbook.Data.Asks)

	return result
}

func Htx_get_orderbook(symbol string) common.OrderBook {
	var Orderbook HtxBook
	var result common.OrderBook

	endpoint := "https://api.huobi.pro/market/depth?type=step0&symbol=" + strings.ToLower(strings.ReplaceAll(symbol, "/", ""))

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Err(err).Str("exchange", "HTX").Msg("error creating request")
		return common.OrderBook{}
	}

	client := GetExchangeClient("HTX")

	body, err := client.MakeRequest(req)
	if err != nil {
		log.Err(err).Str("exchange", "HTX").Msg("error making request")
		return common.OrderBook{}
	}

	err = json.Unmarshal(body, &Orderbook)
	if err != nil {
		log.Err(err).Str("exchange", "HTX").Msg("error unmarshaling config data")
	}

	result.Asks = parseFloat64Pairs(Orderbook.Tick.Asks)
	result.Bids = parseFloat64Pairs(Orderbook.Tick.Bids)

	return result
}

func Okx_get_orderbook(symbol string) common.OrderBook {
	var Orderbook OkxBook
	var result common.OrderBook

	endpoint := "https://www.okx.com/api/v5/market/books?sz=100&instId=" + strings.ReplaceAll(symbol, "/", "-")

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Err(err).Str("exchange", "OKX").Msg("error creating request")
		return common.OrderBook{}
	}

	client := GetExchangeClient("OKX")

	body, err := client.MakeRequest(req)
	if err != nil {
		log.Err(err).Str("exchange", "OKX").Msg("error making request")
		return common.OrderBook{}
	}

	err = json.Unmarshal(body, &Orderbook)
	if err != nil {
		log.Err(err).Str("exchange", "OKX").Msg("error unmarshaling config data")
	}

	if len(Orderbook.Data) == 0 {
		log.Err(err).Str("exchange", "OKX").Str("code", Orderbook.Code).Str("msg", Orderbook.Msg).Msg("order book is nil")
		return result
	}

	data := Orderbook.Data[0]
	result.Bids = parseStringPairs(data.Bids)
	result.Asks = parseStringPairs(data.Asks)

	return result
}

func Poloniex_get_orderbook(symbol string) common.OrderBook {
	var Orderbook PoloniexBook
	var result common.OrderBook

	endpoint := fmt.Sprintf("https://api.poloniex.com/markets/%s/orderBook?limit=100", strings.ReplaceAll(symbol, "/", "_"))

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Err(err).Str("exchange", "POLONIEX").Msg("error creating request")
		return common.OrderBook{}
	}

	client := GetExchangeClient("POLONIEX")

	body, err := client.MakeRequest(req)
	if err != nil {
		log.Err(err).Str("exchange", "POLONIEX").Msg("error making request")
		return common.OrderBook{}
	}

	err = json.Unmarshal(body, &Orderbook)
	if err != nil {
		log.Err(err).Str("exchange", "POLONIEX").Msg("error unmarshaling config data")
	}

	result.Bids = parseFlatStringPairs(Orderbook.Bids)
	result.Asks = parseFlatStringPairs(Orderbook.Asks)

	return result
}

func Kucoin_get_orderbook(symbol string) common.OrderBook {
	var Orderbook BingxBook
	var result common.OrderBook

	symbol = strings.ReplaceAll(symbol, "/", "-")
	endpoint := "https://api.kucoin.com/api/v3/market/orderbook/level2_100?symbol=" + symbol

	body, err := kucoin_makeRequest(endpoint, symbol)
	if err != nil {
		log.Err(err).Str("exchange", "KUCOIN").Msg("error making request")
		return common.OrderBook{}
	}

	err = json.Unmarshal(body, &Orderbook)
	if err != nil {
		log.Err(err).Str("exchange", "KUCOIN").Msg("error unmarshaling config data")
		return common.OrderBook{}
	}

	for count, bid := range Orderbook.Data.Bids {
		price, err := ConvertStringToFloat64(bid[0])
		if err != nil {
			log.Err(err).Str("exchange", "KUCOIN").Msg("error converting bid price")
			continue
		}
		amount, err := ConvertStringToFloat64(bid[1])
		if err != nil {
			log.Err(err).Str("exchange", "KUCOIN").Msg("error converting bid amount")
			continue
		}
		result.Bids = append(result.Bids, common.Order{Price: price, Amount: amount})
		if count == 4 {
			break
		}
	}

	for count, ask := range Orderbook.Data.Asks {
		price, err := ConvertStringToFloat64(ask[0])
		if err != nil {
			log.Err(err).Str("exchange", "KUCOIN").Msg("error converting ask price")
			continue
		}
		amount, err := ConvertStringToFloat64(ask[1])
		if err != nil {
			log.Err(err).Str("exchange", "KUCOIN").Msg("error converting ask amount")
			continue
		}
		result.Asks = append(result.Asks, common.Order{Price: price, Amount: amount})
		if count == 4 {
			break
		}
	}

	return result
}

func Bitmart_get_orderbook(symbol string) common.OrderBook {
	var result common.OrderBook
	var Orderbook BitmartBook

	symbol = strings.ReplaceAll(symbol, "/", "_")
	endpoint := fmt.Sprintf("https://api-cloud.bitmart.com/spot/quotation/v3/books?symbol=%s&limit=50", symbol)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Err(err).Str("exchange", "BITMART").Msg("error creating request")
		return common.OrderBook{}
	}

	client := GetExchangeClient("BITMART")

	body, err := client.MakeRequest(req)
	if err != nil {
		log.Err(err).Str("exchange", "BITMART").Msg("error making request")
		return common.OrderBook{}
	}

	err = json.Unmarshal(body, &Orderbook)
	if err != nil {
		log.Err(err).Str("exchange", "BITMART").Msg("error unmarshaling config data")
		return common.OrderBook{}
	}

	for count, bid := range Orderbook.Data.Bids {
		price, err := ConvertStringToFloat64(bid[0])
		if err != nil {
			log.Err(err).Str("exchange", "BITMART").Msg("error converting bid price")
			continue
		}
		amount, err := ConvertStringToFloat64(bid[1])
		if err != nil {
			log.Err(err).Str("exchange", "BITMART").Msg("error converting bid amount")
			continue
		}
		result.Bids = append(result.Bids, common.Order{Price: price, Amount: amount})
		if count == 4 {
			break
		}
	}

	for count, ask := range Orderbook.Data.Asks {
		price, err := ConvertStringToFloat64(ask[0])
		if err != nil {
			log.Err(err).Str("exchange", "BITMART").Msg("error converting ask price")
			continue
		}
		amount, err := ConvertStringToFloat64(ask[1])
		if err != nil {
			log.Err(err).Str("exchange", "BITMART").Msg("error converting ask amount")
			continue
		}
		result.Asks = append(result.Asks, common.Order{Price: price, Amount: amount})
		if count == 4 {
			break
		}
	}

	return result
}

func Binance_get_orderbook(symbol string) common.OrderBook {
	var result common.OrderBook
	var Orderbook BinanceBook

	symbol = strings.ReplaceAll(symbol, "/", "")
	endpoint := fmt.Sprintf("https://api3.binance.com/api/v3/depth?limit=500&symbol=%s", symbol)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Err(err).Str("exchange", "BINANCE").Msg("error creating request")
		return common.OrderBook{}
	}

	client := GetExchangeClient("BINANCE")

	body, err := client.MakeRequest(req)
	if err != nil {
		log.Err(err).Str("exchange", "BINANCE").Msg("error making request")
		return common.OrderBook{}
	}

	err = json.Unmarshal(body, &Orderbook)
	if err != nil {
		log.Err(err).Str("exchange", "BINANCE").Msg("error unmarshaling config data")
		return common.OrderBook{}
	}

	result.Bids = parseStringPairs(Orderbook.Bids)
	result.Asks = parseStringPairs(Orderbook.Asks)

	return result
}

func Coinex_get_orderbook(symbol string) common.OrderBook {
	var result common.OrderBook
	var Orderbook CoinexBook

	symbol = strings.ReplaceAll(symbol, "/", "")
	endpoint := fmt.Sprintf("https://api.coinex.com/v2/spot/depth?market=%s&limit=50&interval=0", symbol)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Err(err).Str("exchange", "COINEX").Msg("error creating request")
		return common.OrderBook{}
	}

	client := GetExchangeClient("COINEX")

	body, err := client.MakeRequest(req)
	if err != nil {
		log.Err(err).Str("exchange", "COINEX").Msg("error making request")
		return common.OrderBook{}
	}

	err = json.Unmarshal(body, &Orderbook)
	if err != nil {
		log.Err(err).Str("exchange", "COINEX").Msg("error unmarshaling config data")
		return common.OrderBook{}
	}

	result.Bids = parseStringPairs(Orderbook.Data.Depth.Bids)
	result.Asks = parseStringPairs(Orderbook.Data.Depth.Asks)

	return result
}

func Coinw_get_orderbook(symbol string) common.OrderBook {
	var result common.OrderBook
	var Orderbook CoinwBook

	// https://api.coinw.com/api/v1/public?command=returnOrderBook&symbol=BTC_USDT&size=100
	symbol = strings.ReplaceAll(symbol, "/", "_")
	endpoint := fmt.Sprintf("https://api.coinw.com/api/v1/public?command=returnOrderBook&symbol=%s&size=100", symbol)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Err(err).Str("exchange", "COINW").Msg("error creating request")
		return common.OrderBook{}
	}

	client := GetExchangeClient("COINW")

	body, err := client.MakeRequest(req)
	if err != nil {
		log.Err(err).Str("exchange", "COINW").Msg("error making request")
		return common.OrderBook{}
	}

	err = json.Unmarshal(body, &Orderbook)
	if err != nil {
		log.Err(err).Str("exchange", "COINW").Msg("error unmarshaling config data")
		return common.OrderBook{}
	}

	result.Bids = parseStringPairs(Orderbook.Data.Bids)
	result.Asks = parseStringPairs(Orderbook.Data.Asks)

	return result
}

func Xt_get_orderbook(symbol string) common.OrderBook {
	var result common.OrderBook
	var Orderbook XtBook

	// https://sapi.xt.com/v4/public/depth?symbol=BTC_USDT&limit=500
	symbol = strings.ReplaceAll(symbol, "/", "_")
	endpoint := fmt.Sprintf("https://sapi.xt.com/v4/public/depth?symbol=%s&limit=500", symbol)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Err(err).Str("exchange", "XT").Msg("error creating request")
		return common.OrderBook{}
	}

	client := GetExchangeClient("XT")

	body, err := client.MakeRequest(req)
	if err != nil {
		log.Err(err).Str("exchange", "XT").Msg("error making request")
		return common.OrderBook{}
	}

	err = json.Unmarshal(body, &Orderbook)
	if err != nil {
		log.Err(err).Str("exchange", "XT").Msg("error unmarshaling config data")
		return common.OrderBook{}
	}

	result.Bids = parseStringPairs(Orderbook.Result.Bids)
	result.Asks = parseStringPairs(Orderbook.Result.Asks)

	return result
}

func Ascendex_get_orderbook(symbol string) common.OrderBook {
	var result common.OrderBook
	var Orderbook AscendexBook

	// https://ascendex.com/api/pro/v1/depth?symbol=ASD/USDT
	endpoint := fmt.Sprintf("https://ascendex.com/api/pro/v1/depth?symbol=%s&limit=100", symbol)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Err(err).Str("exchange", "ASCENDEX").Msg("error creating request")
		return common.OrderBook{}
	}

	client := GetExchangeClient("ASCENDEX")

	body, err := client.MakeRequest(req)
	if err != nil {
		log.Err(err).Str("exchange", "ASCENDEX").Msg("error making request")
		return common.OrderBook{}
	}

	err = json.Unmarshal(body, &Orderbook)
	if err != nil {
		log.Err(err).Str("exchange", "ASCENDEX").Msg("error unmarshaling config data")
		return common.OrderBook{}
	}

	result.Bids = parseStringPairs(Orderbook.Data.Data.Bids)
	result.Asks = parseStringPairs(Orderbook.Data.Data.Asks)

	return result
}

func Digifinex_get_orderbook(symbol string) common.OrderBook {
	var result common.OrderBook
	var Orderbook DigifinexBook

	// https://openapi.digifinex.com/v3/order_book?symbol=BTC_USDT&limit=150
	symbol = strings.ReplaceAll(symbol, "/", "_")
	endpoint := fmt.Sprintf("https://openapi.digifinex.com/v3/order_book?symbol=%s&limit=150", symbol)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Err(err).Str("exchange", "DIGIFINEX").Msg("error creating request")
		return common.OrderBook{}
	}

	client := GetExchangeClient("DIGIFINEX")

	body, err := client.MakeRequest(req)
	if err != nil {
		log.Err(err).Str("exchange", "DIGIFINEX").Msg("error making request")
		return common.OrderBook{}
	}

	err = json.Unmarshal(body, &Orderbook)
	if err != nil {
		log.Err(err).Str("exchange", "DIGIFINEX").Msg("error unmarshaling config data")
		return common.OrderBook{}
	}

	result.Bids = parseFloat64Pairs(Orderbook.Bids)
	result.Asks = parseFloat64Pairs(Orderbook.Asks)

	return result
}

func Bitrue_get_orderbook(symbol string) common.OrderBook {
	var result common.OrderBook
	var Orderbook BitrueBook

	// https://openapi.bitrue.com/api/v1/depth?limit=500&symbol=BTCUSDT
	symbol = strings.ReplaceAll(symbol, "/", "")
	endpoint := fmt.Sprintf("https://openapi.bitrue.com/api/v1/depth?limit=500&symbol=%s", symbol)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Err(err).Str("exchange", "BITRUE").Msg("error creating request")
		return common.OrderBook{}
	}

	client := GetExchangeClient("BITRUE")

	body, err := client.MakeRequest(req)
	if err != nil {
		log.Err(err).Str("exchange", "BITRUE").Msg("error making request")
		return common.OrderBook{}
	}

	err = Orderbook.UnmarshalJSON(body)
	if err != nil {
		log.Err(err).Str("exchange", "BITRUE").Msg("error unmarshaling config data")
		return common.OrderBook{}
	}

	result.Bids = parseStringPairs(Orderbook.Bids)
	result.Asks = parseStringPairs(Orderbook.Asks)

	return result
}
