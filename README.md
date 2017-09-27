# golang concurrency patterns

## Producer-consumer-collector


## Cancellation

```
goos: linux
goarch: amd64
pkg: github.com/shenwei356/golang-concurrency-patterns
BenchmarkNoCancellation1K-4               100000             13109 ns/op              16 B/op          1 allocs/op
BenchmarkCheckWithSelectChannel1K-4        30000             50678 ns/op             304 B/op          4 allocs/op
BenchmarkCheckWithMonitor1K-4             100000             23215 ns/op            1076 B/op         16 allocs/op
BenchmarkNoCancellation1M-4                  500           2750034 ns/op              16 B/op          1 allocs/op
BenchmarkCheckWithSelectChannel1M-4          100          15489212 ns/op             305 B/op          4 allocs/op
BenchmarkCheckWithMonitor1M-4                500           2867392 ns/op            1086 B/op         16 allocs/op

```

Results show that using `select-case` in loop for cancellation is slow,
while it's much faster by checking in a goroutine.
