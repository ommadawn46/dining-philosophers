package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var actionFlag = flag.Int("action", 0, "select philosophers action")

type philosopher struct {
	name      string
	rightFork chan struct{}
	leftFork  chan struct{}

	ateAmount   int
	stomachSize int
	action      func(*philosopher)
}

func (p *philosopher) takeRightFork() {
	fmt.Println(p.name, "is taking a right fork.")
	<-p.rightFork
}

func (p *philosopher) takeLeftFork() {
	fmt.Println(p.name, "is taking a left fork.")
	<-p.leftFork
}

func (p *philosopher) returnRightFork() {
	fmt.Println(p.name, "is returning a right fork.")
	p.rightFork <- struct{}{}
}

func (p *philosopher) returnLeftFork() {
	fmt.Println(p.name, "is returning a left fork.")
	p.leftFork <- struct{}{}
}

func (p *philosopher) eat() {
	fmt.Println(p.name, "is eating.")
	time.Sleep(time.Duration(rand.Intn(10000)) * time.Nanosecond)
	p.ateAmount++
}

func (p *philosopher) isFull() bool {
	return p.ateAmount >= p.stomachSize
}

func (p *philosopher) run(done chan struct{}) {
	for !p.isFull() {
		p.action(p)
	}
	fmt.Println(p.name, "finished eating.")
	done <- struct{}{}
}

func newPhilosopher(name string, stomachSize int, action func(*philosopher)) *philosopher {
	fork := make(chan struct{}, 1)
	fork <- struct{}{}
	return &philosopher{
		name:        name,
		rightFork:   fork,
		ateAmount:   0,
		stomachSize: stomachSize,
		action:      action,
	}
}

func setupDining(names []string, stomachSize int, actions []func(*philosopher)) ([]*philosopher, error) {
	if len(names) != len(actions) {
		return nil, fmt.Errorf("names length and actions length must be same")
	}

	philosophers := []*philosopher{}
	for i := 0; i < len(names); i++ {
		philosophers = append(
			philosophers,
			newPhilosopher(
				names[i], stomachSize, actions[i],
			),
		)
	}
	for i, philo := range philosophers {
		if i == len(names)-1 {
			philosophers[0].leftFork = philo.rightFork
		} else {
			philosophers[i+1].leftFork = philo.rightFork
		}
	}

	return philosophers, nil
}

func startDining(philosophers []*philosopher) {
	startTime := time.Now()
	doneChannels := []chan struct{}{}
	for _, philo := range philosophers {
		done := make(chan struct{}, 1)
		doneChannels = append(doneChannels, done)
		go philo.run(done)
	}
	for _, done := range doneChannels {
		<-done
	}
	endTime := time.Now()

	fmt.Printf("dining finished in %.2f seconds.", (endTime.Sub(startTime)).Seconds())
}

func allTakeRightFirst(philosopherNames []string, stomachSize int) {
	actions := []func(*philosopher){}

	for i := 0; i < len(philosopherNames); i++ {
		actions = append(
			actions,
			func(p *philosopher) {
				p.takeRightFork()
				p.takeLeftFork()
				p.eat()
				p.returnLeftFork()
				p.returnRightFork()
			},
		)
	}

	philosophers, _ := setupDining(
		philosopherNames, stomachSize, actions,
	)
	startDining(philosophers)
}

func oneTakeLeftFirst(philosopherNames []string, stomachSize int) {
	actions := []func(*philosopher){
		func(p *philosopher) {
			p.takeLeftFork()
			p.takeRightFork()
			p.eat()
			p.returnRightFork()
			p.returnLeftFork()
		},
	}

	for i := 1; i < len(philosopherNames); i++ {
		actions = append(
			actions,
			func(p *philosopher) {
				p.takeRightFork()
				p.takeLeftFork()
				p.eat()
				p.returnLeftFork()
				p.returnRightFork()
			},
		)
	}

	philosophers, _ := setupDining(
		philosopherNames, stomachSize, actions,
	)
	startDining(philosophers)
}

func askWaiter(philosopherNames []string, stomachSize int) {
	waiter := make(chan struct{}, len(philosopherNames)-1)
	actions := []func(*philosopher){}

	for i := 0; i < len(philosopherNames); i++ {
		actions = append(
			actions,
			func(p *philosopher) {
				waiter <- struct{}{}
				p.takeRightFork()
				p.takeLeftFork()
				p.eat()
				p.returnLeftFork()
				p.returnRightFork()
				<-waiter
			},
		)
	}

	philosophers, _ := setupDining(
		philosopherNames, stomachSize, actions,
	)
	startDining(philosophers)
}

func controlByMonitor(philosopherNames []string, stomachSize int) {
	var mntr *monitor
	actions := []func(*philosopher){}

	for i := 0; i < len(philosopherNames); i++ {
		chairIdx := i
		actions = append(
			actions,
			func(p *philosopher) {
				mntr.pickup(chairIdx)
				p.eat()
				mntr.putdown(chairIdx)
			},
		)
	}

	philosophers, _ := setupDining(
		philosopherNames, stomachSize, actions,
	)
	mntr = newMonitor(philosophers)
	startDining(philosophers)
}

func main() {
	philosopherNames := []string{"Socrates", "Plato", "Aristotle", "Kant", "Nietzsche"}
	stomachSize := 1000

	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	switch *actionFlag {
	case 0:
		allTakeRightFirst(philosopherNames, stomachSize) // !!! this pattern will deadlock !!!
	case 1:
		oneTakeLeftFirst(philosopherNames, stomachSize)
	case 2:
		askWaiter(philosopherNames, stomachSize)
	case 3:
		controlByMonitor(philosopherNames, stomachSize)
	}
}
