package main

import (
	"github.com/mmmpa/my_playground/easy"
	"log"
	"time"
	"fmt"
	"context"
	"sync"
)

type Result struct {
	URL      string
	Duration time.Duration
}

func main() {
	start := time.Now().UnixNano()
	urls := genUrls(100)
	maxWorkers := 5

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	worker_in := make(chan string)
	worker_out := make(chan Result)
	task_result := make(chan []Result)

	// worker へ task を送る措置  (今回の場合 url を供給する)
	go send(ctx, urls, worker_in)

	// worker から結果を収集する (終了条件が必要なので全 task である urls を渡す)
	go pick(ctx, urls, worker_out, task_result)

	// worker 発行
	go load(ctx, maxWorkers, worker_in, worker_out)

	results := <-task_result

	log.Printf(
		"finish: result: %+v, total_loaded: %v, total: %v \n",
		results,
		len(results),
		time.Duration(time.Now().UnixNano()-start),
	)

	time.Sleep(time.Second * 3)
}

func genUrls(n int) []string {
	urls := make([]string, 0)

	for i := 1; i <= n; i++ {
		urls = append(urls, fmt.Sprintf("url_%d", i))
	}
	return urls
}

func pick(ctx context.Context, urls []string, out chan Result, re chan []Result) {
	results := make([]Result, 0)
	defer close(re)
	defer func() { re <- results }()

	for i := 0; i < len(urls); i++ {
		select {
		case <-ctx.Done():
			fmt.Println("cancel:", ctx.Err())
			return
		case result, ok := <-out:
			if !ok {
				log.Printf("out is closed")
				return
			}

			results = append(results, result)
		}
	}
}

func send(ctx context.Context, urls []string, in chan string) {
	defer close(in)

	for i, url := range urls {
		select {
		case <-ctx.Done():
			log.Printf("cancel: %v, rest: %v", ctx.Err(), len(urls)-i)
			return
		case in <- url:
			log.Printf("sent: %v\n", url)
		}
	}
}

func load(ctx context.Context, maxWorkers int, in chan string, out chan Result) {
	defer close(out)

	wg := sync.WaitGroup{}

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer log.Print("in closed\n")

			for url := range in {
				log.Printf("load: %v\n", url)
				s := easy.RandomSecsSleep(3)

			L:
				for {
					select {
					case <-ctx.Done():
						log.Printf("loaded: but canceled %v\n", url)
						return
					case out <- Result{
						URL:      url,
						Duration: s,
					}:
						log.Printf("loaded: %v %v\n", url, s)
						break L
					}
				}
			}
		}()
	}

	wg.Wait()
}
