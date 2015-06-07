// Package advertizer provides an advertizer datastructure.
package advertizer

import (
	"container/heap"
)

type event struct {
	id  int64
	val interface{}

	// reserved for the advertizer.
	adv int
	idx int
	inc int64
}

type eventQueue []*event

var _ heap.Interface = (*eventQueue)(nil)

func (eq eventQueue) Len() int { return len(eq) }

func (eq eventQueue) Less(i, j int) bool {
	if eq[i].adv == eq[j].adv {
		return eq[i].inc < eq[j].inc
	}
	return eq[i].adv < eq[j].adv
}

func (eq eventQueue) Swap(i, j int) {
	eq[i], eq[j] = eq[j], eq[i]
	eq[i].idx = i
	eq[j].idx = j
}

func (eq *eventQueue) Push(x interface{}) {
	n := len(*eq)
	event := x.(*event)
	event.idx = n
	*eq = append(*eq, event)
}

func (eq *eventQueue) Pop() interface{} {
	old := *eq
	n := len(old)
	event := old[n-1]
	event.idx = -1 // for safety
	*eq = old[0 : n-1]
	return event
}

func (eq eventQueue) Top() *event {
	if len(eq) < 1 {
		return nil
	}
	return eq[0]
}

// Advertizer is a data structure that allow to advertize an item for
// a number of times before it is dropped. The items are advertized
// in FIFO order for each round.
type Advertizer struct {
	max int
	m   map[int64]*event
	eq  *eventQueue
	inc int64
}

// New returns an initialized Advertizer that advertize n times.
func New(n int) *Advertizer {
	a := &Advertizer{
		max: n,
		m:   make(map[int64]*event),
		eq:  &eventQueue{},
	}
	heap.Init(a.eq)
	return a
}

// Len returns the number of events in the advertizer.
func (a *Advertizer) Len() int {
	return len(*a.eq)
}

// Push pushes a new item in the advertizing queue or updates an existing one,
// resetting its number of advertizing to 0.
func (a *Advertizer) Push(id int64, val interface{}) {
	a.inc++

	if e, ok := a.m[id]; ok {
		e.val = val
		e.adv = 0
		e.inc = a.inc
		heap.Fix(a.eq, e.idx)
		return
	}

	e := &event{
		id:  id,
		val: val,
		adv: 0,
		idx: -1,
		inc: a.inc,
	}
	a.m[id] = e
	heap.Push(a.eq, e)
}

// Advertize return the next item to advertize.
// The item is pushed back if advertize less than the defined
// maximum.
// It returns nil if there is no item to advertize.
func (a *Advertizer) Advertize() (int64, interface{}, bool) {
	e := a.eq.Top()
	if e == nil {
		return 0, nil, false
	}

	if e.adv+1 >= a.max { // Drop the event.
		heap.Remove(a.eq, e.idx) // Same as heap.Pop when e is at top.
		delete(a.m, e.id)
	} else { // Push down the event for later advertizing.
		e.adv++
		heap.Fix(a.eq, e.idx)
	}

	return e.id, e.val, true
}

// Remove removes an item with id.
func (a *Advertizer) Remove(id int64) (interface{}, bool) {
	if e, ok := a.m[id]; ok {
		delete(a.m, id)
		heap.Remove(a.eq, e.idx)
		return e.val, true
	}
	return nil, false
}
