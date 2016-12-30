package main

type Card struct {
	// this is the multiverseid
	// taken from Gatherer
	ID  string
	URL string

	// this is a unique string for
	// the card in a set
	// edge cases with 27a as number
	CardNumber string
	Name       map[string]string
	Set        string
	Mana       string
	Color      string
	Type       string
	Rarity     Rarity

	// total that the card costs
	ConvertedManageCost int

	Power     int
	Toughness int
	Loyality  int

	AbilityText string
	FlavorText  string
	Artist      string
	Ruling      string

	// some cards have a backside
	// we need to represent this
	// see Westvale Abbey
	Backside *Card
}
