package queue

import "testing"

func TestNewQueue(t *testing.T) {
	q := NewQueue[int](10)
	for i := 0; i < 10; i++ {
		q.Enqueue(i)
	}

	for i := 0; i < 10; i++ {
		item, _ := q.Dequeue()
		if item != i {
			t.Fatalf("Wrong item: expected %d got %d\n", i, item)
		}
	}
}

func TestBigQueue(t *testing.T) {
	q := NewQueue[int](10)
	for i := 0; i < 100; i++ {
		q.Enqueue(i)
	}

	for i := 0; i < 100; i++ {
		item, _ := q.Dequeue()
		if item != i {
			t.Fatalf("Wrong item: expected %d got %d\n", i, item)
		}
	}
}

func TestParallelQueue(t *testing.T) {
	q := NewQueue[int](10)
	go func() {
		for i := 0; i < 100; i++ {
			q.Enqueue(i)
		}
	}()

	for i := 0; i < 100; i++ {
		var item int
		var err error
		keep := true
		for keep {
			item, err = q.Dequeue()
			if err == nil {
				keep = false
			}
		}
		if item != i {
			t.Fatalf("Wrong item: expected %d got %d\n", i, item)
		}
	}

}
