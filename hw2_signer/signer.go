package main

// сюда писать код
import (
	"fmt"
	"sort"
	"sync"
	"strconv"
	"strings"
)

var rateLim chan struct{} = make(chan struct{}, 1)

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for input := range in {
		data := strconv.Itoa(input.(int))
		fmt.Printf("%s SingleHash data %[1]s\n", data)
		wg.Add(1)
		go calcSH(out, data, wg)
	}
	wg.Wait()
}

func calcSH(out chan interface{}, data string, wg *sync.WaitGroup){
	defer wg.Done()
	ch0 := make(chan string)
	go func(ch0 chan<- string) {
		rateLim <- struct{}{}
		ch0 <- DataSignerMd5(data)
		<- rateLim
	}(ch0)
	md5Hash := <-ch0
	fmt.Printf("%s SingleHash md5(data) %s\n", data, md5Hash)
	ch1 := make(chan string)
	go func(ch chan<- string) {
		ch <- DataSignerCrc32(md5Hash)
	}(ch1)
	ch2 := make(chan string)
	go func(ch chan<- string) {
		ch <- DataSignerCrc32(data)
	}(ch2)
	crc32md5Hash := <-ch1
	fmt.Printf("%s SingleHash crc32(md5(data)) %s\n", data, crc32md5Hash)
	crc32Hash := <-ch2
	fmt.Printf("%s SingleHash crc32(data) %s\n", data, crc32Hash)
	result := crc32Hash + "~" + crc32md5Hash
	fmt.Printf("%s SingleHash result %s\n", data, result)
	out <- result
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for input := range in {
		data := input.(string)
		wg.Add(1)
		go calcMH(out, data, wg)
	}
	wg.Wait()
}

func calcMH(out chan interface{}, data string, wg0 *sync.WaitGroup) {
	defer wg0.Done()
	wg := &sync.WaitGroup{}
	hashes := make([]string, 6)
	wg.Add(6)
	for i := 0; i < 6; i++ {
		go func(i int, h []string) {
			md5h := DataSignerCrc32(strconv.Itoa(i) + data)
			h[i] = md5h
			fmt.Printf("%s MultiHash: crc32(th+step1) %d %s\n", data, i, md5h)
			wg.Done()
		}(i, hashes)
	}
	wg.Wait()
	result := strings.Join(hashes, "")
	fmt.Printf("%s MultiHash: result %s\n", data, result)
	out <- result
}

func CombineResults(in, out chan interface{}) {
	var results []string
	for hash := range in {
		results = append(results, hash.(string))
	}
	sort.Strings(results)
	combined := strings.Join(results, "_")
	fmt.Printf("CombineResults\n%s", combined)
	out <- combined
}

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	var in chan interface{}
	out := make(chan interface{})
	for _, j := range jobs {
		runJob(j, in, out, wg)
		in = out
		out = make(chan interface{})
	}
	wg.Wait()
}

func runJob(f job, in, out chan interface{}, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		f(in, out)
		close(out)
		wg.Done()
	}()
}
