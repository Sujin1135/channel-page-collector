package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
	"net/url"
	"strings"
	"sync"
	"time"
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
	collector   *colly.Collector
	accessMutex *sync.Mutex
	scrollMutex *sync.Mutex
}

const (
	baseUrl          = "https://www.youtube.com/results"
	htmlPrefix       = "var ytInitialData = "
	scrollDownScript = "window.scrollBy(0, 3000);"
	userAgent        = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
)

func NewWebsiteCollector() *Collector {
	return &Collector{
		collector: colly.NewCollector(
			colly.UserAgent(userAgent),
			colly.AllowedDomains("www.youtube.com", "youtube.com"),
		),
		accessMutex: &sync.Mutex{},
		scrollMutex: &sync.Mutex{},
	}
}

func (c *Collector) CollectWithScrolling(keyword string, numScrolls int) ([]string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	accessErr := c.accessWebsite(ctx, keyword)
	if accessErr != nil {
		return nil, errors.Wrap(accessErr, fmt.Sprintf("failed to access website by keyword as %s", keyword))
	}

	scrollDownErr := c.scrollDownNTimes(numScrolls, ctx)
	if scrollDownErr != nil {
		return nil, errors.Wrap(scrollDownErr, "failed to scroll down")
	}

	subscriberNames, extractErr := c.extractSubscriberNames(ctx)
	if extractErr != nil {
		return nil, errors.Wrap(extractErr, "failed to extract subscriber names")
	}

	fmt.Printf("총 %d개의 '#subscribers' 노드 수집됨\n", len(subscriberNames))

	return subscriberNames, nil
}

func (c *Collector) accessWebsite(ctx context.Context, keyword string) error {
	c.accessMutex.Lock()
	defer c.accessMutex.Unlock()

	fmt.Printf("start to access website by keyword: %s\n", keyword)
	defer fmt.Printf("end to access website by keyword: %s\n", keyword)

	uri := fmt.Sprintf("%s?search_query=%s&sp=%s", baseUrl, url.QueryEscape(keyword), "EgIQAg%253D%253D")
	err := chromedp.Run(ctx,
		chromedp.Navigate(uri),
		chromedp.Evaluate(scrollDownScript, nil),
	)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to navigate to %s", uri))
	}
	time.Sleep(1 * time.Second)
	return nil
}

func (c *Collector) extractSubscriberNames(ctx context.Context) ([]string, error) {
	var allNodes []*cdp.Node
	runErr := chromedp.Run(ctx,
		chromedp.Nodes(`#subscribers`, &allNodes, chromedp.NodeVisible),
	)
	if runErr != nil {
		return nil, errors.Wrap(runErr, "failed to fetch subscriber nodes")
	}

	subscriberNames := make([]string, 0, len(allNodes))

	for _, node := range allNodes {
		for _, children := range node.Children {
			subscriberNames = append(subscriberNames, children.NodeValue)
		}
	}
	return subscriberNames, nil
}

func (c *Collector) scrollDownNTimes(numScrolls int, ctx context.Context) error {
	for i := 0; i < numScrolls; i++ {
		err := c.scrollDown(ctx)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to scroll down by n %d times", numScrolls))
		}
	}
	return nil
}

func (c *Collector) scrollDown(ctx context.Context) error {
	c.scrollMutex.Lock()
	defer c.scrollMutex.Unlock()

	runErr := chromedp.Run(ctx,
		chromedp.Evaluate(scrollDownScript, nil),
	)
	if runErr != nil {
		return errors.Wrap(runErr, "failed to scroll page")
	}

	time.Sleep(1 * time.Second)
	return nil
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
