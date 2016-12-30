package main

type Card struct {
	// this is the multiverseid
	// taken from Gatherer
	ID  string `json:"id"`
	URL string `json:"gatherer_url"`

	// this is a unique string for
	// the card in a set
	// edge cases with 27a as number
	CardNumber string            `json:"card_number"`
	Names      map[string]string `json:"names"`
	Set        string            `json:"set"`
	Mana       string            `json:"mana"`
	Color      string            `json:"color"`
	Type       string            `json:"type"`
	Rarity     Rarity            `json:"rarity"`

	// total that the card costs
	ConvertedManageCost int `json:"converted_mana_cost"`

	Power     int `json:"power"`
	Toughness int `json:"toughness"`
	Loyality  int `json:"loyality"`

	AbilityText string `json:"ability_text"`
	FlavorText  string `json:"flavor_text"`
	Artist      string `json:"artist"`
	Ruling      string `json:"ruling"`

	// some cards have a backside
	// we need to represent this
	// see Westvale Abbey
	Backside *Card `json:"backside,omitempty"`
}
