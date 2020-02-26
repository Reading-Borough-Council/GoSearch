package main

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {

	search := NewSearch()
	search.PopulateJSON("testdata.json", "testsitemap.json")

	text := "so"

	result := search.DoSimpleConcurrentSearch(text, 1)
	fmt.Println(result)

	expected := [1]string{"society"}

	fmt.Println("Run Tests w/ search: " + text)
	for _, str := range expected {
		found := false
		for _, res := range result {
			if str == res.Name {
				found = true
			}
		}

		if !found {
			t.Error("Fail, didn't get " + str)
		}
	}
}
