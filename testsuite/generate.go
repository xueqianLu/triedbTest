package testsuite

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/uuid"
	"math/big"
	"math/rand"
)

func GenerateCustom(count int) (map[string]*types.StateAccount, map[string][]byte) {
	prefix := uuid.NewString()
	d := make(map[string]*types.StateAccount)
	dd := make(map[string][]byte)
	for i := 0; i < count; i++ {
		ref := uuid.NewString()
		addr := fmt.Sprintf("%s%d", prefix, i)
		d[addr] = nil
		dd[addr] = []byte(ref)
	}
	return d, dd
}

func AccountData(acc *types.StateAccount) []byte {
	v, _ := rlp.EncodeToBytes(acc)
	return v
}

func GenerateAccount(count int) (map[string]*types.StateAccount, map[string][]byte) {
	randomPrefix := uuid.NewString()
	d := make(map[string]*types.StateAccount)
	dd := make(map[string][]byte)
	for i := 0; i < count; i++ {
		addr := fmt.Sprintf("%s%d", randomPrefix, i)
		acc := &types.StateAccount{
			CodeHash: []byte{},
			Root:     common.HexToHash(fmt.Sprintf("%x", i+500)),
			Nonce:    uint64(rand.Intn(1000000)) + 10022,
			Balance:  big.NewInt(12339990),
		}
		d[addr] = acc
		dd[addr] = AccountData(acc)
	}
	return d, dd
}
