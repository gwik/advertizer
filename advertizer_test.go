package advertizer

import (
	"container/heap"
	"fmt"
	"testing"
)

// TestEventQueueTop tests that Top event is equal to next popped event.
func TestEventQueueTop(t *testing.T) {
	eq := &eventQueue{}
	heap.Init(eq)

	e1 := &event{id: 1, adv: 1, inc: 0}
	e2 := &event{id: 2, adv: 2, inc: 1}
	e3 := &event{id: 3, adv: 2, inc: 2}
	e4 := &event{id: 4, adv: 0, inc: 3}

	heap.Push(eq, e1)
	heap.Push(eq, e2)
	heap.Push(eq, e3)
	heap.Push(eq, e4)

	for _, v := range []*event{e4, e1, e2, e3} {
		e := eq.Top()

		if e != v {
			t.Errorf("e%d should have been on top, was: e%d", v.id, e.id)
		}

		if p := heap.Pop(eq).(*event); p != v {
			t.Errorf("e%d should have been popped, was: e%d", v.id, p.id)
		}
	}
}

func ExampleAdvertizer() {
	adv := New(2)

	adv.Push(0, "zero")
	adv.Push(1, "one")
	adv.Push(2, "two")
	adv.Push(3, "three")

	id, val, _ := adv.Advertize()
	fmt.Printf("%d: %s\n", id, val)

	id, val, _ = adv.Advertize()
	fmt.Printf("%d: %s\n", id, val)

	id, val, _ = adv.Advertize()
	fmt.Printf("%d: %s\n", id, val)

	id, val, _ = adv.Advertize()
	fmt.Printf("%d: %s\n", id, val)

	id, val, _ = adv.Advertize()
	fmt.Printf("%d: %s\n", id, val)

	id, val, _ = adv.Advertize()
	fmt.Printf("%d: %s\n", id, val)

	adv.Push(1, "one again")

	id, val, _ = adv.Advertize()
	fmt.Printf("%d: %s\n", id, val)

	id, val, _ = adv.Advertize()
	fmt.Printf("%d: %s\n", id, val)

	id, val, _ = adv.Advertize()
	fmt.Printf("%d: %s\n", id, val)

	id, val, _ = adv.Advertize()
	fmt.Printf("%d: %s\n", id, val)

	id, val, _ = adv.Advertize()
	fmt.Printf("%d: %v", id, val)

	// Output: 0: zero
	// 1: one
	// 2: two
	// 3: three
	// 0: zero
	// 1: one
	// 1: one again
	// 2: two
	// 3: three
	// 1: one again
	// 0: <nil>
}

func TestAdvertizerPushAdvertize(t *testing.T) {

	adv := New(10)

	adv.Push(0, nil)
	adv.Push(1, nil)
	adv.Push(2, nil)
	adv.Push(3, nil)

	if l := adv.Len(); l != 4 {
		t.Errorf("Lenght should be 4, was: %d", l)
	}

	events := []int64{0, 1, 2, 3}
	for i := 0; i < len(events)*3; i++ {
		j := i % 4
		id, _, _ := adv.Advertize()
		if id != events[j] {
			t.Errorf("Should have advertize e%d at run %d(%d), was: e%d", events[j], i, j, id)
		}
	}

	adv.Push(1, nil)
	adv.Push(2, nil)

	events = []int64{1, 2, 1, 2, 1, 2, 0, 3}
	for i, v := range events {
		id, _, _ := adv.Advertize()
		if id != v {
			t.Errorf("Should have advertize e%d at run %d, was: e%d", v, i, id)
		}
	}

	// TODO
	// - test insert after advertizing a few times
	// - test drop after max is reached

}

func TestAdvertizeRemove(t *testing.T) {
	adv := New(10)

	adv.Push(0, 0)
	adv.Push(1, 1)

	if val, ok := adv.Remove(1); !ok {
		t.Fatal("should have removed 1")
	} else {
		if val != 1 {
			t.Fatal("val should be returned")
		}
	}
}
