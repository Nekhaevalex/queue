package queue

import (
	"errors"
	"sync"
)

type Queue[T any] struct {
	mutex sync.Mutex
	queue []T
}

// Enqueue item
func (queue *Queue[T]) Enqueue(item T) {
	queue.mutex.Lock()
	defer func() {
		queue.mutex.Unlock()
	}()
	queue.queue = append(queue.queue, item)
}

// Dequeue item. Returns error if no items.
func (queue *Queue[T]) Dequeue() (T, error) {
	queue.mutex.Lock()
	defer func() {
		queue.mutex.Unlock()
	}()
	if len(queue.queue) > 0 {
		item := queue.queue[0]
		queue.queue = queue.queue[1:]
		return item, nil
	}
	var empty T
	return empty, errors.New("no new items")
}

// Returns true if there are any items in queue
func (queue *Queue[T]) HasItems() bool {
	queue.mutex.Lock()
	defer func() {
		queue.mutex.Unlock()
	}()
	return len(queue.queue) > 0
}

// Creates new queue with specified capacity and type T
func NewQueue[T any](capacity int) *Queue[T] {
	queue := new(Queue[T])
	queue.mutex = sync.Mutex{}
	queue.queue = make([]T, 0, capacity)
	return queue
}
