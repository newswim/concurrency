package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	done := make(chan struct{})

	pub := &publisher{done}
	rec := &receiver{
		ch: make(chan int),
	}

	go pub.run(rec.ch)
	go rec.run()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c

	done <- struct{}{}
	close(rec.ch)
	wg.Wait()
}

type publisher struct {
	done chan struct{}
}

func (p *publisher) run(rec chan int) {
	t := time.NewTicker(time.Second)

	i := 0

	for {
		select {
		case <-t.C:
			rec <- i
			i++
		case <-p.done:
			return
		}
	}
}

type receiver struct {
	ch chan int
}

func (r *receiver) run() {
	for i := range r.ch {
		wg.Add(1)
		fmt.Printf("Hello %v\n", i)
		time.Sleep(2 * time.Second)
		fmt.Printf("Hello Again %v\n\n", i)
		wg.Done()
	}
}
