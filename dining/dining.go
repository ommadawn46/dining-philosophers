package dining

import (
	"fmt"
	"time"
)

func setupDining(names []string, stomachSize int, maxEatNS int, actions []func(*philosopher)) ([]*philosopher, error) {
	if len(names) != len(actions) {
		return nil, fmt.Errorf("names length and actions length must be same")
	}

	philosophers := []*philosopher{}
	for i := 0; i < len(names); i++ {
		philosophers = append(
			philosophers,
			newPhilosopher(
				names[i], stomachSize, maxEatNS, actions[i],
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

func AllTakeRightFirst(philosopherNames []string, stomachSize int, maxEatNS int) {
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
		philosopherNames, stomachSize, maxEatNS, actions,
	)
	startDining(philosophers)
}

func OneTakeLeftFirst(philosopherNames []string, stomachSize int, maxEatNS int) {
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
		philosopherNames, stomachSize, maxEatNS, actions,
	)
	startDining(philosophers)
}

func AskWaiter(philosopherNames []string, stomachSize int, maxEatNS int) {
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
		philosopherNames, stomachSize, maxEatNS, actions,
	)
	startDining(philosophers)
}

func ControlByMonitor(philosopherNames []string, stomachSize int, maxEatNS int) {
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
		philosopherNames, stomachSize, maxEatNS, actions,
	)
	mntr = newMonitor(philosophers)
	startDining(philosophers)
}
