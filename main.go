package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gotest/pubsub"
)

func main() {
	p := pubsub.NewPublisher()
	s := pubsub.NewSubscriber("s1", nil)
	s2 := pubsub.NewSubscriber("s2", "test")
	p.Regiser(s)
	p.Regiser(s2)

	for {
		msg := ""
		fmt.Scanf("%s", &msg)

		p.Publish(msg, "")
	}
}

func PCModel() {
	msg := make(chan int, 10)

	go Producer(3, msg)
	go Producer(5, msg)
	go Consumer(msg)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("process exit: %v", <-sig)
}

func Producer(num int, ch chan int) {
	base := 1
	for {
		base += num
		ch <- base
		time.Sleep(1e8)
	}
}

func Consumer(ch chan int) {
	for {
		fmt.Println(<-ch)
	}
}
