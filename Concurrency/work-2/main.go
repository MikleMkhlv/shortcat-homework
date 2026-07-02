package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Semaphore struct {
	C chan struct{}
}

func (s *Semaphore) Acquire() {
	s.C <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.C
}

func worker(inputCh <-chan int, outputCh chan<- int, wg *sync.WaitGroup, sem *Semaphore, indexWorker int) {
	defer wg.Done()
	for item := range inputCh {
		sem.Acquire()
		fmt.Printf("worker-%d set to work\n", indexWorker)
		time.Sleep(time.Second * 2)
		outputCh <- item * 2
		fmt.Printf("worker-%d has finished its task\n", indexWorker)
		sem.Release()
	}
}

func main() {
	inputCh := make(chan int, 10)
	outputCh := make(chan int, 5)
	// outputCh := make(chan int)
	sem := Semaphore{
		C: make(chan struct{}, 5),
	}

	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		indexWorker := i + 1
		wg.Add(1)
		go worker(inputCh, outputCh, &wg, &sem, indexWorker)
	}

	go func() {
		for j := 0; j < 50; j++ {
			inputCh <- rand.Intn(50)
		}
		close(inputCh)
	}()

	go func() {
		wg.Wait()
		close(outputCh)
	}()

	var count int
	readDone := make(chan struct{})

	go func() {
		defer close(readDone)
		for val := range outputCh {
			fmt.Println(val)
			count++
		}
	}()

	<-readDone

	fmt.Printf("result=%d\n", count)
}
