package nsync

import (
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	l := NewNamedMutex()
	l.Lock("test1")
	l.Lock("test2")
	if len(l.mutexMap) != 2 {
		t.Error("Unexpected number of channels.")
	}
	if len(l.mutexMap["test1"]) != 1 {
		t.Error("Lock is not acquired")
	}
	if len(l.mutexMap["test2"]) != 1 {
		t.Error("Lock is not acquired")
	}
	l.Unlock("test1")
	l.Unlock("test2")

	if len(l.mutexMap["test1"]) != 0 {
		t.Error("Lock is not released")
	}
	if len(l.mutexMap["test2"]) != 0 {
		t.Error("Lock is not released")
	}
}

func TestLock2(t *testing.T) {
	l := NewNamedMutex()
	l.Lock("test1")
	go func() {
		time.Sleep(time.Millisecond * 10)
		l.Unlock("test1")
	}()
	l.Lock("test1")
	l.Unlock("test1")
}

func TestTryLock(t *testing.T) {
	l := NewNamedMutex()
	if !l.TryLock("test1") {
		t.Error("Didn't acquire lock")
	}
	if !l.TryLock("test2") {
		t.Error("Didn't acquire lock")
	}

	if l.TryLock("test1") {
		t.Error("Lock should be acquired")
	}
	if l.TryLock("test2") {
		t.Error("Lock should be acquired")
	}
}

func TestTryLockTimeout(t *testing.T) {
	l := NewNamedMutex()
	if !l.TryLockTimeout("test1", time.Millisecond) {
		t.Error("Didn't acquire lock")
	}
	if !l.TryLockTimeout("test2", time.Millisecond) {
		t.Error("Didn't acquire lock")
	}

	if l.TryLockTimeout("test1", time.Millisecond) {
		t.Error("Lock should be acquired")
	}
	if l.TryLockTimeout("test2", time.Millisecond) {
		t.Error("Lock should be acquired")
	}
}
