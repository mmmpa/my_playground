package easy

import (
	"time"
	"math/rand"
)

func RandomSecs(i int) time.Duration {
	rand.Seed(time.Now().UnixNano())
	return time.Duration(rand.Intn(i)) * time.Second
}

func RandomSecsSleep(i int) time.Duration {
	sleep := RandomSecs(i) + time.Second
	time.Sleep(sleep)

	return sleep
}
