package utils_test

import (
	"quote_book/pkg/utils"
	"sync"
	"testing"
)

func TestIDGenerator_SequentialIDs(t *testing.T) {
	startID := 100
	gen := utils.NewIDGenerator(startID)

	for i := 0; i < 10; i++ {
		got := gen.GetID()
		want := startID + i
		if got != want {
			t.Errorf("GetID() = %d; want %d", got, want)
		}
	}
}

func TestIDGenerator_ConcurrentIDs(t *testing.T) {
	gen := utils.NewIDGenerator(0)
	const n = 1000

	var wg sync.WaitGroup
	ids := make(chan int, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := gen.GetID()
			ids <- id
		}()
	}

	wg.Wait()
	close(ids)

	seen := make(map[int]bool, n)
	for id := range ids {
		if id < 0 || id >= n {
			t.Errorf("ID out of range: %d", id)
		}
		if seen[id] {
			t.Errorf("Duplicate ID detected: %d", id)
		}
		seen[id] = true
	}

	if len(seen) != n {
		t.Errorf("Expected %d unique IDs, got %d", n, len(seen))
	}
}
