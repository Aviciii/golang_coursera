package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

// сюда писать код

func main() {
	start := time.Now()
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	go SingleHash(ch1, ch2)
	ch1 <- "0"
	fmt.Printf("%v", <-ch2)
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

func ExecutePipeLine() {

}

var SingleHash = func(in, out chan interface{}) {
	var result string
	ch1 := make(chan interface{}, 1)
	dataRaw := <-in
	data, ok := dataRaw.(string)

	if !ok {
		fmt.Println("cant convert result data to string")
	}

	go func(out chan <- interface{}) {
		out <- DataSignerCrc32(data)
	}(ch1)

	go func(out chan <- interface{}) {
		out <- DataSignerCrc32(DataSignerMd5(data))
		close(out)
	}(ch1)

	for i := range ch1 {
		result += fmt.Sprintf("%v~", i)
	}

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
			//out <- tmp
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