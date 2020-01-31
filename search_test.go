package main

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {

	fmt.Println("Planting Seed")
	search := NewSearch()
	search.PopulateJSON("testdata.json")
	fmt.Println("Tree Grown")

	text := "so"

	result := search.DoSimpleConcurrentSearch(text)

	expected := [4]string{"society", "social", "some", "so"}

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
