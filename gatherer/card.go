package gatherer

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"magic"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/PuerkitoBio/goquery"
)

type card struct {
	column string

	doc *goquery.Document
}

func NewCard(col string, doc *goquery.Document) card {
	return card{
		column: col,

		doc: doc,
	}
}

func (s card) Parse() (*magic.Card, error) {
	// parse the card data and push it around
	data := s.parseCardColumn()

	c := magic.NewCard()

	c.ID = s.getCardID(s.doc.Url.String())
	c.URL = s.doc.Url.String()
	c.CardNumber = s.getCardNumber(data)
	c.Image = s.getCardImage(s.doc)
	c.Names["en"] = s.getCardName(data)
	c.Set = s.getCardSet(data)
	c.Mana = s.getCardMana(data)
	c.Type = s.getCardType(data)
	c.Rarity = s.getCardRarity(data)
	c.ConvertedManageCost = s.getCardConvertedManaCost(data)
	c.Power = s.getCardPower(data)
	c.Toughness = s.getCardToughness(data)
	c.Loyality = s.getCardLoyality(data)
	c.AbilityTexts = s.getCardAbilityTexts(data)
	c.FlavorText = s.getCardFlavorText(data)
	c.Artist = s.getCardArtist(data)
	c.Rulings = s.getCardRulings(s.doc)

	return c, nil
}

func (s card) parseCardColumn() map[string]*goquery.Selection {
	data := make(map[string]*goquery.Selection)

	// we grab the goquery.Selection so that some of the fields can have their html
	// parsed out like mana images.
	s.doc.Find(".cardDetails .row").Each(func(i int, s *goquery.Selection) {
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

func (s card) getCardID(cu string) string {
	log.Println("getting card id", cu)
	u, err := url.Parse(cu)
	if err != nil {
		log.Println(err)

		return ""
	}

	return u.Query().Get("multiverseid")
}

func (s card) getCardNumber(data map[string]*goquery.Selection) string {
	sec, ok := data["Card Number:"]
	if !ok {
		return ""
	}

	return strings.TrimSpace(sec.Text())
}

func (s card) getCardName(data map[string]*goquery.Selection) string {
	sec, ok := data["Card Name:"]
	if !ok {
		return ""
	}

	return strings.TrimSpace(sec.Text())
}

func (s card) getCardImage(doc *goquery.Document) string {
	src, ok := doc.Find(".cardImage img").Attr("src")
	if !ok {
		return ""
	}

	ep, err := url.Parse(doc.Url.String())
	if err != nil {
		log.Println(err)
		return ""
	}

	ep, err = ep.Parse(src)
	if err != nil {
		log.Println(err)
		return ""
	}

	return ep.String()
}

func (s card) getCardSet(data map[string]*goquery.Selection) string {
	sec, ok := data["Expansion:"]
	if !ok {
		return ""
	}

	return strings.TrimSpace(sec.Text())
}

func (s card) getCardType(data map[string]*goquery.Selection) string {
	sec, ok := data["Types:"]
	if !ok {
		return ""
	}

	return strings.TrimSpace(sec.Text())
}

func (s card) getCardRarity(data map[string]*goquery.Selection) string {
	sec, ok := data["Rarity:"]
	if !ok {
		return ""
	}

	return strings.TrimSpace(sec.Text())
}

func (s card) getCardConvertedManaCost(data map[string]*goquery.Selection) int {
	sec, ok := data["Converted Mana Cost:"]
	if !ok {
		return 0
	}

	val, _ := strconv.Atoi(strings.TrimSpace(sec.Text()))
	return val
}

func (s card) getCardPower(data map[string]*goquery.Selection) string {
	ret, ok := data["P/T:"]
	if !ok {
		return ""
	}

	parts := strings.Split(strings.TrimSpace(ret.Text()), "/")
	if len(parts) != 2 {
		return ""
	}

	return strings.TrimSpace(parts[0])
}

func (s card) getCardToughness(data map[string]*goquery.Selection) string {
	ret, ok := data["P/T:"]
	if !ok {
		return ""
	}

	parts := strings.Split(ret.Text(), "/")
	if len(parts) != 2 {
		return ""
	}

	return strings.TrimSpace(parts[1])
}

func (s card) getCardLoyality(data map[string]*goquery.Selection) *int {
	ret, ok := data["Loyalty:"]
	if !ok {
		return nil
	}

	val, _ := strconv.Atoi(strings.TrimSpace(ret.Text()))
	return &val
}

func (s card) getCardFlavorText(data map[string]*goquery.Selection) *string {
	ret, ok := data["Flavor Text:"]
	if !ok {
		return nil
	}

	txt := strings.TrimSpace(ret.Text())
	return &txt
}

func (s card) getCardArtist(data map[string]*goquery.Selection) *string {
	ret, ok := data["Artist:"]
	if !ok {
		return nil
	}

	txt := strings.TrimSpace(ret.Text())
	return &txt
}

func (s card) getCardRulings(doc *goquery.Document) []string {
	var rules []string
	doc.Find(".rulingsTable .rulingsText").Each(func(i int, s *goquery.Selection) {
		rules = append(rules, strings.TrimSpace(s.Text()))
	})
	return rules
}

func (s card) getCardMana(data map[string]*goquery.Selection) string {
	sec, ok := data["Mana Cost:"]
	if !ok {
		return ""
	}

	h, err := sec.Html()
	if err != nil {
		log.Println(err)
		return ""
	}

	mana, err := s.convertManaFromHTML(h)
	if err != nil {
		log.Println(err)
		return ""
	}

	return mana
}

// TODO(styles): figure out why we're seeing both texts the mana and the mana with the string
// do we have multiple ability texts to worry about?
func (s card) getCardAbilityTexts(data map[string]*goquery.Selection) []string {
	var output []string

	text, ok := data["Card Text:"]
	if !ok {
		return output
	}

	h, err := text.Html()
	if err != nil {
		log.Println(err)
		return output
	}

	doc, err := html.Parse(strings.NewReader(h))
	if err != nil {
		log.Println(err)
		return output
	}

	var dirty []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.DataAtom == atom.Div {
			var buf bytes.Buffer
			w := io.Writer(&buf)
			html.Render(w, n)
			dirty = append(dirty, buf.String())
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	for _, d := range dirty {
		out, err := s.convertManaFromHTML(d)
		if err != nil {
			log.Println(err)
			return output
		}

		output = append(output, out)
	}

	return output
}

func (s card) convertManaFromHTML(h string) (string, error) {
	// clean up html so that it's all on the same line
	h = strings.TrimSpace(h)

	var output string
	z := html.NewTokenizer(strings.NewReader(h))
	for {
		tt := z.Next()

		switch {
		// detect the end of the document
		case tt == html.ErrorToken:
			return output, nil

		case tt == html.SelfClosingTagToken:
			t := z.Token()

			// detect all the image tags
			if t.DataAtom == atom.Img {

				// for each image tag, grab it's attributes
				for _, el := range t.Attr {
					if el.Key == "src" {
						// parse out the attribute value, which
						// should be a link that contains the abbreviation
						u, err := url.Parse(el.Val)
						if err != nil {
							return output, err
						}

						m := u.Query().Get("name")
						output += fmt.Sprintf("{%s}", m)
					}
				}
			}
			continue

		// detect any text
		case tt == html.TextToken:
			t := z.Token()
			output += t.Data
		}
	}
}
