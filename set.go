package main

type Set struct {
	Name      string `json:"name"`
	URL       string `json:"gatherer_url"`
	CardCount int    `json:"-"`

	Cards []*Card `json:"-"`
}
