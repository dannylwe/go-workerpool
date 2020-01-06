package main

import (
	"time"
)

func main() {
	// program waits prints "hello" and waits three seconds for world
	dearChannel := make(chan string)
	worldChannel := make(chan string)

	go func() {
		time.Sleep(time.Second * 3)
		worldChannel <- "world"
	}()

	go func() {
		time.Sleep(time.Second * 2)
		dearChannel <- "dear"
	}()

	println("hello", <-dearChannel, <-worldChannel)
}
