package triedbTest

import "math/big"

type AssetInfo struct {
	Name     string
	Contract string
	Decimal  int
}

type Asset struct {
	Coin   string
	Amount *big.Int
}

type UserInfo struct {
	Address string
	Assets  []Asset
}

type OrderBook struct {
}

type Order struct {
}

type StatInfo struct {
}
