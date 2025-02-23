package main

import (
	"channel-page-collector/internal/collector"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	websiteCollector := collector.NewWebsiteCollector()

	var wg sync.WaitGroup
	wg.Add(1)

	start := time.Now()
	subscriberNamesChan := make(chan []string, 1)

	go func() {
		defer wg.Done()
		fmt.Println("start to collect method 1")
		err := websiteCollector.CollectWithScrolling("운동", 5, subscriberNamesChan)
		if err != nil {
			log.Fatal("failed to collect subscriber names", err)
		}
	}()

	fmt.Println("channel data as follow:")
	for v := range subscriberNamesChan {
		fmt.Println(v)
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Println("elapsed:", elapsed)
}
