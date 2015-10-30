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
