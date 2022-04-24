package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var flag = false

func producer(id int, wg *sync.WaitGroup, ch chan<- string) {
	count := 0
	for !flag {
		time.Sleep(1 * time.Second)
		count++
		data := strconv.Itoa(id) + ":" + strconv.Itoa(count)
		fmt.Printf("producer %d produce:", id, data)
		ch <- data
	}
	wg.Done()
}

func consumer(wg *sync.WaitGroup, ch <-chan string) {
	for data := range ch {
		time.Sleep(1 * time.Second)
		fmt.Println("consume:", data)
	}
	wg.Done()
}

func main() {
	ch := make(chan string, 10)
	wgPd := new(sync.WaitGroup)
	wgCs := new(sync.WaitGroup)

	for i := 0; i < 3; i++ {
		wgPd.Add(1)
		go producer(i, wgPd, ch)
	}

	for j := 0; j < 2; j++ {
		wgCs.Add(1)
		go consumer(wgCs, ch)
	}

	go func() {
		time.Sleep(10 * time.Second)
		flag = true
	}()

	wgPd.Wait()
	close(ch)
	wgCs.Wait()
}
