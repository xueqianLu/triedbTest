package types

type Coin struct {
	Name        string
	Contract    string
	Decimal     int
	CanWithdraw bool
	WithdrawFee string
}
