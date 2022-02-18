package main

// сюда писать код
import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func SingleHash(in, out chan interface{}) {
	data := strconv.Itoa((<-in).(int))
	fmt.Printf("%s SingleHash data %[1]s\n", data)
	md5Hash := DataSignerMd5(data)
	fmt.Printf("%s SingleHash md5(data) %s\n", data, md5Hash)
	crc32md5Hash := DataSignerCrc32(md5Hash)
	fmt.Printf("%s SingleHash crc32(md5(data)) %s\n", data, crc32md5Hash)
	crc32Hash := DataSignerCrc32(data)
	fmt.Printf("%s SingleHash crc32(data) %s\n", data, crc32Hash)
	result := crc32Hash + "~" + crc32md5Hash
	fmt.Printf("%s SingleHash result %s\n", data, result)
	out <- result
}

func MultiHash(in, out chan interface{}) {
	data := (<-in).(string)
	result := ""
	for i := 0; i < 6; i++ {
		tempHash := DataSignerCrc32(strconv.Itoa(i) + data)
		result += tempHash
		fmt.Printf("%s MultiHash: crc32(th+step1) %d %s\n", data, i, tempHash)
	}
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
	in := make(chan interface{})
	// var in chan interface{}
	out := make(chan interface{})
	for _, j := range jobs {
		go j(in, out)
		in = out
		out = make(chan interface{})
	}
	fmt.Scanln()
}
