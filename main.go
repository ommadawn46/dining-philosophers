package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/ommadawn46/dining-philosophers/dining"
)

var actionFlag = flag.Int("action", 0, "select philosophers action")

func main() {
	philosopherNames := []string{"Socrates", "Plato", "Aristotle", "Kant", "Nietzsche"}
	stomachSize := 1000
	maxEatNS := 10000

	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	switch *actionFlag {
	case 0:
		dining.AllTakeRightFirst(philosopherNames, stomachSize, maxEatNS) // !!! this pattern will deadlock !!!
	case 1:
		dining.OneTakeLeftFirst(philosopherNames, stomachSize, maxEatNS)
	case 2:
		dining.AskWaiter(philosopherNames, stomachSize, maxEatNS)
	case 3:
		dining.ControlByMonitor(philosopherNames, stomachSize, maxEatNS)
	}
}
