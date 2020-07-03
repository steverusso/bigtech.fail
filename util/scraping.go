// This is free and unencumbered software released
// into the public domain. Please see the UNLICENSE
// file or unlicense.org for more information.
package util

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"golang.org/x/sync/errgroup"
)

type Ctx = context.Context

var (
	errEmptyURL = errors.New("empty URL")
	errEmptySel = errors.New("empty selector")
)

type scraperFunc func(Ctx, *policy) error

var nonDefaultScrapers = map[string]scraperFunc{
	"twitter/commguides": twtrCommGuides,
	"youtube/commguides": ytCommGuides,
}

func ScrapePlatform(ctx Ctx, p *Platform) error {
	if ctx == nil {
		var cancel func()
		ctx, cancel = chromedp.NewContext(context.Background())
		defer cancel()
	}

	for k, pol := range p.Policies {
		// Use a non default scraper if one is set for this policy.
		if fn, ok := nonDefaultScrapers[p.Key()+"/"+k]; ok {
			if err := fn(ctx, pol); err != nil {
				return fmt.Errorf("%s's %s: %v", p.Name, k, err)
			}

			pol.LastCheckedNow()
			continue
		}

		// Use default scraping if no custom scraping is set.
		if err := scrapePolicy(ctx, pol); err != nil {
			return fmt.Errorf("%s's %s: %v", p.Name, k, err)
		}
	}

	return nil
}

func scrapePolicy(ctx Ctx, p *policy) (err error) {
	p.WordCnt, err = getWordCount(ctx, p.URL, p.Sel)
	if err != nil {
		return
	}
	p.LastCheckedNow()
	return
}

func getWordCount(ctx Ctx, url, sel string) (int, error) {
	if url == "" {
		return 0, errEmptyURL
	}
	if sel == "" {
		return 0, errEmptySel
	}

	var text string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Text(sel, &text),
	)

	return len(strings.Fields(text)), err
}

func twtrCommGuides(c Ctx, pol *policy) error {
	ctx, cancel := chromedp.NewContext(c)
	defer cancel()

	pol.WordCnt = 0

	// Get all the links for the "Rules and Policies" index page.
	urls, err := twtrCommGuideURLs(ctx)
	if err != nil {
		return fmt.Errorf("getting twitter policy URLs: %v", err)
	}

	// Send the URLs to an unbuffered channel and close that channel after the
	// last URL value is read.
	urlCh := make(chan string)
	go func() {
		for _, url := range urls {
			urlCh <- url
		}
		close(urlCh)
	}()

	// Start a number of goroutines that read and process URL values from the
	// channel above until the channel is closed.
	var g errgroup.Group
	for i := 0; i < 3; i++ {
		g.Go(func() error {
			ctx, cancel := chromedp.NewContext(ctx)
			defer cancel()

			for {
				url, ok := <-urlCh
				if !ok {
					return nil
				}

				eachCtx, _ := context.WithTimeout(ctx, time.Minute)
				wc, err := getWordCount(eachCtx, url, ".twtr-grid__yang > div:nth-child(1)")
				if err != nil {
					return fmt.Errorf("%s: %v\n", url, err)
				}

				pol.WordCnt += wc
			}
		})
	}

	return g.Wait()
}

// getPolicyLinks extracts all of the links from the Twitter "rules and
// policies" index page. This is the umbrella under which lies the well-known /
// infamous policies such as Hateful Conduct and Abusive Behavior, as well as
// guides and resources for law enforcement.
func twtrCommGuideURLs(ctx Ctx) (urls []string, err error) {
	const prefix = "https://help.twitter.com"

	var nodes []*cdp.Node
	navigate := chromedp.Navigate(prefix + "/en/rules-and-policies")
	getNodes := chromedp.Nodes(".twtr-grid__yang a", &nodes)

	if err = chromedp.Run(ctx, navigate, getNodes); err != nil {
		return
	}

	for _, a := range nodes {
		urls = append(urls, prefix+a.AttributeValue("href"))
	}

	return
}

// ytCommGuides gets the word count from each Google support link that makes up
// YouTube's guidelines.
func ytCommGuides(c Ctx, pol *policy) error {
	const sel = "div.article-content-container"

	ctx, cancel := chromedp.NewContext(c)
	defer cancel()

	urls, err := ytCommGuideURLs(ctx)
	if err != nil {
		return fmt.Errorf("getting youtube policy URLs: %v", err)
	}

	var totalWc int
	for _, url := range urls {
		eachCtx, _ := context.WithTimeout(ctx, time.Minute)

		wc, err := getWordCount(eachCtx, url, sel)
		if err != nil {
			return fmt.Errorf("%s: %v\n", url, err)
		}

		totalWc += wc
	}

	pol.WordCnt = totalWc
	return nil
}

// ytCommGuideURLs extracts all of the Google support links from YouTube's
// "community guidelines" index page as well as the copyright index page (since
// that's it's own section).
func ytCommGuideURLs(ctx Ctx) ([]string, error) {
	targets := map[string]string{
		"https://www.youtube.com/about/policies/#community-guidelines":         "div.tabs__tab:nth-child(1) > div:nth-child(2) > div:nth-child(1) a",
		"https://www.youtube.com/about/copyright/#support-and-troubleshooting": "div.tabs__tab:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div > a",
	}

	var urls []string
	for page, sel := range targets {
		var nodes []*cdp.Node
		navigate := chromedp.Navigate(page)
		getNodes := chromedp.Nodes(sel, &nodes)

		if err := chromedp.Run(ctx, navigate, getNodes); err != nil {
			return nil, fmt.Errorf("%s: %v", page, err)
		}

		for _, a := range nodes {
			if url := a.AttributeValue("href"); isGoogleSupportLink(url) {
				urls = append(urls, url)
			}
		}
	}

	return urls, nil
}

func isGoogleSupportLink(url string) bool {
	return strings.HasPrefix(url, "https://support.google.com")
}
