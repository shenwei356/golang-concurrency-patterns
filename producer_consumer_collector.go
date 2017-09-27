package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sort"
	"sync"
	"time"
)

func main() {
	ncpus := runtime.NumCPU()
	runtime.GOMAXPROCS(ncpus)

	// collector
	out := make(chan int, ncpus)
	done := make(chan int)
	go func() {
		buffer := make(map[int]int)
		id := 1
		for i := range out {
			buffer[i] = i
			// fmt.Printf("read from out: %d, buffer: %v\n", i, buffer)
			if _, ok := buffer[id]; ok {
				fmt.Println("  result of", id)
				delete(buffer, id)
				id++
			}
		}
		// sort
		keys := make([]int, len(buffer))
		i := 0
		for k := range buffer {
			keys[i] = k
			i++
		}
		sort.Ints(keys)
		// fmt.Println("remaining data in buffer:", keys)
		for _, k := range keys {
			fmt.Println("  result of", buffer[k])
		}
		done <- 1
	}()

	// producer
	in := make(chan int, ncpus)
	go func() {
		for i := 1; i < 10; i++ {
			in <- i
		}
		close(in)
	}()

	// comsumer
	var wg sync.WaitGroup
	// tokens := make(chan int, ncpus) # for cases jobs not coming from chan
	for i := range in {
		// fmt.Println("read from in", i)
		// tokens <- 1
		wg.Add(1)
		go func(i int) {
			defer func() {
				wg.Done()
				// <-tokens
			}()
			fmt.Println("work with", i)
			time.Sleep(time.Millisecond * time.Duration(rand.Int31n(10000)))
			fmt.Println("done with", i)
			out <- i
		}(i)
	}
	wg.Wait()
	close(out)
	<-done // wait reducer
}
