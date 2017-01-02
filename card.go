package magic

type Card struct {
	// this is the multiverseid
	// taken from Gatherer
	ID    string `json:"id"`
	URL   string `json:"gatherer_url"`
	Image string `json:"image"`

	// this is a unique string for
	// the card in a set
	// edge cases with 27a as number
	CardNumber string            `json:"card_number"`
	Names      map[string]string `json:"names"`
	Set        string            `json:"set"`
	Mana       string            `json:"mana"`
	Color      string            `json:"color"`
	Type       string            `json:"type"`
	Rarity     string            `json:"rarity"`

	// total that the card costs
	ConvertedManageCost int `json:"converted_mana_cost"`

	Power     int `json:"power"`
	Toughness int `json:"toughness"`
	Loyality  int `json:"loyality"`

	AbilityTexts []string `json:"ability_texts"`
	FlavorText   string   `json:"flavor_text"`
	Artist       string   `json:"artist"`
	Rulings      []string `json:"rulings"`

	// some cards have a backside
	// we need to represent this
	// see Westvale Abbey
	Backside *Card `json:"backside,omitempty"`
}

func (c Card) String() string {
	name, ok := c.Names["en"]
	if !ok {
		return "unknown name"
	}

	return name
}
