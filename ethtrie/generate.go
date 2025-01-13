package ethtrie

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"math/big"
)

func generateAccount(count int) map[string]*types.StateAccount {
	randomPrefix := uuid.NewString()
	d := make(map[string]*types.StateAccount)
	for i := 0; i < count; i++ {
		addr := fmt.Sprintf("%s%d", randomPrefix, i)
		d[addr] = &types.StateAccount{
			CodeHash: []byte{},
			Root:     common.HexToHash(fmt.Sprintf("%x", i+500)),
			Nonce:    10022,
			Balance:  big.NewInt(12339990),
		}
	}
	return d
}
