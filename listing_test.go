package main

import (
	"math/rand"
	"testing"
)

func BenchmarkGenListingInfo(b *testing.B) {
	b.N = 1000000
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var info templateInfo
			genListingInfoCRC(&info, rand.Uint64())
		}
	})
}
