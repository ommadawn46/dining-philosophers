package dining

import (
	"fmt"
	"time"
)

// PrintFlag (true: do print, false: no print)
var PrintFlag = true

var philosopherNames = []string{
	"Socrates", "Plato", "Aristotle", "Kant", "Nietzsche",
	"Confucius", "Averroes", "Buddha", "Abelard", "Adorno",
	"Bacon", "Barthes", "Bataille", "Baudrillard", "Beauvoir",
	"Benjamin", "Berkeley", "Butler", "Camus", "Chomsky",
	"Cixous", "Deleuze", "Derrida", "Descartes", "Dewey",
	"Foucault", "Gadamer", "Habermas", "Haraway", "Hegel",
	"Heidegger", "Hobbes", "Hume", "Husserl", "Irigaray",
	"James", "Immanuel", "Kristeva", "Tzu", "Levinas",
	"Locke", "Lyotard", "Ponty", "Mill", "Moore",
	"Quine", "Rand", "Rousseau", "Sartre", "Schopenhauer",
	"Spinoza", "Wittgenstein", "Aquinas", "Arendt", "Augustine",
}

func setupDining(philosopherN int, stomachSize int, maxEatNS int, actions []func(*Philosopher)) ([]*Philosopher, error) {
	if philosopherN != len(actions) {
		return nil, fmt.Errorf("names length and actions length must be same")
	}
	if philosopherN > len(philosopherNames) {
		return nil, fmt.Errorf("philosopherN must be smaller than %d", len(philosopherNames))
	}

	philosophers := []*Philosopher{}
	for i := 0; i < philosopherN; i++ {
		philosophers = append(
			philosophers,
			newPhilosopher(
				philosopherNames[i], stomachSize, maxEatNS, actions[i],
			),
		)
	}
	for i, philo := range philosophers {
		if i == philosopherN-1 {
			philosophers[0].leftFork = philo.rightFork
		} else {
			philosophers[i+1].leftFork = philo.rightFork
		}
	}

	return philosophers, nil
}

// StartDining starts philosophers' dining
func StartDining(philosophers []*Philosopher) {
	for _, philo := range philosophers {
		philo.init()
	}

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

	if PrintFlag {
		fmt.Printf("dining finished in %.2f seconds.", (endTime.Sub(startTime)).Seconds())
	}
}

// PrepareAllRightFirst prepare all philosophers to take right fork first
func PrepareAllRightFirst(philosopherN int, stomachSize int, maxEatNS int) ([]*Philosopher, error) {
	actions := []func(*Philosopher){}

	for i := 0; i < philosopherN; i++ {
		actions = append(
			actions,
			func(p *Philosopher) {
				p.takeRightFork()
				p.takeLeftFork()
				p.eat()
				p.returnLeftFork()
				p.returnRightFork()
			},
		)
	}

	return setupDining(
		philosopherN, stomachSize, maxEatNS, actions,
	)
}

// PrepareOneLeftFirst prepare a philosopher to take left fork first
func PrepareOneLeftFirst(philosopherN int, stomachSize int, maxEatNS int) ([]*Philosopher, error) {
	actions := []func(*Philosopher){
		func(p *Philosopher) {
			p.takeLeftFork()
			p.takeRightFork()
			p.eat()
			p.returnRightFork()
			p.returnLeftFork()
		},
	}

	for i := 1; i < philosopherN; i++ {
		actions = append(
			actions,
			func(p *Philosopher) {
				p.takeRightFork()
				p.takeLeftFork()
				p.eat()
				p.returnLeftFork()
				p.returnRightFork()
			},
		)
	}

	return setupDining(
		philosopherN, stomachSize, maxEatNS, actions,
	)
}

// PrepareAskWaiter prepare philosophers to ask the waiter
func PrepareAskWaiter(philosopherN int, stomachSize int, maxEatNS int) ([]*Philosopher, error) {
	waiter := make(chan struct{}, philosopherN-1)
	actions := []func(*Philosopher){}

	for i := 0; i < philosopherN; i++ {
		actions = append(
			actions,
			func(p *Philosopher) {
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

	return setupDining(
		philosopherN, stomachSize, maxEatNS, actions,
	)
}

// PrepareControlByMonitor prepare philosophers controlled by the monitor
func PrepareControlByMonitor(philosopherN int, stomachSize int, maxEatNS int) ([]*Philosopher, error) {
	var mntr *monitor
	actions := []func(*Philosopher){}

	for i := 0; i < philosopherN; i++ {
		chairIdx := i
		actions = append(
			actions,
			func(p *Philosopher) {
				mntr.pickup(chairIdx)
				p.eat()
				mntr.putdown(chairIdx)
			},
		)
	}

	philosophers, err := setupDining(
		philosopherN, stomachSize, maxEatNS, actions,
	)
	mntr = newMonitor(philosophers)
	return philosophers, err
}
