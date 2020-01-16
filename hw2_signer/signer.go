package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

// сюда писать код

var in chan interface{}
var out chan interface{}

func main() {
	//inputData := []int{0, 1, 1, 2, 3, 5, 8}
	in = make(chan interface{}, 1)
	out = make(chan interface{}, 1)
	reader := bufio.NewReader(os.Stdin)
	arg, _ := reader.ReadString('\n')

	fmt.Println(arg)
	start := time.Now()
	go SingleHash(in, out)
	go MultiHash(in, out)
	out <- "0"
	//for r := range out {
	//	log.Println(r)
	//}

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
	data, ok := (<-in).(string)
	fmt.Println(data)

	if !ok {
		fmt.Println("cant convert result data to string")
	}

	result := ""

	for i := 0; i <= 5; i++ {
		go func(th int) {
			fmt.Println(th)
			out <- DataSignerCrc32(strconv.Itoa(th) + data)
			//runtime.Gosched()
		}(i)
	}

	for s := range out {
		fmt.Println(s)
	}

	out <- result
}

var CombineResults = func(){}