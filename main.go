package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ommadawn46/dining-philosophers/dining"
)

var philosopherNFlag = flag.Int("n", 5, "number of philosophers")
var stomachSizeFlag = flag.Int("stomach", 1000, "philosopher's stomach size")
var maxEatNSFlag = flag.Int("eatns", 5000, "time required for eating")
var solutionFlag = flag.Int("solution", 0, "dining philosopher problem solution")

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	philosopherN := *philosopherNFlag
	stomachSize := *stomachSizeFlag
	maxEatNS := *maxEatNSFlag

	var philosophers []*dining.Philosopher
	var err error
	switch *solutionFlag {
	case 0:
		// !!! this solution will deadlock !!!
		philosophers, err = dining.PrepareAllRightFirst(philosopherN, stomachSize, maxEatNS)
	case 1:
		philosophers, err = dining.PrepareOneLeftFirst(philosopherN, stomachSize, maxEatNS)
	case 2:
		philosophers, err = dining.PrepareAskWaiter(philosopherN, stomachSize, maxEatNS)
	case 3:
		philosophers, err = dining.PrepareControlByMonitor(philosopherN, stomachSize, maxEatNS)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dining.StartDining(philosophers)
}
