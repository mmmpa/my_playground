package main

import (
	"github.com/mmmpa/my_playground/easy"
	"log"
	"time"
	"fmt"
)

func main() {
	start := time.Now().UnixNano()
	urls := genUrls(100)
	maxWorkers := 5

	worker_in := make(chan string)
	worker_out := make(chan time.Duration)
	task_result := make(chan time.Duration)

	// worker へ task を送る措置  (今回の場合 url を供給する)
	go send(urls, worker_in)

	// worker から結果を収集する (終了条件が必要なので全 task である urls を渡す)
	go pick(urls, worker_out, task_result)

	// worker 発行
	for i := 0; i < maxWorkers; i++ {
		go load(worker_in, worker_out)
	}

	log.Printf(
		"finish: worker_total: %v, total: %v \n",
		<-task_result,
		time.Duration(time.Now().UnixNano()-start),
	)
}

func genUrls(n int) []string {
	urls := make([]string, 0)

	for i := 1; i <= n; i++ {
		urls = append(urls, fmt.Sprintf("url_%d", i))
	}
	return urls
}

func pick(urls []string, out chan time.Duration, re chan time.Duration) {
	total := time.Duration(0)

	for i := 0; i < len(urls); i++ {
		s, ok := <-out

		if !ok {
			log.Printf("out is closed")
			break
		}

		total += s
	}

	re <- total
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
