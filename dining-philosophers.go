package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var eatN = 1000
var philosopherNames = []string{"Socrates", "Plato", "Aristotle", "Kant", "Nietzsche"}

var actionFlag = flag.Int("action", 0, "select philosophers action")

type philosopher struct {
	name      string
	rightFork chan struct{}
	leftFork  chan struct{}

	action func(*philosopher)
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
}

func (p *philosopher) run(done chan struct{}) {
	for i := 0; i < eatN; i++ {
		p.action(p)
	}
	done <- struct{}{}
}

func newPhilosopher(name string, action func(*philosopher)) *philosopher {
	fork := make(chan struct{}, 1)
	fork <- struct{}{}
	return &philosopher{
		name:      name,
		rightFork: fork,
		action:    action,
	}
}

func setupDining(names []string, actions []func(*philosopher)) ([]*philosopher, error) {
	if len(names) != len(actions) {
		return nil, fmt.Errorf("names length and actions length must be same")
	}

	philosophers := []*philosopher{}
	for i := 0; i < len(names); i++ {
		philosophers = append(
			philosophers,
			newPhilosopher(
				names[i], actions[i],
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

func allTakeRightFirst() {
	takeRightForkFirst := func(p *philosopher) {
		p.takeRightFork()
		p.takeLeftFork()
		p.eat()
		p.returnLeftFork()
		p.returnRightFork()
	}
	pholosophers, _ := setupDining(
		philosopherNames,
		[]func(*philosopher){
			takeRightForkFirst,
			takeRightForkFirst,
			takeRightForkFirst,
			takeRightForkFirst,
			takeRightForkFirst,
		},
	)
	startDining(pholosophers)
}

func oneTakeLeftFirst() {
	takeRightForkFirst := func(p *philosopher) {
		p.takeRightFork()
		p.takeLeftFork()
		p.eat()
		p.returnLeftFork()
		p.returnRightFork()
	}
	takeLeftForkFirst := func(p *philosopher) {
		p.takeLeftFork()
		p.takeRightFork()
		p.eat()
		p.returnRightFork()
		p.returnLeftFork()
	}
	pholosophers, _ := setupDining(
		philosopherNames,
		[]func(*philosopher){
			takeLeftForkFirst,
			takeRightForkFirst,
			takeRightForkFirst,
			takeRightForkFirst,
			takeRightForkFirst,
		},
	)
	startDining(pholosophers)
}

func allAskWaiter() {
	waiter := make(chan struct{}, len(philosopherNames)-1)
	waitWaitersOk := func(p *philosopher) {
		waiter <- struct{}{}
		p.takeRightFork()
		p.takeLeftFork()
		p.eat()
		p.returnLeftFork()
		p.returnRightFork()
		<-waiter
	}
	pholosophers, _ := setupDining(
		philosopherNames,
		[]func(*philosopher){
			waitWaitersOk,
			waitWaitersOk,
			waitWaitersOk,
			waitWaitersOk,
			waitWaitersOk,
		},
	)
	startDining(pholosophers)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	flag.Parse()
	switch *actionFlag {
	case 0:
		allTakeRightFirst() // !!! this pattern will deadlock !!!
	case 1:
		oneTakeLeftFirst()
	case 2:
		allAskWaiter()
	}
}
