package gatherer

import (
	"magic"

	"github.com/PuerkitoBio/goquery"
)

type single struct{}

func NewSingle() single {
	return single{}
}

func (s single) Parse(doc *goquery.Document) (*magic.Card, error) {
	col := "#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_cardComponent0"
	return NewCard(col, doc).Parse()
}

// import (
// 	"fmt"
// 	"log"
// 	"magic"
// 	"net/url"
// 	"strconv"
// 	"strings"
// 	"time"
//
// 	"golang.org/x/net/html"
// 	"golang.org/x/net/html/atom"
//
// 	"github.com/PuerkitoBio/goquery"
// )
//
// type single struct{}
//
// func (s single) Parse(doc *goquery.Document, c *magic.Card) error {
// 	data := parseCardColumn(doc, "#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_cardComponent0")
//
// 	t := time.Now()
//
// 	c.ID = getCardID(c.URL)
// 	c.CardNumber = s.getCardNumber(data)
// 	c.Image = s.getCardImage(doc)
// 	c.Names["en"] = s.getCardName(data)
// 	c.Set = s.getCardSet(data)
// 	c.Mana = s.getCardMana(data)
// 	// c.Color = doc.Find("")
// 	c.Type = s.getCardType(data)
// 	c.Rarity = s.getCardRarity(data)
// 	c.ConvertedManageCost = s.getCardConvertedManaCost(data)
// 	c.Power = s.getCardPower(data)
// 	c.Toughness = s.getCardToughness(data)
// 	c.Loyality = s.getCardLoyality(data)
// 	c.AbilityText = s.getCardAbilityText(data)
// 	c.FlavorText = s.getCardFlavorText(data)
// 	c.Artist = s.getCardArtist(data)
// 	c.Rulings = s.getCardRulings(doc)
//
// 	log.Println(time.Now().Sub(t).String())
//
// 	return nil
// }
//
// func (s single) getCardNumber(data map[string]*goquery.Selection) string {
// 	sec, ok := data["Card Number:"]
// 	if !ok {
// 		return ""
// 	}
//
// 	return strings.TrimSpace(sec.Text())
// }
//
// func (s single) getCardName(data map[string]*goquery.Selection) string {
// 	sec, ok := data["Card Name:"]
// 	if !ok {
// 		return ""
// 	}
//
// 	return strings.TrimSpace(sec.Text())
// }
//
// func (s single) getCardImage(doc *goquery.Document) string {
// 	src, ok := doc.Find(".cardImage img").Attr("src")
// 	if !ok {
// 		return ""
// 	}
//
// 	ep, err := url.Parse(doc.Url.String())
// 	if err != nil {
// 		log.Println(err)
// 		return ""
// 	}
//
// 	ep, err = ep.Parse(src)
// 	if err != nil {
// 		log.Println(err)
// 		return ""
// 	}
//
// 	return ep.String()
// }
//
// func (s single) getCardSet(data map[string]*goquery.Selection) string {
// 	sec, ok := data["Expansion:"]
// 	if !ok {
// 		return ""
// 	}
//
// 	return strings.TrimSpace(sec.Text())
// }
//
// func (s single) getCardType(data map[string]*goquery.Selection) string {
// 	sec, ok := data["Types:"]
// 	if !ok {
// 		return ""
// 	}
//
// 	return strings.TrimSpace(sec.Text())
// }
//
// func (s single) getCardRarity(data map[string]*goquery.Selection) string {
// 	sec, ok := data["Rarity:"]
// 	if !ok {
// 		return ""
// 	}
//
// 	return strings.TrimSpace(sec.Text())
// }
//
// func (s single) getCardConvertedManaCost(data map[string]*goquery.Selection) int {
// 	sec, ok := data["Converted Mana Cost:"]
// 	if !ok {
// 		return 0
// 	}
//
// 	val, _ := strconv.Atoi(strings.TrimSpace(sec.Text()))
// 	return val
// }
//
// func (s single) getCardPower(data map[string]*goquery.Selection) string {
// 	ret, ok := data["P/T:"]
// 	if !ok {
// 		return ""
// 	}
//
// 	parts := strings.Split(strings.TrimSpace(ret.Text()), "/")
// 	if len(parts) != 2 {
// 		return ""
// 	}
//
// 	return strings.TrimSpace(parts[0])
// }
//
// func (s single) getCardToughness(data map[string]*goquery.Selection) string {
// 	ret, ok := data["P/T:"]
// 	if !ok {
// 		return ""
// 	}
//
// 	parts := strings.Split(ret.Text(), "/")
// 	if len(parts) != 2 {
// 		return ""
// 	}
//
// 	return strings.TrimSpace(parts[1])
// }
//
// func (s single) getCardLoyality(data map[string]*goquery.Selection) *int {
// 	ret, ok := data["Loyalty:"]
// 	if !ok {
// 		return nil
// 	}
//
// 	val, _ := strconv.Atoi(strings.TrimSpace(ret.Text()))
// 	return &val
// }
//
// func (s single) getCardFlavorText(data map[string]*goquery.Selection) *string {
// 	ret, ok := data["Flavor Text:"]
// 	if !ok {
// 		return nil
// 	}
//
// 	txt := strings.TrimSpace(ret.Text())
// 	return &txt
// }
//
// func (s single) getCardArtist(data map[string]*goquery.Selection) *string {
// 	ret, ok := data["Artist:"]
// 	if !ok {
// 		return nil
// 	}
//
// 	txt := strings.TrimSpace(ret.Text())
// 	return &txt
// }
//
// func (s single) getCardRulings(doc *goquery.Document) []string {
// 	var rules []string
// 	doc.Find(".rulingsTable .rulingsText").Each(func(i int, s *goquery.Selection) {
// 		rules = append(rules, strings.TrimSpace(s.Text()))
// 	})
// 	return rules
// }
//
// func (s single) getCardMana(data map[string]*goquery.Selection) string {
// 	sec, ok := data["Mana Cost:"]
// 	if !ok {
// 		return ""
// 	}
//
// 	h, err := sec.Html()
// 	if err != nil {
// 		log.Println(err)
// 		return ""
// 	}
//
// 	mana, err := s.convertManaFromHTML(h)
// 	if err != nil {
// 		log.Println(err)
// 		return ""
// 	}
//
// 	return mana
// }
//
// // TODO(styles): figure out why we're seeing both texts the mana and the mana with the string
// // do we have multiple ability texts to worry about?
// func (s single) getCardAbilityText(data map[string]*goquery.Selection) string {
// 	text, ok := data["Card Text:"]
// 	if !ok {
// 		return ""
// 	}
//
// 	h, err := text.Html()
// 	if err != nil {
// 		log.Println(err)
// 		return ""
// 	}
//
// 	out, err := s.convertManaFromHTML(h)
// 	if err != nil {
// 		log.Println(err)
// 		return ""
// 	}
//
// 	return out
// }
//
// func (s single) convertManaFromHTML(h string) (string, error) {
// 	// clean up html so that it's all on the same line
// 	h = strings.TrimSpace(h)
//
// 	var output string
// 	z := html.NewTokenizer(strings.NewReader(h))
// 	for {
// 		tt := z.Next()
//
// 		switch {
// 		// detect the end of the document
// 		case tt == html.ErrorToken:
// 			return output, nil
//
// 		case tt == html.SelfClosingTagToken:
// 			t := z.Token()
//
// 			// detect all the image tags
// 			if t.DataAtom == atom.Img {
//
// 				// for each image tag, grab it's attributes
// 				for _, el := range t.Attr {
// 					if el.Key == "src" {
// 						// parse out the attribute value, which
// 						// should be a link that contains the abbreviation
// 						u, err := url.Parse(el.Val)
// 						if err != nil {
// 							return output, err
// 						}
//
// 						m := u.Query().Get("name")
// 						output += fmt.Sprintf("{%s}", m)
// 					}
// 				}
// 			}
// 			continue
//
// 		// detect any text
// 		case tt == html.TextToken:
// 			t := z.Token()
// 			output += t.Data
// 		}
// 	}
// }
