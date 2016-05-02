package nsync

import (
	"testing"
	"time"
)

func TestTryMutexLock(t *testing.T) {
	l := NewTryMutex()
	l.Lock()
	if len(l.c) != 1 {
		t.Error("Failed to acquire lock")
	}
	l.Unlock()
}

func TestTryMutexTryLock(t *testing.T) {
	l := NewTryMutex()
	if !l.TryLock() {
		t.Error("Lock must be acquired")
	}
	if l.TryLock() {
		t.Error("Lock must not be acquired")
	}
}

func TestTryMutexTryLock2(t *testing.T) {
	l := NewTryMutex()
	l.Lock()
	go func() {
		time.Sleep(time.Millisecond * 10)
		l.Unlock()
	}()
	l.Lock()
	l.Unlock()
}

func TestTryMutexLockTimeout(t *testing.T) {
	l := NewTryMutex()
	l.Lock()
	st := time.Now().UnixNano()
	l.TryLockTimeout(time.Millisecond * 10)
	et := time.Now().UnixNano()
	if et-st < int64(time.Millisecond*9) {
		t.Error("Wrong timeout")
	}
}
