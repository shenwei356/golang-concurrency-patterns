package main

import (
	"runtime"
	"sync"
	"testing"
)

func NoCancellation(n int) {
	var wg sync.WaitGroup

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < n; j++ {
				if n == int(n/3) {
					return
				}
			}
		}()
	}

	wg.Wait()
}

// func CheckWithSelectChannel0(n int) {
// 	var wg sync.WaitGroup
//
// 	done := make(chan struct{})
// 	var once sync.Once
//
// 	for i := 0; i < runtime.NumCPU(); i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
//
// 			select {
// 			case <-done:
// 				return
// 			default:
// 			}
//
// 			for j := 0; j < n; j++ {
// 				select {
// 				case <-done:
// 					return
// 				default:
// 				}
//
// 				if n == int(n/3) {
// 					once.Do(func() { close(done) })
// 					return
// 				}
// 			}
// 		}()
// 	}
//
// 	wg.Wait()
// }

func CheckWithSelectChannel(n int) {
	var wg sync.WaitGroup

	done := make(chan struct{})

	toStop := make(chan struct{}, 1)
	doneDone := make(chan struct{})
	go func() {
		<-toStop
		close(done)
		doneDone <- struct{}{}
	}()

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			select {
			case <-done:
				return
			default:
			}

			for j := 0; j < n; j++ {
				select {
				case <-done:
					return
				default:
				}

				if n == int(n/3) {
					select {
					case toStop <- struct{}{}:
					default:
					}
					return
				}
			}
		}()
	}

	wg.Wait()
	toStop <- struct{}{}
	<-doneDone
}

func CheckWithMonitor(n int) {
	var wg sync.WaitGroup

	done := make(chan struct{})

	toStop := make(chan struct{}, 1)
	doneDone := make(chan struct{})
	go func() {
		<-toStop
		close(done)
		doneDone <- struct{}{}
	}()

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)

		go func() {
			stopMonitor := make(chan struct{})
			doneMonitor := make(chan struct{})
			defer func() {
				close(stopMonitor)
				<-doneMonitor
				wg.Done()
			}()

			var escape bool

			// monitor
			go func() {
				for {
					select {
					case <-done:
						escape = true
						doneMonitor <- struct{}{}
						return
					case <-stopMonitor:
						doneMonitor <- struct{}{}
						return
					default:
					}
				}
			}()

			select {
			case <-done:
				return
			default:
			}

			for j := 0; j < n; j++ {
				if escape { // faster than select-case
					return
				}

				if n == int(n/3) {
					select {
					case toStop <- struct{}{}:
					default:
					}
					return
				}
			}
		}()
	}

	wg.Wait()
	toStop <- struct{}{}
	<-doneDone
}

func BenchmarkNoCancellation1K(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NoCancellation(1000)
	}
}

// func BenchmarkCheckWithSelectChannel01K(b *testing.B) {
// 	for n := 0; n < b.N; n++ {
// 		CheckWithSelectChannel(1000)
// 	}
// }

func BenchmarkCheckWithSelectChannel1K(b *testing.B) {
	for n := 0; n < b.N; n++ {
		CheckWithSelectChannel(1000)
	}
}

func BenchmarkCheckWithMonitor1K(b *testing.B) {
	for n := 0; n < b.N; n++ {
		CheckWithMonitor(1000)
	}
}

func BenchmarkNoCancellation1M(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NoCancellation(1000000)
	}
}

// func BenchmarkCheckWithSelectChannel01M(b *testing.B) {
// 	for n := 0; n < b.N; n++ {
// 		CheckWithSelectChannel(1000000)
// 	}
// }

func BenchmarkCheckWithSelectChannel1M(b *testing.B) {
	for n := 0; n < b.N; n++ {
		CheckWithSelectChannel(1000000)
	}
}

func BenchmarkCheckWithMonitor1M(b *testing.B) {
	for n := 0; n < b.N; n++ {
		CheckWithMonitor(1000000)
	}
}
