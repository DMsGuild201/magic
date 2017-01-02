package gatherer

import (
	"fmt"
	"log"
	"magic"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	api = "http://gatherer.wizards.com/"
)

type Service interface {
	SetList() ([]*magic.Set, error)
	GetCards(*magic.Set) ([]*magic.Card, error)
	ScrapeCard(*magic.Card) error
}

type gatherer struct {
	*http.Client
}

func New() gatherer {
	jar, _ := cookiejar.New(nil)

	// we need to set the client
	// to allow us to scrape without limits
	// 11=7 is the specific variable
	var cookies []*http.Cookie
	cookies = append(cookies, &http.Cookie{
		Name:   "CardDatabaseSettings",
		Value:  "0=1&1=28&2=0&14=1&3=13&4=0&5=1&6=15&7=0&8=1&9=1&10=19&11=7&12=8&15=1&16=0&13=",
		Path:   "/",
		Domain: "gatherer.wizards.com",
	})

	u, _ := url.Parse("http://gatherer.wizards.com")
	jar.SetCookies(u, cookies)

	return gatherer{
		Client: &http.Client{
			Jar: jar,
		},
	}
}

func (g gatherer) ScrapeSets() ([]*magic.Set, error) {
	rsp, err := g.Get(fmt.Sprintf("%sPages/Default.aspx", api))
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		return nil, err
	}

	var sets []*magic.Set
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

		set := magic.Set{
			Name: name,
			URL:  ep.String(),
		}

		sets = append(sets, &set)
	})

	return sets, nil
}

func (g gatherer) GetCards(set *magic.Set) ([]*magic.Card, error) {
	doc, err := goquery.NewDocument(set.URL)
	if err != nil {
		return nil, err
	}

	var cards []*magic.Card
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

		cards = append(cards, &magic.Card{
			Names: make(map[string]string),
			URL:   u.String(),
		})
	})

	return cards, nil
}

func (g gatherer) ScrapeCard(c *magic.Card) error {
	doc, err := goquery.NewDocument(c.URL)
	if err != nil {
		return err
	}

	// detect the card parser
	// flip card, normal card
	p, err := getCardParser(doc)
	if err != nil {
		return err
	}

	if err := p.Parse(doc, c); err != nil {
		return err
	}

	return nil
}

func parseCardColumn(doc *goquery.Document, col string) map[string]*goquery.Selection {
	data := make(map[string]*goquery.Selection)

	// we grab the goquery.Selection so that some of the fields can have their html
	// parsed out like mana images.
	doc.Find(".cardDetails .row").Each(func(i int, s *goquery.Selection) {
		label := strings.TrimSpace(s.Find(".label").Text())
		value := s.Find(".value")

		// we couldn't find the data for the card
		if value == nil {
			return
		}

		data[label] = value
	})

	return data
}
