package gatherer

import (
	"errors"
	"magic"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrNoParserFound = errors.New("no content parser found")
)

type CardParser interface {
	Parse(doc *goquery.Document) (*magic.Card, error)
}

// detect if we have special card cases
// examples are flip cards which contain
// two cards per page
func getCardParser(doc *goquery.Document) (CardParser, error) {
	switch doc.Find(".cardComponent").Length() {
	case 2: // this is a flip card
		return NewFlip(), nil

	case 0: // this is a normal card
		return NewSingle(), nil

	default:
		return nil, ErrNoParserFound
	}
}
