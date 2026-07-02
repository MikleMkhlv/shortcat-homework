package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func worker(inputCh <-chan int, stopCh <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	timer := time.NewTimer(time.Millisecond * 200)
	for {
		select {
		case value, ok := <-inputCh:
			if !ok {
				fmt.Println("chan inputCh closed")
				return
			}
			time.Sleep(time.Millisecond * 200)
			fmt.Println(value)

			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
		case <-stopCh:
			fmt.Println("stop signal")
			return
		case <-timer.C:
			fmt.Println("stop with timeout")
			return
		}
	}
}

func main() {
	wg := sync.WaitGroup{}
	stopCh := make(chan struct{})
	intutCh := make(chan int, 10)

	go func() {
		for i := 0; i < 50; i++ {
			intutCh <- rand.IntN(50)
		}
		close(intutCh)
	}()

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go worker(intutCh, stopCh, &wg)
	}

	time.Sleep(time.Second * 5)
	close(stopCh)
	wg.Wait()
}
