package nsync

import (
	"sync"
	"testing"
	"time"
)

func TestCWaitGroup(t *testing.T) {
	cwg := NewControlWaitGroup(2)
	cwg.Do(func() { time.Sleep(time.Second) })
	cwg.Do(func() { time.Sleep(time.Second) })
	if cwg.Working() != 2 {
		t.Error("Two working thread should be running.")
	}
	cwg.Wait()
	if cwg.Working() != 0 {
		t.Error("No workers should be running")
	}

	cwg.Do(func() { time.Sleep(time.Second) })
	cwg.Do(func() { time.Sleep(time.Second) })
	go cwg.Do(func() { time.Sleep(time.Second) })

	time.Sleep(time.Millisecond * 200)
	if cwg.Waiting() != 1 {
		t.Error("1 waiting worker should be there")
	}
	if cwg.Working() != 2 {
		t.Error("Two working thread should be running.")
	}
	cwg.Wait()
	if cwg.Working() != 0 {
		t.Error("No workers should be running")
	}
}

func TestAbort(t *testing.T) {
	var mu sync.Mutex
	var a int
	f := func() {
		mu.Lock()
		a++
		mu.Unlock()
		time.Sleep(time.Second)
	}
	cwg := NewControlWaitGroup(2)
	cwg.Do(f)
	cwg.Do(f)
	go cwg.Do(f)
	time.Sleep(time.Millisecond * 200)
	cwg.Abort()
	cwg.Wait()
	if a != 2 {
		t.Errorf("Only two jobs should be completed. Actual: %d", a)
	}
}
