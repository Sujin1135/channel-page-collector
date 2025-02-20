package main

import (
	"channel-page-collector/internal/collector"
)

func main() {
	websiteCollector := collector.NewWebsiteCollector()
	websiteCollector.Collect("운동")
}
