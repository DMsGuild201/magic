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

type flip struct{}

func (f flip) Parse(doc *goquery.Document, c *magic.Card) error {
	c.ID = getCardID(c.URL)
	c.CardNumber = f.getFrontCardNumber(doc)
	c.Names["en"] = f.getFrontCardName(doc)
	c.Set = f.getFrontCardSet(doc)
	c.Mana = f.getFrontCardMana(doc)
	c.Type = f.getFrontCardType(doc)
	c.Rarity = f.getFrontCardRarity(doc)
	c.ConvertedManageCost = f.getFrontCardConvertedManaCost(doc)
	c.Power = f.getFrontCardPower(doc)
	c.Toughness = f.getFrontCardToughness(doc)
	c.Loyality = f.getFrontCardLoyality(doc)
	c.AbilityTexts = f.getFrontCardAbilityText(doc)
	// c.FlavorText = s.getCardFlavorText(doc)
	c.Artist = f.getFrontCardArtist(doc)
	// c.Rulings = s.getCardRulings(doc)

	// create another card
	// same multiverseid
	c.Backside = &magic.Card{
		ID:         c.ID,
		URL:        c.URL,
		CardNumber: f.getBackCardNumber(doc),
		Image:      f.getBackCardImage(doc),
		Names: map[string]string{
			"en": f.getBackCardName(doc),
		},
		Set:                 f.getBackCardSet(doc),
		Mana:                f.getBackCardMana(doc),
		Type:                f.getBackCardType(doc),
		Rarity:              f.getBackCardRarity(doc),
		ConvertedManageCost: f.getBackCardConvertedManaCost(doc),
		Power:               f.getBackCardPower(doc),
		Toughness:           f.getBackCardToughness(doc),
		Loyality:            f.getBackCardLoyality(doc),
		AbilityTexts:        f.getBackCardAbilityText(doc),
		// FlavorText: "",
		Artist: f.getBackCardArtist(doc),
		// Rulings: []string{""},
	}

	return nil
}

// the front card on the flip card
func (f flip) getFrontCardNumber(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl03_numberRow .value").Text())
}

func (f flip) getFrontCardName(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_nameRow .value").Text())
}

func (f flip) getFrontCardSet(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl03_currentSetSymbol .value").Text())
}

func (f flip) getFrontCardMana(doc *goquery.Document) string {
	var mana []string
	doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl03_manaRow .value img").Each(func(i int, s *goquery.Selection) {
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

func (f flip) getFrontCardImage(doc *goquery.Document) string {
	src, ok := doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl02_Td1 img").Attr("src")
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

func (f flip) getFrontCardType(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl03_typeRow .value").Text())
}

func (f flip) getFrontCardRarity(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl03_rarityRow .value").Text())
}

func (f flip) getFrontCardConvertedManaCost(doc *goquery.Document) int {
	val, _ := strconv.Atoi(strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl03_cmcRow .value").Text()))
	return val
}

func (f flip) getFrontCardPower(doc *goquery.Document) int {
	parts := strings.Split(strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl02_ptRow .value").Text()), "/")
	if len(parts) != 2 {
		return 0
	}

	val, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	return val
}

func (f flip) getFrontCardToughness(doc *goquery.Document) int {
	parts := strings.Split(strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl02_ptRow .value").Text()), "/")
	if len(parts) != 2 {
		return 0
	}

	val, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	return val
}

func (f flip) getFrontCardLoyality(doc *goquery.Document) int {
	val, _ := strconv.Atoi(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl02_ptRow .value").Text())
	return val
}

func (f flip) getFrontCardAbilityText(doc *goquery.Document) []string {
	var texts []string
	doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl03_textRow .value .cardtextbox").Each(func(i int, s *goquery.Selection) {
		texts = append(texts, s.Text())
	})
	return texts
}

func (f flip) getFrontCardArtist(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl03_ArtistCredit .value").Text())
}

// get back side of the card
func (f flip) getBackCardName(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl04_nameRow .value").Text())
}

func (f flip) getBackCardMana(doc *goquery.Document) string {
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

func (f flip) getBackCardNumber(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl04_numberRow .value").Text())
}

func (f flip) getBackCardImage(doc *goquery.Document) string {
	src, ok := doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl03_Td1 img").Attr("src")
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

func (f flip) getBackCardSet(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("").Text())
}

func (f flip) getBackCardType(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl04_typeRow .value").Text())
}

func (f flip) getBackCardRarity(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl04_rarityRow .value").Text())
}

func (f flip) getBackCardConvertedManaCost(doc *goquery.Document) int {
	val, _ := strconv.Atoi(strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl04_cmcRow .value").Text()))
	return val
}

func (f flip) getBackCardPower(doc *goquery.Document) int {
	parts := strings.Split(strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl03_ptRow .value").Text()), "/")
	if len(parts) != 2 {
		return 0
	}

	val, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	return val
}

func (f flip) getBackCardToughness(doc *goquery.Document) int {
	parts := strings.Split(strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl03_ptRow .value").Text()), "/")
	if len(parts) != 2 {
		return 0
	}

	val, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	return val
}

func (f flip) getBackCardLoyality(doc *goquery.Document) int {
	return 0
}

func (f flip) getBackCardAbilityText(doc *goquery.Document) []string {
	var texts []string
	// doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl03_textRow .value .cardtextbox").Each(func(i int, s *goquery.Selection) {
	// 	texts = append(texts, s.Text())
	// })
	return texts
}

func (f flip) getBackCardArtist(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl04_ArtistCredit .value").Text())
}
