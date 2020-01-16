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
	"regexp"
	"strconv"
	"strings"
	"unicode"

	strip "github.com/grokify/html-strip-tags-go"
)

// Page Json struct
type Page struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Site struct {
	ID  int    `json:ID`
	URL string `json:url`
}

// Node a node
type Node struct {
	Children []*Node
	ID       []int
	Value    rune
}

// Result with id for article
type Result struct {
	Rendered string
	Title    string
	Name     string
	ID       []int
}

// NewNode create a node w/ no childrena
func NewSearch() *Node {
	node := Node{}
	return &node
}

func NewResultArray() []Result {
	result := make([]Result, 0)
	idArr := make([]int, 0)

	initial := Result{Name: "", ID: idArr}
	result = append(result, initial)
	return result
}

func (search *Node) AddWord(word string, id int) {
	//start at base node
	node := search

	//add to tries
	//for each character in a word look for it in the top level
	for _, wordChar := range word {
		exists := false
		//scan branches
		for b := 0; b < len(node.Children); b++ {
			thisChar := node.Children[b].Value

			//traverse
			if thisChar == wordChar {
				exists = true
				node = node.Children[b]
				break
			}
		}

		//add new node to children and move to it
		if !exists {
			//create node with character position for no particular reason
			newNode := &Node{Value: wordChar}
			node.Children = append(node.Children, newNode)
			//traverse
			node = newNode
		}
	} // end char

	// word end, complete but id should
	// be array as there may be multiple articles with
	// the same words
	node.ID = append(node.ID, id)
}

// NewSearch construct a search trie
func (search *Node) PopulateJSON(filePath string) {
	//get data array from json
	var pages = loadData(filePath)

	//now for each page
	for p := 0; p < len(pages); p++ {

		//now add for each word of title type
		title := strings.Fields(pages[p].Title)
		for _, word := range title {
			search.AddWord(strings.ToLower(word), pages[p].ID)
		}

		//now add for each word of title type
		// content := strings.Fields(pages[p].Content)
		// for _, word := range content {
		// 	search.AddWord(word, pages[p].ID)
		// }
	}
}

// DoSearch scan through node trie and return all possibilities
func (search *Node) DoSearch(term string) []Result {
	result := NewResultArray()

	//scan leaves
	//move through tree until end of search term or not found
	for _, char := range term {
		found := false

		//look for matching node
		for index := 0; index < len(search.Children); index++ {

			thisChar := search.Children[index].Value

			//move along
			if thisChar == rune(unicode.ToLower(char)) {
				search = search.Children[index]
				result[0].Name = result[0].Name + string(unicode.ToLower(thisChar))
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

	if len(node.ID) > 0 {
		item := Result{Name: str, ID: node.ID}
		result = append(result, item)
	}

	if len(node.Children) > 0 {
		for _, child := range node.Children {
			result = append(result, getTree(child, str+string(child.Value))...)
		}
	}

	return result
}

// loadData, does what it says, loads json file returns array of 'pages'
func loadData(path string) []Page {
	var dirtyPages []Page
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println("Can't read file")
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &dirtyPages)
	var pages []Page

	// get rid of trailing commas, full stops and colour codes
	regex := regexp.MustCompile("\\&#\\d*;|\\.^.|\\,|\\/^.|\\?|\\;|\\)|\\(|\\:")

	for _, page := range dirtyPages {
		title := strip.StripTags(page.Title)
		content := strip.StripTags(page.Content)

		title = regex.ReplaceAllString(title, "")
		content = regex.ReplaceAllString(content, "")

		cleanPage := Page{
			ID:      page.ID,
			Content: content,
			Title:   title}

		pages = append(pages, cleanPage)
	}

	fmt.Println("Page count: " + strconv.Itoa(len(pages)))

	defer jsonFile.Close()
	return pages
}

func loadSiteMap(path string) []Site {
	var siteMap []Site
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println("Can't read file")
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &siteMap)

	fmt.Println("Page count: " + strconv.Itoa(len(siteMap)))

	defer jsonFile.Close()
	return siteMap
}
