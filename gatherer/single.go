package gatherer

import (
	"fmt"
	"log"
	"magic"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type single struct{}

func (s single) Parse(doc *goquery.Document, c *magic.Card) error {
	c.ID = getCardID(c.URL)
	c.CardNumber = s.getCardNumber(doc)
	c.Image = s.getCardImage(doc)
	c.Names["en"] = s.getCardName(doc)
	c.Set = s.getCardSet(doc)
	c.Mana = s.getCardMana(doc)
	// c.Color = doc.Find("")
	c.Type = s.getCardType(doc)
	c.Rarity = s.getCardRarity(doc)
	c.ConvertedManageCost = s.getCardConvertedManaCost(doc)
	c.Power = s.getCardPower(doc)
	c.Toughness = s.getCardToughness(doc)
	c.Loyality = s.getCardLoyality(doc)
	c.AbilityTexts = s.getCardAbilityText(doc)
	c.FlavorText = s.getCardFlavorText(doc)
	c.Artist = s.getCardArtist(doc)
	c.Rulings = s.getCardRulings(doc)

	return nil
}

func (s single) getCardNumber(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_numberRow .value").Text())
}

func (s single) getCardName(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_nameRow .value").Text())
}

func (s single) getCardImage(doc *goquery.Document) string {
	src, ok := doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_leftColumn img").Attr("src")
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

func (s single) getCardSet(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_currentSetSymbol a").Text())
}

func (s single) getCardType(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_typeRow .value").Text())
}

func (s single) getCardRarity(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_rarityRow .value").Text())
}

func (s single) getCardConvertedManaCost(doc *goquery.Document) int {
	val, _ := strconv.Atoi(strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_cmcRow .value").Text()))
	return val
}

func (s single) getCardPower(doc *goquery.Document) int {
	parts := strings.Split(strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ptRow .value").Text()), "/")
	if len(parts) != 2 {
		return 0
	}

	val, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	return val
}

func (s single) getCardToughness(doc *goquery.Document) int {
	parts := strings.Split(strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ptRow .value").Text()), "/")
	if len(parts) != 2 {
		return 0
	}

	val, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	return val
}

func (s single) getCardLoyality(doc *goquery.Document) int {
	val, _ := strconv.Atoi(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ptRow .value").Text())
	return val
}

func (s single) getCardFlavorText(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_FlavorText").Text())
}

func (s single) getCardArtist(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ArtistCredit").Text())
}

func (s single) getCardRulings(doc *goquery.Document) []string {
	var rules []string
	doc.Find(".rulingsText").Each(func(i int, s *goquery.Selection) {
		rules = append(rules, strings.TrimSpace(s.Text()))
	})
	return rules
}

func (s single) getCardMana(doc *goquery.Document) string {
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

func (s single) getCardAbilityText(doc *goquery.Document) []string {
	var texts []string
	doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_textRow .value .cardtextbox").Each(func(i int, s *goquery.Selection) {
		texts = append(texts, s.Text())
	})
	return texts
}
