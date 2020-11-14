package deque

import ()

// minCapacity is the smallest capacity that deque may have
// must be power of 2 for bitwise modules: x % n == x & (n - 1)
const minCapacity = 16

// Deque represents a single instance of the deque data structure
// 这是一个不带空头节点、带空尾节点的环形队列
type Deque struct {
	buf    []interface{}
	head   int
	tail   int
	count  int
	minCap int
}

// Len returns the number of the elements currently stored in the queue.
func (q *Deque) Len() int {
	return q.count
}

// resize resizes the deque to fix exactly twice its current contents.
// this is used to grow the queue when it is full, and also to shrink it
// when it is only a quarter full
func (q *Deque) resize() {
	newBuf := make([]interface{}, q.count<<1)
	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		// more element
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

// growIfFull resizes up if the buffer is full.
func (q *Deque) growIfFull() {
	if len(q.buf) == 0 {
		if q.minCap == 0 {
			q.minCap = minCapacity
		}
		q.buf = make([]interface{}, q.minCap)
		return
	}
	if q.count == len(q.buf) {
		q.resize()
	}
}

// shrinkIfExcess resize down if the buffer 1/4 full
func (q *Deque) shrinkIfExcess() {
	if len(q.buf) > q.minCap && (q.count<<2) == len(q.buf) {
		q.resize()
	}
}

// next returns the next buffer position wrapping around
func (q *Deque) next(i int) int {
	return (i + 1) & (len(q.buf) - 1)
}

// prev returns the previous buffer position wrapping around buffer
func (q *Deque) prev(i int) int {
	return (i - 1) & (len(q.buf) - 1)
}

// SetMinCapacity sets a minimum capacity of 2^minCapacityExp.
// If the value of the minimum is less than or equal to the minmum allowed,
// then capacity is set to the minimum allowed.
// This may be called at anytime to set a new minimum caoacity.
// Setting a larger minimum capacity may be used to prevent resizing when the
// number of stored items changes frequently across a wide range.
func (q *Deque) SetMinCapacity(minCapacityExp uint) {
	if 1<<minCapacityExp > minCapacity {
		q.minCap = 1 << minCapacityExp
	} else {
		q.minCap = minCapacity
	}
}

// Clear removes all elements from the queue, but retains the current capacity.
// This is usefull when repeatedly resuing the queue at hight frequency to avoid
// GC during reuse. The queue will not be resized smaller as long as items are only added.
// Only when items are removed is the queue subject to getting resized smaller
func (q *Deque) Clear() {
	modBits := len(q.buf) - 1
	for h := q.head; h != q.tail; h = (h + 1) & modBits {
		q.buf[h] = nil
	}
	q.head = 0
	q.tail = 0
	q.tail = 0
}

// At returns the element at index i in the queue without removing the element
// from the queue. This method accepts only non-negative index values.
// At(0) refers to the first element and is the same as Front(). At(len()-1) refers
// to the last element and is the same as Back(). If the index is invalid, the call panics.
func (q *Deque) At(i int) interface{} {
	if i < 0 || i >= q.count {
		panic("deque: At() called with index out of range")
	}
	// bitwise modules
	return q.buf[(q.head+i)&(len(q.buf)-1)]
}

// Front returns the element at the front of the queue.
// This is the element that would be returned by popfront.
// This call panic if the queue is empty
func (q *Deque) Front() interface{} {
	if q.count <= 0 {
		panic("deque: Front() called when empty")
	}
	return q.buf[q.head]
}

// PopBack removes and returns the element from the back of the queue.
// Implements LIFO when used with PushBack(). If the queue is empty,
// the call panics.
func (q *Deque) PopBack() interface{} {
	if q.count <= 0 {
		panic("deque: PopBack() called on empty queue")
	}
	// calculate new tail position
	q.tail = q.prev(q.tail)
	// remove value at tail
	ret := q.buf[q.tail]
	q.buf[q.tail] = nil
	q.count--

	q.shrinkIfExcess()
	return ret
}

// PopFront removes and retuens the element from the front of the queue.
// Implements FIFO when used with PushBack(). If the queue is empty,
// the call panics.
func (q *Deque) PopFront() interface{} {
	if q.count <= 0 {
		panic("deque: PopFront() called on empty queue")
	}
	ret := q.buf[q.head]
	q.buf[q.head] = nil
	q.head = q.next(q.head)
	q.count--
	q.shrinkIfExcess()

	return ret
}

// PushFront prepends an element to the front of the queue
func (q *Deque) PushFront(elem interface{}) {
	q.growIfFull()
	q.head = q.prev(q.head)
	q.buf[q.head] = elem
	q.count++
}

// PushBack appends an element to the back of the queue. implements FIFO when
// elements are removed with PopFront(), and LIFO when elements are removed with PopBack
func (q *Deque) PushBack(elem interface{}) {
	q.growIfFull()
	q.buf[q.tail] = elem
	q.tail = q.next(q.tail)
	q.count++
}
