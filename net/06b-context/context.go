package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func dealLongWithCtx(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	randTime := time.Second * 2
	timer := time.NewTimer(randTime)
	fmt.Printf("wait for %s \n", randTime)

	select {
	case <-timer.C:
		fmt.Println("Done")
	case <-ctx.Done():
		fmt.Println("Canceled")
	}
	fmt.Printf("I'm finished")
}

func main() {
	wg := &sync.WaitGroup{}
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Second)
	// ctx, cancel := context.WithTimeout(ctx, time.Second)

	wg.Add(1)
	go dealLongWithCtx(wg, ctx)
	wg.Wait()
}
