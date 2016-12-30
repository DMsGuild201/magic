package main

import "log"

func main() {
	var (
		sets  []*Set
		cards []*Card
	)

	g := NewGatherer()
	sets, err := g.ScrapeSets()
	if err != nil {
		panic(err)
	}

	log.Println("found", len(sets))
	for _, s := range sets {
		log.Println(s.Name, s.URL)

		cards, err = g.GetCards(s)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, card := range cards {
			log.Println(card.URL)
		}

		log.Println("found", len(cards))
	}
	return
}
