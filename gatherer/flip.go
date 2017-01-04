package gatherer

import (
	"magic"

	"github.com/PuerkitoBio/goquery"
)

type flip struct{}

func NewFlip() flip {
	return flip{}
}

func (s flip) Parse(doc *goquery.Document) (*magic.Card, error) {
	frontCol := "#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_cardComponent0"
	backCol := "#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_cardComponent1"

	front, err := NewCard(frontCol, doc).Parse()
	if err != nil {
		return nil, err
	}

	back, err := NewCard(backCol, doc).Parse()
	if err != nil {
		return nil, err
	}
	front.Backside = back

	return front, nil
}
