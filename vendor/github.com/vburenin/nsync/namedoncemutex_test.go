package nsync

import (
	"testing"
	"time"
)

func TestOnceMutex(t *testing.T) {
	nm := NewOnceMutex()

	if !nm.Lock() {
		t.Error("Didn't acquire lock")
	}
	go func() {
		time.Sleep(time.Millisecond * 10)
		nm.Unlock()
	}()

	if nm.Lock() {
		t.Error("Locked acquired twice!")
	}

	if nm.Lock() {
		t.Error("Locked acquired twice!")
	}
}

func TestNamedOnceMutex(t *testing.T) {
	nm := NewNamedOnceMutex()

	if !nm.Lock(1) {
		t.Error("Didn't acquire lock")
	}
	go func() {
		time.Sleep(time.Millisecond * 10)
		nm.Unlock(1)
	}()

	if nm.Lock(1) {
		t.Error("Locked acquired twice!")
	}

	if !nm.Lock(1) {
		t.Error("Not acquired again!")
	}
}
