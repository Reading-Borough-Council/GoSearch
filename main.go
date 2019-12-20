package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// main live feedback
func main() {
	fmt.Println("Planting Seed")
	search := NewSearch("data.json")
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Print Base Level\n")
	for i := 0; i < len(search.Children); i++ {
		fmt.Printf("%d: %c\n", i, search.Children[i].Value)
	}

	fmt.Println("Tree Grown")

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')

		result := search.DoSearch(text)

		for index := 0; index < len(result); index++ {
			fmt.Println(strconv.Itoa(result[index].ID) + ": " + result[index].Name)
		}
	}
}
