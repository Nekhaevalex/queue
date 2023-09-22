package queue

type ChannelQueue[T any] struct {
	reserve  *Queue[T]
	queue    chan T
	capacity int
}

type sink uint8

const (
	channel sink = iota
	reserve
)

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func (queue *ChannelQueue[T]) getToCopyAndSink() (int, sink) {
	queueFreeSpace := queue.capacity - len(queue.queue)
	reserveLength := func() int {
		defer queue.reserve.mutex.Unlock()
		queue.reserve.mutex.Lock()
		return len(queue.reserve.queue)
	}()
	toCopy := min(queueFreeSpace, reserveLength)

	queueSpaceWillRemain := queueFreeSpace > toCopy
	sink := func() sink {
		if queueSpaceWillRemain {
			return channel
		}
		return reserve
	}()

	return toCopy, sink
}

func (queue *ChannelQueue[T]) Enqueue(item T) {
	toCopy, sink := queue.getToCopyAndSink()
	for i := 0; i < toCopy; i++ {
		item, err := queue.reserve.Dequeue()
		if err != nil {
			return
		}
		queue.queue <- item
	}
	switch sink {
	case channel:
		queue.queue <- item
	case reserve:
		queue.reserve.Enqueue(item)
	}
}

func (queue *ChannelQueue[T]) Dequeue() <-chan T {
	toCopy, _ := queue.getToCopyAndSink()
	for i := 0; i < toCopy; i++ {
		item, _ := queue.reserve.Dequeue()
		queue.queue <- item
	}
	return queue.queue
}

func (queue *ChannelQueue[T]) IsEmpty() bool {
	return len(queue.queue) == 0 && queue.reserve.IsEmpty()
}

func NewChannelQueue[T any](capacity int) *ChannelQueue[T] {
	queue := new(ChannelQueue[T])
	queue.reserve = NewQueue[T](capacity)
	queue.capacity = capacity
	queue.queue = make(chan T, capacity)
	return queue
}
