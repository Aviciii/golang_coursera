package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// сюда писать код

var in chan interface{}
var out chan interface{}
var wg sync.WaitGroup

func main() {
	start := time.Now()
	//inputData := []int{0, 1, 1, 2, 3, 5, 8}
	in = make(chan interface{}, 1)
	out = make(chan interface{}, 1)
	reader := bufio.NewReader(os.Stdin)
	arg, _ := reader.ReadString('\n')

	fmt.Println(arg)
	wg.Add(2)
	go SingleHash(in, out)
	go MultiHash(in, out)
	out <- "0"
	wg.Wait()


	//fmt.Println(<-in)
	//fmt.Println(<-out)
	end := time.Since(start)
	//go MultiHash(ch2, ch1)
	//fmt.Printf("%v\n", <-ch2)
	//fmt.Printf("%v", <-ch1)

	log.Printf("run time: %v", end)
	fmt.Scanln()

	//result := <-ch2
	//fmt.Printf("%v", result)
	//MultiHash(data)
}

func ExecutePipeLine(pipeline []job) {
	for _, f := range pipeline {
		go f(in, out)
	}
}

var SingleHash = func(in, out chan interface{}) {
	defer wg.Done()
	data, ok := (<-out).(string)

	if !ok {
		fmt.Println("cant convert result data to string")
	}

	ch1 := make(chan string, 1)

	go func(ch1 chan string){
		ch1 <- DataSignerCrc32(DataSignerMd5(data))
	}(ch1)

	go func(ch1 chan string){
		ch1 <- DataSignerCrc32(data)
		runtime.Gosched()
	}(ch1)

	val := fmt.Sprintf("%v~%v", <-ch1, <-ch1)
	in <- val
	close(ch1)
}

var MultiHash = func(in, out chan interface{}) {
	defer wg.Done()
	data, ok := (<-in).(string)
	fmt.Println(data)
	hashCount := 6
	var chans [6]chan string

	if !ok {
		fmt.Println("cant convert result data to string")
	}

	for i := 0; i < hashCount; i++ {
		go func(th int) {
			mOut <- DataSignerCrc32(strconv.Itoa(th) + data)
			//runtime.Gosched()
		}(i)
	}

	for hashCount > 0 {
		fmt.Println(<-mOut)
		hashCount--
	}

	out <- "done"
}

var CombineResults = func(){}