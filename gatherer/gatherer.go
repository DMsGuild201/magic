package gatherer

import (
	"fmt"
	"log"
	"magic"
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
	rsp, err := g.Get(set.URL)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
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
	rsp, err := g.Get(c.URL)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		return err
	}

	c.ID = getCardID(c.URL)
	c.CardNumber = getCardNumber(doc)
	c.Names["en"] = getCardName(doc)
	c.Set = getCardSet(doc)
	c.Mana = getCardMana(doc)
	// c.Color = doc.Find("")
	c.Type = getCardType(doc)
	c.Rarity = getCardRarity(doc)
	c.ConvertedManageCost = getCardConvertedManaCost(doc)
	c.Power = getCardPower(doc)
	c.Toughness = getCardToughness(doc)
	c.Loyality = getCardLoyality(doc)

	c.AbilityTexts = getCardAbilityText(doc)
	c.FlavorText = getCardFlavorText(doc)
	c.Artist = getCardArtist(doc)
	c.Rulings = getCardRulings(doc)

	return nil
}

func getCardID(cu string) string {
	u, err := url.Parse(cu)
	if err != nil {
		log.Println(err)

		return ""
	}

	return u.Query().Get("multiverseid")
}

func getCardNumber(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_numberRow .value").Text())
}

func getCardName(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_nameRow .value").Text())
}

func getCardSet(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_currentSetSymbol a").Text())
}

func getCardType(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_typeRow .value").Text())
}

func getCardRarity(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_rarityRow .value").Text())
}

func getCardConvertedManaCost(doc *goquery.Document) int {
	val, _ := strconv.Atoi(strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_cmcRow .value").Text()))
	return val
}

func getCardPower(doc *goquery.Document) int {
	parts := strings.Split(strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ptRow .value").Text()), "/")
	if len(parts) != 2 {
		return 0
	}

	val, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	return val
}

func getCardToughness(doc *goquery.Document) int {
	parts := strings.Split(strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ptRow .value").Text()), "/")
	if len(parts) != 2 {
		return 0
	}

	val, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	return val
}

func getCardLoyality(doc *goquery.Document) int {
	val, _ := strconv.Atoi(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ptRow .value").Text())
	return val
}

func getCardFlavorText(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_FlavorText").Text())
}

func getCardArtist(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ArtistCredit").Text())
}

func getCardRulings(doc *goquery.Document) []string {
	var rules []string
	doc.Find(".rulingsText").Each(func(i int, s *goquery.Selection) {
		rules = append(rules, strings.TrimSpace(s.Text()))
	})
	return rules
}

func getCardMana(doc *goquery.Document) string {
	var mana []string
	doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_manaRow .value img").Each(func(i int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if ok {
			u, err := url.Parse(src)
			if err != nil {
				log.Println(err)
			}

			m := u.Query().Get("name")
			mana = append(mana, fmt.Sprintf("{%s}", m))
		}
	})

	return strings.Join(mana, "")
}

func getCardAbilityText(doc *goquery.Document) []string {
	var texts []string
	doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_textRow .value .cardtextbox").Each(func(i int, s *goquery.Selection) {
		texts = append(texts, s.Text())
	})
	return texts
}
