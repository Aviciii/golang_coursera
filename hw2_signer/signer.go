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
	in = make(chan interface{}, 1)
	out = make(chan interface{}, 1)
	reader := bufio.NewReader(os.Stdin)
	arg, _ := reader.ReadString('\n')

	fmt.Println(arg)
	start := time.Now()
	go SingleHash(in, out)
	//go MultiHash(in, out)
	in <- "0"
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
	data, ok := (<-in).(string)

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

	out <- fmt.Sprintf("%v~%v", <-ch1, <-ch1)
	//out <- DataSignerCrc32(data) + "~" +
}

var MultiHash = func(in, out chan interface{}) {
	//var result, tmp string
	dataRaw := <-in
	data, ok := dataRaw.(string)
	if !ok {
		fmt.Println("cant convert result data to string")
	}
	fmt.Println(data)
	for i := 0; i <= 5; i++ {
		go func(out chan <- interface{}) {
			out <- DataSignerCrc32(strconv.Itoa(i) + data)
		}(out)
		fmt.Printf("%v MultiHash: crc32(th+step1)) %v %v\n", data, i, <-out)
		//result += fmt.Sprintf("%v", tmp)
	}
	close(out)
	//for s := range out {
	//
	//}
	//fmt.Printf("%v", result)
	//out <- result
}

var CombineResults = func(){}