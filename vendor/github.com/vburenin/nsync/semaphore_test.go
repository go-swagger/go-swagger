package nsync

import (
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	s := NewSemaphore(2)
	s.Acquire()
	s.Acquire()
	if s.Value() != 2 {
		t.Error("Unexpected size")
	}
	go func() {
		time.Sleep(time.Millisecond * 20)
		s.Release()
	}()
	s.Acquire()
	if s.Value() != 2 {
		t.Error("Unexpected size")
	}
}

func TestSemaphoreTry(t *testing.T) {
	s := NewSemaphore(2)
	if !s.TryAcquire() {
		t.Error("Failed to acquire")
	}
	if !s.TryAcquire() {
		t.Error("Failed to acquire")
	}
	if s.TryAcquire() {
		t.Error("Should not acuire!!!")
	}
}

func TestSemaphoreTryTimeout(t *testing.T) {
	s := NewSemaphore(2)
	if !s.TryAcquire() {
		t.Error("Failed to acquire")
	}
	if !s.TryAcquire() {
		t.Error("Failed to acquire")
	}
	st := time.Now().UnixNano() / 1000000
	if s.TryAcquireTimeout(time.Millisecond * 15) {
		t.Error("Should not acuire!!!")
	}
	d := (time.Now().UnixNano() / 1000000) - st
	if d < 10 {
		t.Error("Error in timeout")
	}
	if s.Value() != 2 {
		t.Error("Unexpected size")
	}
}
