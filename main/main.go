package main

import (
	"./pool"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var sum int32

func myFunc(i interface{}) error {
	n := i.(int32)
	atomic.AddInt32(&sum, n)
	fmt.Printf("run with %d\n", n)
	return nil
}

func demoFunc() error {
	for i := 0; i < 10000; i++ {

	}
	return nil
}

func main() {
	p, _ := pool.NewProcessPool(1000)
	defer p.PoolClose()
	start := time.Now()
	runTimes := 100000

	// use the common pool
	var wg sync.WaitGroup
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		var args []interface{}
		args = append(args, int32(i))
		p.ProcessSubmit(func(args []interface{}) {
			myFunc(args[0])
			wg.Done()
		}, args)
	}
	/*for i := 0; i < runTimes; i++ {
		wg.Add(1)
		var args []interface{}
		args = append(args, int32(i))
		go func(args []interface{}){
		myFunc(args[0])
		wg.Done()
		}(args)
	}*/
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", p.Running())
	fmt.Printf("total goroutines: %d\n", p.Size())
	fmt.Printf("finish all tasks.\n")
	fmt.Printf("time: %v\n", time.Now().Sub(start))
	// use the pool with a function
	// set 10 the size of goroutine pool and 1 second for expired duration
}
