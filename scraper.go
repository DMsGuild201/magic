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

		if cards, err = g.GetCards(s); err != nil {
			log.Println(err)
			continue
		}

		for _, card := range cards {
			if err = g.ScrapeCard(card); err != nil {
				log.Println(err)
				continue
			}

			log.Printf("%+v", card)
		}

		break
	}

	sstor, err := JsonStorage("sets")
	if err != nil {
		panic(err)
	}

	if err := sstor.Save(sets); err != nil {
		panic(err)
	}
	cstor, err := JsonStorage("cards")
	if err != nil {
		panic(err)
	}

	if err := cstor.Save(cards); err != nil {
		panic(err)
	}
	return
}
