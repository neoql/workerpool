package workerpool

type element struct {
	prev  *element
	next  *element
	value *worker
}

type workerlist struct {
	root element
}

func (wl *workerlist) Front() *element {
	return wl.root.next
}

func (wl *workerlist) Back() *element {
	return wl.root.prev
}

func (wl *workerlist) PushBack(w *worker) {
	ele := &element{prev: nil, next: nil, value: w}
	back := wl.Back()

	if back == nil {
		wl.root.prev = ele
		wl.root.next = ele
		return
	}
	back.next = ele
	ele.prev = back

	wl.root.prev = ele
}

func (wl *workerlist) PopBack() (w *worker, ok bool) {
	ele := wl.Back()
	if ele != nil {
		w = ele.value
		ok = true

		if wl.root.prev == wl.root.next {
			wl.root.prev = nil
			wl.root.next = nil
		} else {
			ele.prev.next = nil
			wl.root.prev = ele.prev
			ele.prev = nil
		}
	}
	return
}

func (wl *workerlist) ResetFront(e *element) {
	if e == nil {
		return
	}

	prev := e.prev
	if prev != nil {
		prev.next = nil
		e.prev = nil
	}

	wl.root.next = e
}
