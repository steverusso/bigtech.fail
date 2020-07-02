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
