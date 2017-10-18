package main

import (
	"github.com/mmmpa/my_playground/easy"
	"log"
	"time"
)

func main() {
	ch := make(chan time.Duration)

	go async(ch)
	go async(ch)
	go async(ch)

	total := <-ch + <-ch + <-ch

	log.Printf("finish: total: %v\n", total)
}

func async(ch chan time.Duration) {
	s := easy.RandomSecsSleep(3)
	log.Printf("sleep: %v\n", s)

	ch <- s
}
