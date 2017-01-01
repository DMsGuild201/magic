package magic

type Set struct {
	Name      string `json:"name"`
	URL       string `json:"gatherer_url"`
	CardCount int    `json:"-"`

	Cards []*Card `json:"cards"`
}

func (s Set) String() string {
	return s.Name
}
