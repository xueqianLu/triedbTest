package cosmos

import "testing"

func BenchIavlCommit(b *testing.B) {
	tree := newIVAL(b.TempDir(), false)
	defer tree.Close()

	for i := 0; i < b.N; i++ {
		tree.Commit()
	}
}
