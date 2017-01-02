package main

import (
	"log"
	"magic"
	"magic/gatherer"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	c := &magic.Card{
		URL:   "http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=955",
		Names: make(map[string]string),
	}
	e := gatherer.New().ScrapeCard(c)
	if e != nil {
		panic(e)
	}

	spew.Dump(*c)
	return

	var (
		sets  []*magic.Set
		cards []*magic.Card
	)

	g := gatherer.New()
	sets, err := g.ScrapeSets()
	if err != nil {
		panic(err)
	}

	log.Println("found", len(sets))
	counter := 5
	for _, s := range sets {
		log.Println("processing set", s)

		if cards, err = g.GetCards(s); err != nil {
			log.Println(err)
			continue
		}

		for _, card := range cards {
			if err = g.ScrapeCard(card); err != nil {
				log.Println(err)
				continue
			}

			log.Println("scraped", card, card.URL)

			// append the card to the set too
			s.Cards = append(s.Cards, card)
		}

		counter--
		if counter == 0 {
			break
		}
	}

	// save the set data
	setStore, err := JsonStorage("sets")
	if err != nil {
		panic(err)
	}

	if err := setStore.Save(sets); err != nil {
		panic(err)
	}

	// save the card data
	cardStore, err := JsonStorage("cards")
	if err != nil {
		panic(err)
	}

	if err := cardStore.Save(cards); err != nil {
		panic(err)
	}
	return
}
