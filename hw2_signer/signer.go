package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// сюда писать код
var SingleHash = func(in, out chan interface{}) {
	var data string
	var i uint8
	var crc32md5Chans []chan interface{}
	var crc32Chans []chan interface{}

	for o := range in {
		data = fmt.Sprintf("%v", o)

		md5 := DataSignerMd5(data)
		crc32md5Chans = append(crc32md5Chans, make(chan interface{}))
		crc32Chans = append(crc32Chans, make(chan interface{}))

		go func(ch chan interface{}, md5 string) {
			ch <- DataSignerCrc32(md5)
			close(ch)
		}(crc32md5Chans[i], md5)

		go func(ch chan interface{}, data string) {
			ch <- DataSignerCrc32(data)
			close(ch)
		}(crc32Chans[i], data)

		i++
	}

	for i > 0 {
		out <- fmt.Sprintf("%s~%s", <-crc32Chans[i-1], <-crc32md5Chans[i-1])
		i--
	}
}

var MultiHash = func(in, out chan interface{}) {
	i, j, thMax := 0, 0, 6
	var RrChans [][]chan string
	var result string

	for shData := range in {
		data := fmt.Sprintf("%v", shData)
		RrChans = append(RrChans, make([]chan string, thMax))
		for i = 0; i < thMax; i++ {
			RrChans[j][i] = make(chan string)
			go func(result chan string, i int, data string) {
				th := strconv.Itoa(i)
				thData := DataSignerCrc32(th + data)
				result <- thData
				close(result)
			}(RrChans[j][i], i, data)
		}

		j++
	}

	for _, rChan := range RrChans {
		result = ""
		for _, r := range rChan {
			result += <-r
		}
		out <- result
	}
}

var CombineResults = func(in, out chan interface{}) {
	var result []string

	for mData := range in {
		data := fmt.Sprintf("%v", mData)
		result = append(result, data)
	}

	sort.Strings(result)
	strResult := strings.Join(result, "_")

	out <- strResult
}

var ExecutePipeline = func(pipeline ...job) {
	chanOutSlice := make([]chan interface{}, len(pipeline))
	chanInSlice := make([]chan interface{}, len(pipeline)+1)
	wg := &sync.WaitGroup{}

	for i, f := range pipeline {
		chanOutSlice[i] = make(chan interface{})
		chanInSlice[i] = make(chan interface{})

		f := f
		i := i
		go func() {
			f(chanInSlice[i], chanOutSlice[i])
			close(chanOutSlice[i])
		}()
	}

	for n, ch := range chanOutSlice {
		wg.Add(1)
		go func(ch chan interface{}, chanInSlice []chan interface{}, n int, wg *sync.WaitGroup) {
			defer wg.Done()
			for val := range ch {
				chanInSlice[n+1] <- val
			}
			if chanInSlice[n+1] != nil {
				close(chanInSlice[n+1])
			}
		}(ch, chanInSlice, n, wg)
	}
	wg.Wait()
}
