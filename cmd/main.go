package main

import (
	"channel-page-collector/internal/collector"
)

func main() {
	websiteCollector := collector.NewWebsiteCollector()
	keyword := "운동"
	websiteCollector.CollectWithScrolling(keyword, 5)
}
