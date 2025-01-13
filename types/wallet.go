package types

import "math/big"

type Wallet struct {
	User          string
	Coin          string
	Balance       *big.Int
	FrozenBalance *big.Int
}
