package main

import (
	"github.com/mmmpa/my_playground/easy"
	"log"
	"time"
	"sync"
)

func main() {
	// url を生成する自前関数です
	urls := easy.GenUrls(5)

	start := time.Now().UnixNano()
	maxWorkers := 2

	worker_in := provide(urls)
	worker_out := fetch(maxWorkers, worker_in)
	task_result := receive(worker_out)

	log.Printf(
		"finish: worker_total: %v, total: %v \n",
		<-task_result,
		time.Duration(time.Now().UnixNano()-start),
	)
}

func receive(out chan time.Duration) chan time.Duration {
	re := make(chan time.Duration)

	go func() {
		defer close(re)

		total := time.Duration(0)
		defer func() { re <- total }()

		for s := range out {
			total += s
		}
	}()

	return re
}

func provide(urls []string) chan string {
	in := make(chan string)

	go func() {
		defer close(in)

		for _, url := range urls {
			log.Printf("send: %v\n", url)
			in <- url
		}
	}()

	return in
}

func fetch(maxWorkers int, in chan string) chan time.Duration {
	out := make(chan time.Duration)

	go func() {
		defer close(out)

		wg := sync.WaitGroup{}
		for i := 0; i < maxWorkers; i++ {
			wg.Add(1)

			// worker 本体
			go func() {
				defer wg.Done()
				for url := range in {
					log.Printf("start fetching: %v\n", url)
					s := easy.RandomSecsSleep(3)
					log.Printf("fetched: %v %v\n", url, s)

					out <- s
				}
				log.Print("in closed\n")
			}()
		}
		wg.Wait()
	}()

	return out
}
