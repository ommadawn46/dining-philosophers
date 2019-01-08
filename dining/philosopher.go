package dining

import (
	"fmt"
	"math/rand"
	"time"
)

type philosopher struct {
	name      string
	rightFork chan struct{}
	leftFork  chan struct{}

	ateAmount   int
	stomachSize int
	maxEatNS    int
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
	time.Sleep(time.Duration(rand.Intn(p.maxEatNS)) * time.Nanosecond)
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

func newPhilosopher(name string, stomachSize int, maxEatNS int, action func(*philosopher)) *philosopher {
	fork := make(chan struct{}, 1)
	fork <- struct{}{}
	return &philosopher{
		name:        name,
		rightFork:   fork,
		ateAmount:   0,
		stomachSize: stomachSize,
		maxEatNS:    maxEatNS,
		action:      action,
	}
}
