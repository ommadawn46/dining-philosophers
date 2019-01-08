package dining

import (
	"sync"
)

type state int

const (
	thinking state = iota
	hungry
	eating
)

type monitor struct {
	mutex        *sync.Mutex
	conds        []*sync.Cond
	states       []state
	philosophers []*philosopher
	n            int
}

func newMonitor(philosophers []*philosopher) *monitor {
	mutex := &sync.Mutex{}
	conds := []*sync.Cond{}
	states := []state{}
	n := len(philosophers)
	for i := 0; i < n; i++ {
		conds = append(conds, sync.NewCond(mutex))
		states = append(states, thinking)
	}
	return &monitor{mutex, conds, states, philosophers, n}
}

func (m *monitor) test(i int) {
	left, right := (i-1+m.n)%m.n, (i+1)%m.n
	if !m.philosophers[i].isFull() && m.states[left] != eating && m.states[right] != eating {
		m.states[i] = eating
		m.conds[i].Signal()
	}
}

func (m *monitor) pickup(i int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.states[i] = hungry
	m.test(i)
	if m.states[i] != eating {
		m.conds[i].Wait()
	}

	m.philosophers[i].takeRightFork()
	m.philosophers[i].takeLeftFork()
}

func (m *monitor) putdown(i int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.philosophers[i].returnLeftFork()
	m.philosophers[i].returnRightFork()

	m.states[i] = thinking

	left, right := (i-1+m.n)%m.n, (i+1)%m.n
	m.test(left)
	m.test(right)
}
