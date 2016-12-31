package main

import (
	"log"
	"magic"
	"magic/gatherer"
)

func main() {
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
