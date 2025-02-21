package main

import (
	"channel-page-collector/internal/collector"
	"fmt"
	"log"
)

func main() {
	websiteCollector := collector.NewWebsiteCollector()
	channelIDs, err := websiteCollector.Collect("운동")
	if err != nil {
		log.Fatal(err)
	}
	for _, channelID := range channelIDs {
		fmt.Println(channelID)
	}
}
