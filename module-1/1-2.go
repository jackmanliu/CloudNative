package main

import (
	"context"
	"fmt"
	"time"
)

func producer(ch chan<- int) {
	ticker := time.NewTicker(time.Second)
	for count := range ticker.C {
		var seed = count.Second()
		select {
		case ch <- seed:
			fmt.Println("produce:", seed)
		default:
			fmt.Println("The channel is full")
		}
	}
}

func consumer(ch <-chan int) {
	ticker := time.NewTicker(time.Second)
	indicator := 10
	for _ = range ticker.C {
		if indicator == 0 {
			time.Sleep(9 * time.Second)
		}
		select {
		case n := <-ch:
			fmt.Println("      consume:", n)
		default:
			fmt.Println("The channel is empty now...")
		}
		indicator--
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*25)
	defer cancel()
	goods := make(chan int, 10)
	go producer(goods)
	go consumer(goods)
	select {
	case <-ctx.Done():
		fmt.Println("main process exit!")
	}
}
