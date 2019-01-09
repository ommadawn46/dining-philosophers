package dining

import (
	"fmt"
	"math/rand"
	"time"
)

// Philosopher can only eat when have both left and right forks
type Philosopher struct {
	name      string
	rightFork chan struct{}
	leftFork  chan struct{}

	ateAmount   int
	stomachSize int
	maxEatNS    int
	action      func(*Philosopher)
}

func (p *Philosopher) takeRightFork() {
	if PrintFlag {
		fmt.Println(p.name, "is taking a right fork.")
	}
	<-p.rightFork
}

func (p *Philosopher) takeLeftFork() {
	if PrintFlag {
		fmt.Println(p.name, "is taking a left fork.")
	}
	<-p.leftFork
}

func (p *Philosopher) returnRightFork() {
	if PrintFlag {
		fmt.Println(p.name, "is returning a right fork.")
	}
	p.rightFork <- struct{}{}
}

func (p *Philosopher) returnLeftFork() {
	if PrintFlag {
		fmt.Println(p.name, "is returning a left fork.")
	}
	p.leftFork <- struct{}{}
}

func (p *Philosopher) eat() {
	if PrintFlag {
		fmt.Printf("%s is eating (%d/%d).\n", p.name, p.ateAmount+1, p.stomachSize)
	}
	time.Sleep(time.Duration(rand.Intn(p.maxEatNS)) * time.Nanosecond)
	p.ateAmount++
}

func (p *Philosopher) isFull() bool {
	return p.ateAmount >= p.stomachSize
}

func (p *Philosopher) run(done chan struct{}) {
	for !p.isFull() {
		p.action(p)
	}
	if PrintFlag {
		fmt.Println(p.name, "finished eating.")
	}
	done <- struct{}{}
}

func (p *Philosopher) init() {
	p.ateAmount = 0
	select {
	case <-p.rightFork:
	default:
	}
	p.rightFork <- struct{}{}
}

func newPhilosopher(name string, stomachSize int, maxEatNS int, action func(*Philosopher)) *Philosopher {
	p := &Philosopher{
		name:        name,
		rightFork:   make(chan struct{}, 1),
		stomachSize: stomachSize,
		maxEatNS:    maxEatNS,
		action:      action,
	}
	p.init()
	return p
}
