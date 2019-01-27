package workerpool

import (
	"testing"
)

func TestWorkerlist(t *testing.T) {
	a := &worker{}
	b := &worker{}
	c := &worker{}

	var wl workerlist

	if wl.Front() != nil || wl.Back() != nil {
		t.Fatal("Front and Back should be nil")
	}

	_, ok := wl.PopBack()
	if ok {
		t.Fatal(ok)
	}

	wl.PushBack(a)
	if wl.Front().value != a || wl.Back().value != a {
		t.Fatal("Front and Back should be a")
	}

	if wl.Front().prev != nil || wl.Back().next != nil {
		t.Fatal(wl.Front().prev, wl.Back().next)
	}

	wl.PushBack(b)
	if wl.Front().value != a || wl.Back().value != b {
		t.Fatal("Front should be a and Back should be b")
	}

	wl.PushBack(c)
	if wl.Front().value != a || wl.Back().value != c {
		t.Fatal("Front should be a and Back should be c")
	}

	v, ok := wl.PopBack()
	if !ok || v != c {
		t.Fatal("wrong return value")
	}
	if wl.Front().value != a || wl.Back().value != b {
		t.Fatal("Front should be a and Back should be b")
	}

	v, ok = wl.PopBack()
	if !ok || v != b {
		t.Fatal("wrong return value", v, ok)
	}
	if wl.Front().value != a || wl.Back().value != a {
		t.Fatal("Front should be a and Back should be b")
	}
	wl.Back()
	v, ok = wl.PopBack()
	if !ok || v != a {
		t.Fatal("wrong return value", v, ok)
	}
	if wl.Front() != nil || wl.Back() != nil {
		t.Fatal("Front and Back should be nil", wl.Front(), wl.Back())
	}

	_, ok = wl.PopBack()
	if ok {
		t.Fatal()
	}

	wl.PushBack(a)
	if wl.Front().value != a || wl.Back().value != a {
		t.Fatal("Front and Back should be a")
	}

	wl.PushBack(b)
	wl.PushBack(c)

	wl.ResetFront(wl.Front().next)
	x, _ := wl.PopBack()
	y, _ := wl.PopBack()
	_, ok = wl.PopBack()

	if x != c || y != b || ok != false{
		t.Fatal()
	}
}
