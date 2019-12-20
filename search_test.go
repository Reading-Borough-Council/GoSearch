package main

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {

	fmt.Println("Planting Seed")
	search := NewSearch("data.json")
	fmt.Println("Tree Grown")

	text := "be"

	result := search.DoSearch(text)
	fmt.Println("Result for: " + text)
	for _, res := range result {
		fmt.Println(res)
	}

	expected := [3]string{"bean", "bear", "bertie"}

	fmt.Println("Tests")
	for _, str := range expected {
		found := false
		fmt.Println(str)
		for _, res := range result {
			if str == res {
				fmt.Println("Found")
				found = true
			}
		}

		if !found {
			t.Error("Fail")
		}
	}
}
