package req

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"mad-scanner.com/scriner/common"
	logger "mad-scanner.com/scriner/pkg/log"
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

type PoloniexBook struct {
	Asks []string
	Bids []string
}

func ConvertStringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

var (
	json = jsoniter.ConfigFastest
	log  = logger.GetLogger()
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

	return requestBytes(req, "kucoin")
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

	body, err := requestBytes(req, "mexc")
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

	body, err := requestBytes(req, "gate")
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

	body, err := requestBytes(req, "bybit")
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

	body, err := requestBytes(req, "bingx")
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

	body, err := requestBytes(req, "bitget")
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

	body, err := requestBytes(req, "htx")
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

	body, err := requestBytes(req, "okx")
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

	body, err := requestBytes(req, "poloniex")
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
