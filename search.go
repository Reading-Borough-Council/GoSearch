package main

/*
Author: magicmilo
Date: 20/12/2019

This package is used for instant comprehensive search
Data is loaded from json
Partial strings are received and queried against articles and posts data
If a suitable match is found at an adequate strength then the article/post Title
and an extract is returned

This already exists:
It's fun

Populating:
Provide a json file with parameters of id, title and content

Bugs:
- Double characters sometime duplicated i.e equip
2260: equipment,
2077: equipping
1315: equippment
1247: equippment.</td>
1247: equippped
828: equipppment
846: equipppment.
736: equipppment.</p>
840: equipppment</a>
1201: equipppment</a></p>
1174: equipppment</p>
846: equipppment</h2>
781: equipppment</li>
860: equipppment,
993: equipppped

*/

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
				//scan branches
				for k := 0; k < len(node.Children); k++ {
					thisChar := node.Children[k].Value

					//traverse
					if thisChar == rune {
						exists = true
						node = node.Children[k]
					}
				}

				//add new node to children and move to it
				if !exists {
					//create node with character position for no particular reason
					newNode := NewNode(i, rune, false)
					node.Children = append(node.Children, newNode)
					node = newNode
				}
			}

			node.Complete = true
			node.ID = pages[i].ID
		}
	}

	return &baseNode
}

// DoSearch scan through node trie and return all possibilities
func (search *Node) DoSearch(term string) []Result {
	result := make([]string, 0)

	fmt.Println("Searching with " + term)
	fmt.Println(strconv.Itoa(len(search.Children)) + " leaves.")

	prefix := ""
	result = append(result, prefix)

	//scan leaves
	//move through tree until end of search term or not found
	for termIndex := 0; termIndex < len(term)-1; termIndex++ {
		found := false

		//look for matching node
		for index := 0; index < len(search.Children); index++ {

			thisChar := search.Children[index].Value

			//move along
			if thisChar == rune(term[termIndex]) {
				search = search.Children[index]
				result[0] = prefix
				found = true
				break
			}
		}

		if found {
			prefix = prefix + string(search.Value)
		}
	}

	//return results with node from end of term and prefix
	return getTree(search, result[0])
}

// getTree from end of term node find all branches
func getTree(node *Node, str string) []Result {
	result := make([]Result, 0)
	str = str + string(node.Value)

	if node.Complete {
		item := Result{Name: str, ID: node.ID}
		result = append(result, item)
	}

	if len(node.Children) > 0 {
		for index := 0; index < len(node.Children); index++ {
			result = append(result, getTree(node.Children[index], str)...)
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
