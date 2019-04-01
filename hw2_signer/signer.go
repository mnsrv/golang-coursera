package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func startWorker(job job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)
	job(in, out)
}

// ExecutePipeline обеспечивает конвейерную обработку функций-воркеров
func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})
	out := make(chan interface{})

	for _, job := range jobs {
		wg.Add(1)
		go startWorker(job, in, out, wg)
		in = out
		out = make(chan interface{})
	}

	wg.Wait()
}

// SingleHash считает значение crc32(data)+"~"+crc32(md5(data))
// ( конкатенация двух строк через ~),
// где data - то что пришло на вход (по сути - числа из первой функции)
func SingleHash(in, out chan interface{}) {
	for i := range in {
		value, ok := i.(int)
		if !ok {
			panic("SingleHash: failed type assertion")
		}
		data := strconv.Itoa(value)
		result := DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data))
		out <- result
	}
}

// MultiHash считает значение crc32(th+data))
// (конкатенация цифры, приведённой к строке и строки),
// где th=0..5 ( т.е. 6 хешей на каждое входящее значение ),
// потом берёт конкатенацию результатов в порядке расчета (0..5),
// где data - то что пришло на вход (и ушло на выход из SingleHash)
func MultiHash(in, out chan interface{}) {
	for i := range in {
		data, ok := i.(string)
		if !ok {
			panic("MultiHash: failed type assertion")
		}
		var result string
		for th := 0; th < 6; th++ {
			result += DataSignerCrc32(strconv.Itoa(th) + data)
		}
		out <- result
	}
}

// CombineResults получает все результаты,
// сортирует (https://golang.org/pkg/sort/),
// объединяет отсортированный результат через _ (символ подчеркивания) в одну строку
func CombineResults(in, out chan interface{}) {
	var hashes []string
	for i := range in {
		data, ok := i.(string)
		if !ok {
			panic("CombineResults: failed type assertion")
		}
		hashes = append(hashes, data)
	}
	sort.Strings(hashes)

	result := strings.Join(hashes, "_")
	out <- result
}

func main() {}
