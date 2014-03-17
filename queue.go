package gogadgets

import (
	"sync"
)

type queuenode struct {
	data *Message
	next *queuenode
}

type Queue struct {
	head  *queuenode
	tail  *queuenode
	count int
	lock  *sync.Mutex
	cond  *sync.Cond
}

func NewQueue() *Queue {
	q := &Queue{}
	q.lock = &sync.Mutex{}
	q.cond = sync.NewCond(&sync.Mutex{})
	return q
}

func (q *Queue) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.count
}

func (q *Queue) Push(item *Message) {
	q.lock.Lock()
	defer q.lock.Unlock()
	n := &queuenode{data: item}

	if q.tail == nil {
		q.tail = n
		q.head = n
	} else {
		q.tail.next = n
		q.tail = n
	}
	q.count++
	q.cond.Signal()
}

func (q *Queue) Get() *Message {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.head == nil {
		return nil
	}
	n := q.head
	q.head = n.next

	if q.head == nil {
		q.tail = nil
	}
	q.count--
	return n.data
	
}

func (q *Queue) Wait() {
	q.cond.Wait()
}

func (q *Queue) Lock() {
	q.cond.L.Lock()
}

func (q *Queue) Unlock() {
	q.cond.L.Unlock()
}
