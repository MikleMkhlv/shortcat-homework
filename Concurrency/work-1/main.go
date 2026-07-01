package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func withDataRace() {
	var counter int
	before := time.Now()
	for i := 0; i < 1000; i++ {
		go func() {
			counter++
		}()
	}
	after := time.Now()
	fmt.Println(counter)
	fmt.Printf("result=%d: time=%v\n", counter, after.Sub(before))
}

func withoutDataRace() {
	var counter int
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	before := time.Now()
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()
	after := time.Now()
	mu.Lock()
	fmt.Printf("result=%d: time=%v\n", counter, after.Sub(before))
	mu.Unlock()
}

func withAtomic() {
	var counter int64
	wg := sync.WaitGroup{}
	before := time.Now()
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}
	wg.Wait()
	after := time.Now()
	fmt.Printf("result=%d: time=%v\n", counter, after.Sub(before))
}

func main() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		// withDataRace()
		// withoutDataRace()
		withAtomic()
	}()

	wg.Wait()
	fmt.Println("Done")
}
