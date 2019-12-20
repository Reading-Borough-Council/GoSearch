package main

/*
This package is used for instant comprehensive search
Data is loaded from json
Partial strings are received and queried against articles and posts data
If a suitable match is found at an adequate strength then the article/post Title
and an extract is returned

Populating:

*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Json struct
type Page struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// A Node
type Node struct {
	Children []*Node
	ID       int
	Value    rune
	Complete bool
}

// A NewNode
func NewNode(iD int, value rune, complete bool) *Node {
	node := Node{
		ID:       iD,
		Value:    value,
		Complete: complete}
	return &node
}

// NewSearch construct a search trie
func NewSearch(filePath string) *Node {
	baseNode := *NewNode(0, '#', false)

	//get data array from json
	var pages = getData(filePath)

	//now for each page
	for i := 0; i < len(pages); i++ {

		//now add for each word
		words := strings.Fields(pages[i].Content)
		for j := 0; j < len(words); j++ {

			//start at base node
			node := &baseNode

			//add to tries
			//for each character in a word look for it in the top level
			for i, rune := range words[j] {
				exists := false
				//scan leaves
				for k := 0; k < len(node.Children); k++ {
					thisChar := node.Children[k].Value

					//move along
					if thisChar == rune {
						exists = true
						node = node.Children[k]
					}
				}

				//add new node to children and move to it
				if !exists {
					newNode := NewNode(i, rune, false)
					node.Children = append(node.Children, newNode)
					node = newNode
				}
			}

			node.Complete = true
			node.ID =
		}
	}

	return &baseNode
}

//scan through node trie and return all possibilities
func (search *Node) DoSearch(term string) []string {
	result := make([]string, 0)

	fmt.Println("Searching with " + term)
	fmt.Println(strconv.Itoa(len(search.Children)) + " leaves.")

	prefix := ""
	result = append(result, prefix)

	//scan leaves
	//move through tree until end of search term or not found
	for termIndex := 0; termIndex < len(term)-1; termIndex++ {
		found := false

		for index := 0; index < len(search.Children); index++ {

			thisChar := search.Children[index].Value

			//fmt.Println(string(thisChar) + ", " + string(term[termIndex]))
			//move along
			if thisChar == rune(term[termIndex]) {
				search = search.Children[index]
				result[0] = prefix
				found = true
				break
			}
		}

		if !found {
			return result
		} else {
			prefix = prefix + string(search.Value)
		}
	}

	//return result
	return getTree(search, result[0])
}

func getTree(node *Node, str string) []string {
	result := make([]string, 0)
	str = str + string(node.Value)

	if node.Complete {
		result = append(result, str)
	}

	if len(node.Children) > 0 {
		for index := 0; index < len(node.Children); index++ {
			result = append(result, getTree(node.Children[index], str)...)
		}
	}

	return result
}

func getData(path string) []Page {
	var pages []Page
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println("Can't read file")
	} else {
		fmt.Println("File open")
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &pages)

	fmt.Println("Page count: " + strconv.Itoa(len(pages)))
	fmt.Println("Done")

	defer jsonFile.Close()
	return pages
}
