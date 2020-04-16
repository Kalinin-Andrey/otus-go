package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

func emulateLongOperation(ctx context.Context, id int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randVal := r.Intn(5000)
	randTime := time.Duration(randVal) * time.Millisecond
	log.Printf("Job %d Will be evalutated for %d", id, randVal)
	timer := time.NewTimer(randTime)

	select {
	case <-timer.C:
		log.Printf("Successfully finished job %d", id)
	case <-ctx.Done():
		log.Printf("id %d timed out", id)
	}
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			emulateLongOperation(ctx, id)
			//cancel()
		}(i)
	}

	wg.Wait()
}
