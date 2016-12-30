package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

var (
	api = "http://gatherer.wizards.com/"
)

type GathererService interface {
	SetList() ([]*Set, error)
	GetCards(*Set) ([]*Card, error)
	ScrapeCard(string) (*Card, error)
}

type gatherer struct {
	*http.Client
}

func NewGatherer() gatherer {
	return gatherer{
		Client: &http.Client{},
	}
}

func (g gatherer) ScrapeSets() ([]*Set, error) {
	rsp, err := g.Get(fmt.Sprintf("%sPages/Default.aspx", api))
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		return nil, err
	}

	var sets []*Set
	doc.Find("#ctl00_ctl00_MainContent_Content_SearchControls_setAddText option").Each(func(i int, s *goquery.Selection) {
		name, ok := s.Attr("value")
		if !ok {
			log.Println("failed to parse value...")
			return
		}

		// skip the first empty value
		if name == "" {
			return
		}

		// we need to create the URL that
		// has a list of all the cards in it
		ep, err := url.Parse(api)
		if err != nil {
			log.Println(err)
			return
		}
		ep.Path = "Pages/Search/Default.aspx"

		q := ep.Query()
		q.Set("set", fmt.Sprintf(`["%s"]`, name))
		ep.RawQuery = q.Encode()

		set := Set{
			Name: name,
			URL:  ep.String(),
		}

		sets = append(sets, &set)
	})

	return sets, nil
}

func (g gatherer) GetCards(set *Set) ([]*Card, error) {
	rsp, err := g.Get(set.URL)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		return nil, err
	}

	var cards []*Card
	doc.Find(".middleCol a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok {
			log.Println("failed to get href for card")
			return
		}

		u, err := url.Parse(set.URL)
		if err != nil {
			log.Println(err)
			return
		}

		u, err = u.Parse(href)
		if err != nil {
			log.Println(err)
			return
		}

		cards = append(cards, &Card{
			URL: u.String(),
		})
	})

	return cards, nil
}

func (g gatherer) ScrapeCard(ep string) (*Card, error) {
	return nil, nil
}
