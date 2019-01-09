package dining

import (
	"testing"
)

func init() {
	PrintFlag = false
}

var stomachSize = 100
var maxEatNS = 1000000

func BenchmarkOneLeftFirst5(b *testing.B) {
	philosophers, _ := PrepareOneLeftFirst(5, stomachSize, maxEatNS)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StartDining(philosophers)
	}
}

func BenchmarkOneLeftFirst20(b *testing.B) {
	philosophers, _ := PrepareOneLeftFirst(20, stomachSize, maxEatNS)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StartDining(philosophers)
	}
}

func BenchmarkOneLeftFirst50(b *testing.B) {
	philosophers, _ := PrepareOneLeftFirst(50, stomachSize, maxEatNS)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StartDining(philosophers)
	}
}

func BenchmarkAskWaiter5(b *testing.B) {
	philosophers, _ := PrepareAskWaiter(5, stomachSize, maxEatNS)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StartDining(philosophers)
	}
}

func BenchmarkAskWaiter20(b *testing.B) {
	philosophers, _ := PrepareAskWaiter(20, stomachSize, maxEatNS)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StartDining(philosophers)
	}
}

func BenchmarkAskWaiter50(b *testing.B) {
	philosophers, _ := PrepareAskWaiter(50, stomachSize, maxEatNS)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StartDining(philosophers)
	}
}

func BenchmarkControlByMonitor5(b *testing.B) {
	philosophers, _ := PrepareControlByMonitor(5, stomachSize, maxEatNS)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StartDining(philosophers)
	}
}

func BenchmarkControlByMonitor20(b *testing.B) {
	philosophers, _ := PrepareControlByMonitor(20, stomachSize, maxEatNS)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StartDining(philosophers)
	}
}

func BenchmarkControlByMonitor50(b *testing.B) {
	philosophers, _ := PrepareControlByMonitor(50, stomachSize, maxEatNS)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StartDining(philosophers)
	}
}
