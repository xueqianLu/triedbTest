package types

type OrderBook struct {
	Symbol              string
	BaseCoin            string
	BaseCoinMinDecimal  int
	QuoteCoin           string
	QuoteCoinMinDecimal int
	Published           bool
	Disabled            bool
}
