package main

import (
	"github.com/mmmpa/my_playground/easy"
	"log"
	"time"
	"fmt"
	"sync"
)

func main() {
	urls := genUrls(3)
	maxWorkers := 3

	var waiter sync.WaitGroup
	in := make(chan string)
	out := make(chan time.Duration)

	// worker 発行
	for i := 0; i < maxWorkers; i++ {
		waiter.Add(1)
		go async(&waiter, in, out)
	}

	// worker へ url を送る措置
	go func() {
		defer close(in)
		for _, url := range urls {
			// すべての送信が受け付けられるまで続く (受信されるまでブロックされる)
			log.Printf("send: %v\n", url)
			in <- url
		}
	}()

	total := time.Duration(0)
	for i := 0; i < len(urls); i++ {
		s, ok := <-out

		// timeout 措置を取る場合は out を close する
		if !ok {
			log.Printf("out is closed")
			break
		}

		total += s
	}

	// すべての worker が終了した判定措置
	// これいらなさそう
	waiter.Wait()

	log.Printf("finish: total: %v\n", total)
}

func genUrls(n int) []string {
	urls := make([]string, 0)

	for i := 1; i <= n; i++ {
		urls = append(urls, fmt.Sprintf("url_%d", i))
	}
	return urls
}

func async(waiter *sync.WaitGroup, in chan string, out chan time.Duration) {
	defer waiter.Done()
	for {
		url, ok := <-in

		if !ok {
			return
		}

		s := easy.RandomSecsSleep(3)
		log.Printf("sleep: %v on %v\n", s, url)

		out <- s
	}
}
