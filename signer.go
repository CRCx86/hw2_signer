package main

import (
	"sort"
	"strconv"
	"sync"
)

// сюда писать код

func SingleHash(in, out chan interface{}) {

	//mutex := &sync.Mutex{}

	for i := range in {

		var s string

		i1 := i.(int)
		//var result string
		s = strconv.Itoa(i1)

		ch1 := funcSingleSolt(s)
		//mutex.Lock()
		md5 := DataSignerMd5(s)
		//mutex.Unlock()
		ch2 := funcSingleSolt(md5)

		resCh := collectSingleHash(ch1, ch2)
		//result = <-resCh
		//fmt.Println(result)
		out <- resCh

	}

}

func collectSingleHash(ch1, ch2 chan string) chan string {
	ch := make(chan string)
	go func(ch1, ch2 chan string, ch chan string) {
		s := <-ch1 + "~" + <-ch2
		//fmt.Println(s)
		ch <- s
	}(ch1, ch2, ch)
	return ch
}

func funcSingleSolt(s string) chan string {
	ch := make(chan string)
	go func(s1 string, c chan string) {
		crc32 := DataSignerCrc32(s1)
		//fmt.Println(crc32)
		c <- crc32
	}(s, ch)
	return ch
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	multiHashIn := make(chan interface{}, 100)
	go multiHashCollectValues(wg, multiHashIn, out)
	for i := range in {
		s := <- i.(chan string)
		chans := multiHashEx(s)
		wg.Add(1)
		multiHashIn <- chans
	}
	wg.Wait()
	close(multiHashIn)
}

func multiHashCollectValues(wg *sync.WaitGroup, multiHashIn chan interface{}, out chan interface{}) {

	var result string
	for i := range multiHashIn {
		ch := i.([]<-chan string)
		result = ""
		for _, c := range ch {
			c2 := <-c
			result = result + c2
		}
		//fmt.Println(result)
		out <- result
		wg.Done()
	}
}

func multiHashEx(s string) []<-chan string {
	chans := make([]<-chan string, 0)
	for j := 0; j <= 5; j++ {
		chans = append(chans, funcGetSoltMultiData(j, s))
	}
	return chans
}

func funcGetSoltMultiData(j int, s string) chan string {

	out := make(chan string)
	go func(st string, out chan<- string) {
		crc32 := DataSignerCrc32(st)
		out <- crc32
	}(strconv.Itoa(j)+s, out)
	return out
}

func CombineResults(in, out chan interface{}) {

	var result string
	strings := make([]string, 0)
	for i := range in {
		var s string
		s = i.(string)
		strings = append(strings, s)
		//fmt.Println(strings)
	}

	sort.Strings(strings)

	for _, s := range strings {
		if len(result) > 0 {
			result += "_" + s
		} else {
			result = s
		}
	}
	//fmt.Println(result)

	out <- result
}

func ExecutePipeline(jobs ...job) {

	wg := &sync.WaitGroup{}
	in := make(chan interface{}, 100)
	for _, job := range jobs {
		wg.Add(1)
		out := make(chan interface{}, 100)
		//go ExecutePipelineFn(wg, in, out, job)
		ExecutePipelineFn(wg, in, out, job)
		in = out
	}
	wg.Wait()
}

//func ExecutePipelineFn(wg *sync.WaitGroup, in chan interface{}, out chan interface{}, j job) {
//	j(in, out)
//	close(out)
//	wg.Done()
//}

func ExecutePipelineFn(wg *sync.WaitGroup, in chan interface{}, out chan interface{}, j job) {
	go func(in, out chan interface{}, j job) {
		j(in, out)
		close(out)
		wg.Done()
	}(in, out, j)
}
