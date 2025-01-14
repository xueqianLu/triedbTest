package testsuite

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Snapshot struct {
	Root    common.Hash
	Account types.StateAccount
}

type SnapshotSet struct {
	snapshots []Snapshot
}

func NewSnapshotSet() *SnapshotSet {
	return &SnapshotSet{
		snapshots: make([]Snapshot, 0, 1000),
	}
}

func (s *SnapshotSet) AddSnapshot(root common.Hash, account types.StateAccount) {
	s.snapshots = append(s.snapshots, Snapshot{root, account})
}

type VerifierFunc func(sp Snapshot) bool

func (s *SnapshotSet) RangeSnapshot(verifier VerifierFunc) (int, int) {
	var success, failed int
	for i := len(s.snapshots) - 1; i >= 0; i-- {
		if verifier(s.snapshots[i]) {
			success++
		} else {
			failed++
		}
	}
	return success, failed
}
