package main

import (
	"github.com/mmmpa/my_playground/easy"
	"log"
)

func main() {
	go async()
	go async()
	go async()

	log.Print("finish\n")
}

func async() {
	s := easy.RandomSecsSleep(3)
	log.Printf("sleep: %v\n", s)
}
