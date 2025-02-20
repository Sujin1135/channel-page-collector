package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"strings"
	"time"
)

func main() {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
		colly.AllowedDomains("www.youtube.com", "youtube.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		time.Sleep(1 * time.Second)
	})

	// a 태그의 href 속성에서 채널 ID 추출
	c.OnHTML(`script`, func(e *colly.HTMLElement) {
		content := e.Text
		prefix := "var ytInitialData = "

		if strings.Contains(content, prefix) {
			startIdx := strings.Index(content, prefix)
			jsonPart := content[startIdx+len(prefix):]
			jsonPart = strings.TrimSpace(jsonPart)

			if strings.HasSuffix(jsonPart, ";") {
				jsonPart = strings.TrimSuffix(jsonPart, ";")
			}

			var ytInitialData interface{}
			if err := json.Unmarshal([]byte(jsonPart), &ytInitialData); err != nil {
				log.Printf("JSON 파싱 오류: %v", err)
				return
			}
			dataMap, ok := ytInitialData.(map[string]interface{})
			if !ok {
				log.Println("ytInitialData의 타입이 예상과 다릅니다.")
				return
			}

			contentsMap := dataMap["contents"].(map[string]interface{})
			twoColumn := contentsMap["twoColumnSearchResultsRenderer"].(map[string]interface{})
			primaryContents := twoColumn["primaryContents"].(map[string]interface{})
			sectionListRenderer := primaryContents["sectionListRenderer"].(map[string]interface{})
			contents := sectionListRenderer["contents"].([]interface{})
			for _, content := range contents {
				if val, ok := content.(map[string]interface{})["itemSectionRenderer"]; ok {
					itemRenderer := val.(map[string]interface{})
					innerContents := itemRenderer["contents"].([]interface{})
					for _, inner := range innerContents {
						if channelRenderer, ok := inner.(map[string]interface{})["channelRenderer"].(map[string]interface{}); ok {
							if channelId, ok := channelRenderer["channelId"].(string); ok {
								fmt.Println(channelId)
							}
						}
					}
				}
			}
		}
	})

	c.Visit("https://www.youtube.com/results?search_query=%EB%B7%B0%ED%8B%B0&sp=EgIQAg%253D%253D")
	s := "gopher"
	fmt.Println("Hello and welcome, %s!", s)
}
