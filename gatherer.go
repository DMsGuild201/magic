package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	api = "http://gatherer.wizards.com/"
)

type GathererService interface {
	SetList() ([]*Set, error)
	GetCards(*Set) ([]*Card, error)
	ScrapeCard(*Card) error
}

type gatherer struct {
	*http.Client
}

func NewGatherer() gatherer {
	jar, _ := cookiejar.New(nil)

	// we need to set the client
	// to allow us to scrape without limits
	// 11=7 is the specific variable
	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:   "CardDatabaseSettings",
		Value:  "0=1&1=28&2=0&14=1&3=13&4=0&5=1&6=15&7=0&8=1&9=1&10=19&11=7&12=8&15=1&16=0&13=",
		Path:   "/",
		Domain: "gatherer.wizards.com",
	}
	cookies = append(cookies, cookie)

	u, _ := url.Parse("http://gatherer.wizards.com")
	jar.SetCookies(u, cookies)

	return gatherer{
		Client: &http.Client{
			Jar: jar,
		},
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

func (g gatherer) ScrapeCard(c *Card) error {
	rsp, err := g.Get(c.URL)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		return err
	}

	u, err := url.Parse(c.URL)
	if err != nil {
		return err
	}

	c.ID = u.Query().Get("multiverseid")
	// c.CardNumber
	c.CardNumber = strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_numberRow .value").Text())
	// c.Names["en"] = doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_nameRow .value").Text()
	c.Set = strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_currentSetSymbol a").Text())
	// c.Mana
	// c.Color = doc.Find("")
	c.Type = strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_typeRow .value").Text())
	c.Rarity = strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_rarityRow .value").Text())
	c.ConvertedManageCost, _ = strconv.Atoi(strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_cmcRow .value").Text()))
	c.Power, _ = strconv.Atoi(doc.Find("").Text())
	c.Toughness, _ = strconv.Atoi(doc.Find("").Text())
	c.Loyality, _ = strconv.Atoi(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ptRow .value").Text())

	// c.AbilityText
	c.FlavorText = strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_FlavorText").Text())
	c.Artist = strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ArtistCredit").Text())
	// c.Rulings

	// ID
	// URL
	//
	// CardNumber
	// Names
	// Set
	// Mana
	// Color
	// Type
	// Rarity
	// ConvertedManaCost
	//
	// Power
	// Toughness
	// Loyality
	//
	// AbilityText
	// FlavorText
	// Artist
	// Ruling
	//
	// Backside

	return nil
}
