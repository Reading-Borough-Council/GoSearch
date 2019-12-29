package main

// Author: Milo Bascombe (magicmilo)
// Date: 20/12/2019
// Copyright 2019 Reading Borough Council

// Trie search api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

//Page Json struct
type Page struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

//Node a node
type Node struct {
	Children []*Node
	ID       int
	Value    rune
	Complete bool
}

//Result with id for article
type Result struct {
	Name string
	ID   int
}

//NewNode create a node w/ no children
func NewNode(ID int, value rune, complete bool) *Node {
	node := Node{
		ID:       ID,
		Value:    value,
		Complete: complete}
	return &node
}

// NewSearch construct a search trie
func NewSearch(filePath string) *Node {
	baseNode := *NewNode(0, '#', false)

	//get data array from json
	var pages = loadData(filePath)

	//now for each page
	for p := 0; p < len(pages); p++ {

		//now add for each word
		words := strings.Fields(pages[p].Content)
		for w := 0; w < len(words); w++ {

			//start at base node
			node := &baseNode

			//add to tries
			//for each character in a word look for it in the top level
			for c, rune := range words[w] {
				exists := false
				//scan branches
				for b := 0; b < len(node.Children); b++ {
					thisChar := node.Children[b].Value

					//traverse
					if thisChar == rune {
						exists = true
						node = node.Children[b]
						break
					}
				}

				//add new node to children and move to it
				if !exists {
					//create node with character position for no particular reason
					newNode := NewNode(c, rune, false)
					node.Children = append(node.Children, newNode)
					//traverse
					node = newNode
				}
			} // end char

			// word end, complete but id should
			// be array as there may be multiple articles with
			// the same words 
			node.Complete = true
			node.ID = pages[p].ID
		}
	}

	return &baseNode
}

// DoSearch scan through node trie and return all possibilities
func (search *Node) DoSearch(term string) []Result {
	result := make([]Result, 0)

	fmt.Println("Searching with " + term)

	initial := Result{Name: "", ID: 0}
	result = append(result, initial)

	//scan leaves
	//move through tree until end of search term or not found
	for termIndex := 0; termIndex < len(term); termIndex++ {
		found := false

		//look for matching node
		for index := 0; index < len(search.Children); index++ {

			thisChar := search.Children[index].Value

			//move along
			if thisChar == rune(term[termIndex]) {
				search = search.Children[index]
				result[0].Name = result[0].Name + string(thisChar)
				found = true
				break
			}
		}

		if !found {
			return result
		}
	}

	//return results with node from end of term and prefix
	return getTree(search, result[0].Name)
}

// getTree from end of term node find all branches
func getTree(node *Node, str string) []Result {
	result := make([]Result, 0)	

	if node.Complete {
		item := Result{Name: str, ID: node.ID}
		result = append(result, item)
	}

	if len(node.Children) > 0 {
		for index := 0; index < len(node.Children); index++ {
			result = append(result, getTree(node.Children[index], str + string(node.Children[index].Value))...)
		}
	}

	return result
}

// loadData, does what it says
func loadData(path string) []Page {
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
