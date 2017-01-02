package gatherer

import (
	"fmt"
	"log"
	"magic"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"

	"github.com/PuerkitoBio/goquery"
)

type single struct{}

func (s single) Parse(doc *goquery.Document, c *magic.Card) error {
	data := parseCardColumn(doc, "#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_cardComponent0")

	c.ID = getCardID(c.URL)
	c.CardNumber = s.getCardNumber(data)
	c.Image = s.getCardImage(doc)
	c.Names["en"] = s.getCardName(data)
	c.Set = s.getCardSet(data)
	c.Mana = s.getCardMana(data)
	// c.Color = doc.Find("")
	c.Type = s.getCardType(data)
	c.Rarity = s.getCardRarity(data)
	c.ConvertedManageCost = s.getCardConvertedManaCost(data)
	c.Power = s.getCardPower(data)
	c.Toughness = s.getCardToughness(data)
	c.Loyality = s.getCardLoyality(data)
	c.AbilityTexts = s.getCardAbilityText(data)
	c.FlavorText = s.getCardFlavorText(data)
	c.Artist = s.getCardArtist(data)
	c.Rulings = s.getCardRulings(doc)

	return nil
}

func (s single) getCardNumber(data map[string]*goquery.Selection) string {
	sec, ok := data["Card Number:"]
	if !ok {
		return ""
	}

	return strings.TrimSpace(sec.Text())
}

func (s single) getCardName(data map[string]*goquery.Selection) string {
	sec, ok := data["Card Name:"]
	if !ok {
		return ""
	}

	return strings.TrimSpace(sec.Text())
}

func (s single) getCardImage(doc *goquery.Document) string {
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

func (s single) getCardSet(data map[string]*goquery.Selection) string {
	sec, ok := data["Expansion:"]
	if !ok {
		return ""
	}

	return strings.TrimSpace(sec.Text())
}

func (s single) getCardType(data map[string]*goquery.Selection) string {
	sec, ok := data["Types:"]
	if !ok {
		return ""
	}

	return strings.TrimSpace(sec.Text())
}

func (s single) getCardRarity(data map[string]*goquery.Selection) string {
	sec, ok := data["Rarity:"]
	if !ok {
		return ""
	}

	return strings.TrimSpace(sec.Text())
}

func (s single) getCardConvertedManaCost(data map[string]*goquery.Selection) int {
	sec, ok := data["Converted Mana Cost:"]
	if !ok {
		return 0
	}

	val, _ := strconv.Atoi(strings.TrimSpace(sec.Text()))
	return val
}

func (s single) getCardPower(data map[string]*goquery.Selection) *int {
	ret, ok := data["P/T:"]
	if !ok {
		return nil
	}

	parts := strings.Split(strings.TrimSpace(ret.Text()), "/")
	if len(parts) != 2 {
		return nil
	}

	val, _ := strconv.Atoi(parts[0])
	return &val
}

func (s single) getCardToughness(data map[string]*goquery.Selection) *int {
	ret, ok := data["P/T:"]
	if !ok {
		return nil
	}

	parts := strings.Split(ret.Text(), "/")
	if len(parts) != 2 {
		return nil
	}

	val, _ := strconv.Atoi(parts[1])
	return &val
}

func (s single) getCardLoyality(data map[string]*goquery.Selection) *int {
	ret, ok := data["Loyalty:"]
	if !ok {
		return nil
	}

	val, _ := strconv.Atoi(strings.TrimSpace(ret.Text()))
	return &val
}

func (s single) getCardFlavorText(data map[string]*goquery.Selection) *string {
	ret, ok := data["Flavor Text:"]
	if !ok {
		return nil
	}

	txt := strings.TrimSpace(ret.Text())
	return &txt
}

func (s single) getCardArtist(data map[string]*goquery.Selection) *string {
	ret, ok := data["Artist:"]
	if !ok {
		return nil
	}

	txt := strings.TrimSpace(ret.Text())
	return &txt
}

func (s single) getCardRulings(doc *goquery.Document) []string {
	var rules []string
	doc.Find(".rulingsTable .rulingsText").Each(func(i int, s *goquery.Selection) {
		rules = append(rules, strings.TrimSpace(s.Text()))
	})
	return rules
}

func (s single) getCardMana(data map[string]*goquery.Selection) string {
	var mana []string

	sec, ok := data["Mana Cost:"]
	if !ok {
		return ""
	}

	h, err := sec.Html()
	if err != nil {
		log.Println(err)
		return ""
	}

	if doc, err := html.Parse(strings.NewReader(h)); err == nil {
		var parser func(*html.Node)

		parser = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "img" {
				for _, el := range n.Attr {
					if el.Key == "src" {
						u, err := url.Parse(el.Val)
						if err != nil {
							log.Println(err)
							continue
						}

						m := u.Query().Get("name")
						mana = append(mana, fmt.Sprintf("{%s}", m))
					}
				}
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				parser(c)
			}
		}

		parser(doc)
	}

	return strings.Join(mana, "")
}

// TODO(styles): figure out why we're seeing both texts the mana and the mana with the string
// do we have multiple ability texts to worry about?
func (s single) getCardAbilityText(data map[string]*goquery.Selection) []string {
	var texts []string

	text, ok := data["Card Text:"]
	if !ok {
		return texts
	}

	h, err := text.Html()
	if err != nil {
		log.Println(err)
		return texts
	}

	if doc, err := html.Parse(strings.NewReader(h)); err == nil {
		var parser func(*html.Node)

		var output string
		parser = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "img" {
				for _, el := range n.Attr {
					if el.Key == "src" {
						u, err := url.Parse(el.Val)
						if err != nil {
							log.Println(err)
							continue
						}

						m := u.Query().Get("name")
						output += fmt.Sprintf("{%s}", m)
					}
				}
			}

			if n.Type == html.TextNode {
				output += n.Data
			}

			texts = append(texts, output)

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				parser(c)
			}
		}

		parser(doc)
	}

	return texts
}
