package queue

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestSimpleCase(t *testing.T) {
	queue := NewChannelQueue[int](10)
	for i := 0; i < 10; i++ {
		queue.Enqueue(i)
	}
	for i := 0; i < 10; i++ {
		item := <-queue.Dequeue()
		if item != i {
			t.Fatalf("Wrong item: expected %d got %d\n", i, item)
		}
	}
}

func TestReserveCase(t *testing.T) {
	queue := NewChannelQueue[int](10)
	for i := 0; i < 15; i++ {
		queue.Enqueue(i)
	}
	for i := 0; i < 15; i++ {
		item := <-queue.Dequeue()
		if item != i {
			t.Fatalf("Wrong item: expected %d got %d\n", i, item)
		}
	}
}

func TestParallelCase(t *testing.T) {
	queue := NewChannelQueue[int](10)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 0; i < 100; i++ {
			queue.Enqueue(i)
		}
		wg.Done()
	}()
	for i := 0; i < 100; i++ {
		item := <-queue.Dequeue()
		if item != i {
			t.Fatalf("Wrong item: expected %d got %d\n", i, item)
		}
	}
	wg.Wait()
}

func BenchmarkMultiParallelCase(t *testing.B) {
	queue := NewChannelQueue[int](1024)
	timer := time.After(10 * time.Second)
	for th := 0; th < 100; th++ {
		go func(thnum int) {
			locTimer := time.Tick(100 * time.Millisecond)
			for {
				<-locTimer
				queue.Enqueue(10)
				log.Printf("[%d] Added to queue: %d/%d, Reserve: %d\n", thnum, len(queue.queue), queue.capacity, len(queue.reserve.queue))
			}

		}(th)
	}
	loop := true
	for loop {
		select {
		case <-timer:
			loop = false
		default:
			item := <-queue.Dequeue()
			log.Printf("Removed from queue: %d/%d, Reserve: %d\n", len(queue.queue), queue.capacity, len(queue.reserve.queue))
			if item < 0 {
				t.Fatalf("Wrong item: got %d\n", item)
			}
		}
	}
}
