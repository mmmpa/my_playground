package main

import (
	"github.com/mmmpa/my_playground/easy"
	"log"
	"time"
	"fmt"
)

func main() {
	urls := genUrls(30)
	maxWorkers := 3

	in := make(chan string)
	out := make(chan time.Duration)

	// worker へ url を送る措置
	go send(urls, in)

	// worker 発行
	for i := 0; i < maxWorkers; i++ {
		go load(in, out)
	}

	log.Printf("finish: total: %v\n", sum(urls, out))
}

func genUrls(n int) []string {
	urls := make([]string, 0)

	for i := 1; i <= n; i++ {
		urls = append(urls, fmt.Sprintf("url_%d", i))
	}
	return urls
}

func sum(urls []string, out chan time.Duration) time.Duration {
	total := time.Duration(0)

	for i := 0; i < len(urls); i++ {
		s, ok := <-out

		if !ok {
			log.Printf("out is closed")
			break
		}

		total += s
	}

	return total
}

func send(urls []string, in chan string) {
	for _, url := range urls {
		log.Printf("send: %v\n", url)
		in <- url
	}
	close(in)
}

func load(in chan string, out chan time.Duration) {
	for {
		url, ok := <-in

		if !ok {
			return
		}

		log.Printf("load: %v\n", url)
		s := easy.RandomSecsSleep(3)
		log.Printf("loaded: %v %v\n", url, s)

		out <- s
	}
}
