package main

type Set struct {
	Name      string
	URL       string
	CardCount int

	Cards []*Card
}
