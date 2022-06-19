package queue

import "sync"

type Queue struct {
	first *queueItem
	last  *queueItem
	m     sync.Mutex
}

type queueItem struct {
	next *queueItem
	data interface{}
}

func (q *Queue) AddToHead(newItem interface{}) {
	q.m.Lock()
	defer q.m.Unlock()
	q.first = &queueItem{
		next: q.first,
		data: newItem,
	}
	if q.last == nil {
		q.last = q.first
	}
}

func (q *Queue) AddToTail(newItem interface{}) {
	q.m.Lock()
	defer q.m.Unlock()
	if q.last == nil {
		q.last = &queueItem{
			data: newItem,
		}
		q.first = q.last
		return
	}
	q.last.next = &queueItem{
		data: newItem,
	}
}

func (q *Queue) Get() interface{} {
	q.m.Lock()
	defer q.m.Unlock()
	defer func() {
		q.first = q.first.next
	}()
	return q.first.data
}

func (q *Queue) List() []interface{} {
	q.m.Lock()
	defer q.m.Unlock()
	res := make([]interface{}, 0)
	for next := q.first; next != nil; next = next.next {
		res = append(res, next.data)
	}
	return res
}
