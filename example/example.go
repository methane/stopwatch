package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/methane/stopwatch"
)

func main() {
	wg := sync.WaitGroup{}

	myfunc := func() {
		defer wg.Done()
		defer stopwatch.Start("aaa").Stop()
		time.Sleep(time.Millisecond * 10)
		defer stopwatch.Start("bbb").Stop()
		time.Sleep(time.Millisecond * 10)
	}

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go myfunc()
	}
	wg.Wait()
	fmt.Println(stopwatch.Show())
}

// Output:
// name	count	avg	total
// aaa /Users/inada-n/go1.5/src/github.com/methane/stopwatch/example/example.go:16	100	25.500379ms	2.550037924s
// bbb /Users/inada-n/go1.5/src/github.com/methane/stopwatch/example/example.go:18	100	12.799502ms	1.279950208s
