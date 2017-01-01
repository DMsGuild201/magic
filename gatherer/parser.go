package gatherer

import (
	"errors"
	"log"
	"magic"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrNoParserFound = errors.New("no content parser found")
)

type CardParser interface {
	Parse(doc *goquery.Document, c *magic.Card) error
}

// detect if we have special card cases
// examples are flip cards which contain
// two cards per page
func getCardParser(doc *goquery.Document) (CardParser, error) {
	switch doc.Find(".cardComponent").Length() {
	case 2: // this is a flip card
		log.Println("flip card..........")
		return flip{}, nil

	case 0: // this is a normal card
		return single{}, nil

	default:
		return nil, ErrNoParserFound
	}
}

func getCardID(cu string) string {
	u, err := url.Parse(cu)
	if err != nil {
		log.Println(err)

		return ""
	}

	return u.Query().Get("multiverseid")
}
