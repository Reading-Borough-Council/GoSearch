package main

import (
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	search := NewSearch()
	search.PopulateJSON("data.json", "sitemap.json")
	searchResults := make([]SearchResult, 0)

	rawResults := search.DoSimpleConcurrentSearch(strings.ToLower("appl"), 100)

	//array of possible words with ids
	//[{apple: [12, 43, 62]}, {application: [1, 43, 52]}]

	for _, r := range rawResults { //for each word
		for _, loc := range r.Location { //for each id
			title := search.getArticleTitle(loc.ID)
			url := search.getArticleURL(loc.ID)
			searchResult := SearchResult{ID: loc.ID, Rendered: title, URL: url}
			searchResults = append(searchResults, searchResult)
		}
	}

	println("test finished")
}
