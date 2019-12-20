package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Planting Seed")
	search := NewSearch("data2.json")
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Print Base Level\n")
	for i := 0; i < len(search.Children); i++ {
		fmt.Printf("%d: %c\n", i, search.Children[i].Value)
	}

	fmt.Println("Tree Grown")

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')

		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		result := search.DoSearch(text)
		fmt.Println(result)
	}
}
