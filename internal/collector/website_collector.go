package collector

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

type ChannelResponse struct {
	Contents struct {
		TwoColumnSearchResultsRenderer struct {
			PrimaryContents struct {
				SectionListRenderer struct {
					SectionListRendererContents []struct {
						ItemSectionRenderer struct {
							ItemSectionRendererContents []struct {
								ChannelRenderer struct {
									ChannelID string `json:"channelId"`
								} `json:"channelRenderer"`
							} `json:"contents"`
						} `json:"itemSectionRenderer"`
					} `json:"contents"`
				} `json:"sectionListRenderer"`
			} `json:"primaryContents"`
		} `json:"twoColumnSearchResultsRenderer"`
	} `json:"contents"`
}

type Collector struct {
	collector *colly.Collector
}

const (
	baseUrl    = "https://www.youtube.com/results"
	htmlPrefix = "var ytInitialData = "
)

func NewWebsiteCollector() *Collector {
	return &Collector{collector: colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
		colly.AllowedDomains("www.youtube.com", "youtube.com"),
	)}
}

func (c *Collector) Collect(keyword string) ([]string, error) {
	var channelIDs []string
	var collectErr error

	c.collector.OnHTML("script", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, htmlPrefix) {
			response, err := c.extractChannelResponse(e.Text)
			if err != nil {
				collectErr = errors.Wrap(err, "extract channel response error:")
				return
			}

			for _, content := range response.Contents.TwoColumnSearchResultsRenderer.PrimaryContents.SectionListRenderer.SectionListRendererContents {
				for _, inner := range content.ItemSectionRenderer.ItemSectionRendererContents {
					channelIDs = append(channelIDs, inner.ChannelRenderer.ChannelID)
				}
			}
		}
	})

	c.collector.Visit(fmt.Sprintf("%s?search_query=%s&sp=%s", baseUrl, url.QueryEscape(keyword), "EgIQAg%253D%253D"))
	c.collector.Wait()

	if collectErr != nil {
		return nil, collectErr
	}

	return channelIDs, nil
}

func (c *Collector) extractChannelResponse(content string) (*ChannelResponse, error) {
	startIdx := strings.Index(content, htmlPrefix)
	jsonPart := content[startIdx+len(htmlPrefix):]
	jsonPart = strings.TrimSpace(jsonPart)

	if strings.HasSuffix(jsonPart, ";") {
		jsonPart = strings.TrimSuffix(jsonPart, ";")
	}

	var response ChannelResponse
	if err := json.Unmarshal([]byte(jsonPart), &response); err != nil {
		return nil, errors.New("occurred an error when extract channel response")
	}

	return &response, nil
}
