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
	wg.Add(2)

	start := time.Now()

	go func() {
		defer wg.Done()
		fmt.Println("start to collect method 1")
		subscriberNames, err := websiteCollector.CollectWithScrolling("운동", 5)
		if err != nil {
			log.Fatal("failed to collect subscriber names", err)
		}

		fmt.Println("ended to collect method 1")
		fmt.Println("subscriber names as follow:")
		for _, name := range subscriberNames {
			fmt.Println(name)
		}
	}()
	go func() {
		defer wg.Done()
		fmt.Println("start to collect method 2")
		subscriberNames, err := websiteCollector.CollectWithScrolling("뷰티", 5)
		if err != nil {
			log.Fatal("failed to collect subscriber names", err)
		}

		fmt.Println("ended to collect method 2")
		fmt.Println("subscriber names as follow:")
		for _, name := range subscriberNames {
			fmt.Println(name)
		}
	}()
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Println("elapsed:", elapsed)
}
