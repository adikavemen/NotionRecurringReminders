package main

type ResultsQuery struct {
	has_more    bool   `json:"has_more"`
	next_cursor string `json:"next_cursor"`
}
