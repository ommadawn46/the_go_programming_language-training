package popcount

import (
	"testing"
)

var temp int

func BenchmarkPopCountTable(b *testing.B) {
	temp = 0
	for i := uint64(0); i < uint64(b.N); i++ {
		temp += PopCountTable(i)
	}
}

func BenchmarkPopCountLoop(b *testing.B) {
	temp = 0
	for i := uint64(0); i < uint64(b.N); i++ {
		temp += PopCountLoop(i)
	}
}

func BenchmarkPopCountBitShift(b *testing.B) {
	temp = 0
	for i := uint64(0); i < uint64(b.N); i++ {
		temp += PopCountBitShift(i)
	}
}

func BenchmarkPopCountLowerBitClear(b *testing.B) {
	temp = 0
	for i := uint64(0); i < uint64(b.N); i++ {
		temp += PopCountLowerBitClear(i)
	}
}
