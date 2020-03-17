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
	"sort"
	"strconv"
	"strings"
	"unicode"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/rookii/paicehusk"
)

const SEARCHLIMIT = 256
const MINTERM = 3

type search struct {
	Root    *node
	Pages   []page
	SiteMap []site
}

// node a node
type node struct {
	Children []*node
	Location []location
	Value    rune
}

// Location of word in site(ID) and body(Position)
type location struct {
	ID       int
	Position int
}

// Result with id for article
type result struct {
	Rendered string
	Title    string
	Name     string
	Location []location
}

// Page Json struct
type page struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Site
type site struct {
	ID  int    `json:ID`
	URL string `json:url`
}

// NewSearch create a node w/ no children
func NewSearch() *search {
	node := node{}
	search := search{Root: &node}
	return &search
}

// NewResultArray make
func NewResultArray() []result {
	locArr := make([]location, 0)
	initial := result{Name: "", Location: locArr}

	r := make([]result, 0)
	result := append(r, initial)
	return result
}

func NewNode(char rune) *node {
	return &node{Value: char}
}

func NewResult(r, t, n string, loc []location) *result {
	return &result{
		Rendered: r,
		Title:    t,
		Name:     n,
		Location: loc}
}

func NewPage(id int, title, content string) *page {
	return &page{
		ID:      id,
		Title:   title,
		Content: content}
}

// func (search *Search) DoComplexSearch(query string, count int) []result {
// 	search.D
// }33

// DoSimpleConcurrentSearch split up input and run search
// Get results for each individual term
// Return concurrent terms
func (search *search) DoSimpleConcurrentSearch(query string, count int) []result {
	results := NewResultArray()
	terms := strings.Split(query, " ")

	temp := make([]result, 0)
	output := make([]result, 0)

	if len(terms[0]) < MINTERM {
		return output
	}

	//get results for first term
	termResults := search.WordSearch(terms[0])

	resultCount := 0

	//for word/partial results i.e (app) => {application,applicator,appropo,...}
	for _, termResult := range termResults {
		results = append(results, termResult)

		for _, loc := range termResult.Location {

			location := make([]location, 1)
			location[0] = loc

			newResult := NewResult(termResult.Rendered,
				termResult.Title,
				termResult.Name,
				location)

			temp = append(temp, *newResult)
			resultCount += 1
		}

		if resultCount > RESULTLIMIT {
			break
		}
	}

	//now keep matching following terms
	followerCount := len(terms) - 1

	if followerCount > 0 {
		followers := terms[1:]

		for _, result := range temp {

			articleID := result.Location[0].ID
			articlePos := result.Location[0].Position

			text := strings.Split(strings.ToLower(search.getArticleTitle(articleID)), " ")
			valid := true

			//now check each following word
			for offset, term := range followers {
				txtIndex := articlePos + offset + 1

				if txtIndex < len(text) {
					match := text[txtIndex]

					if !strings.HasPrefix(match, term) {
						valid = false
						break
					} else {
						result.Rendered = result.Rendered + " " + term
					}
				} else {
					valid = false
					break
				}
			}

			if valid {
				output = append(output, result)
			}
		}
	} else {
		output = temp
	}

	//Now we have an array of all the results with concurrent terms
	//i.e hello world
	//if no results then return single word results by score
	//i.e if article contains hello and world
	//if still no results then return for single word results

	//Order results by match score i.e
	//Search: council tax
	//return:
	//1) council tax (2/2)
	//2) what is council tax (2/4)
	//3) i like chocolate and paying council tax (2/7)

	//If first term is beginning of sentence then prefer

	return output
}

// DoStemmedConcurrentSearch split up input and run search
// 1.
// Run Search on first stemmed term returns []result
// look for concurrent terms by location
// 2.
// Run search on all stemmed terms
// Look for concurrency
// 3.
// Run search on all
// Get results for each individual term
// Return concurrent terms
func (search *search) DoStemmedConcurrentSearch(query string, count int) []result {
	results := NewResultArray()
	terms := strings.Split(query, " ")

	for i, t := range terms {
		terms[i] = paicehusk.DefaultRules.Stem(t)
	}

	temp := make([]result, 0)

	//get results for first term
	termResults := search.WordSearch(terms[0])

	//now keep matching following terms
	followerCount := len(terms) - 1

	//for word/partial results i.e (app) => {application,applicator,appropo,...}
	for _, termResult := range termResults {
		results = append(results, termResult)

		for _, loc := range termResult.Location {

			location := make([]location, 1)
			location[0] = loc

			newResult := NewResult(termResult.Rendered,
				termResult.Title,
				termResult.Name,
				location)

			temp = append(temp, *newResult)
		}

	}

	output := make([]result, 0)

	if followerCount > 0 {
		followers := terms[1:]

		for _, result := range temp {

			articleID := result.Location[0].ID
			articlePos := result.Location[0].Position

			text := strings.Split(strings.ToLower(search.getArticleTitle(articleID)), " ")
			valid := true

			//now check each following word
			for offset, term := range followers {
				txtIndex := articlePos + offset + 1

				if txtIndex < len(text) {
					// match := text[txtIndex]
					match := paicehusk.DefaultRules.Stem(text[txtIndex])

					if !strings.HasPrefix(match, term) {
						valid = false
						break
					} else {
						result.Rendered = result.Rendered + " " + term
					}
				} else {
					valid = false
					break
				}
			}

			if valid {
				output = append(output, result)
			}
		}
	} else {
		output = temp
	}

	//Now we have an array of all the results with concurrent terms
	//i.e hello world
	//if no results then return single word results by score
	//i.e if article contains hello and world
	//if still no results then return for single word results

	//Order results by match score i.e
	//Search: council tax
	//return:
	//1) council tax (2/2)
	//2) what is council tax (2/4)
	//3) i like chocolate and paying council tax (2/7)

	//If first term is beginning of sentence then prefer

	return output
}

// WordSearch scan through node trie and return all possibilities
func (search *search) WordSearch(term string) []result {
	result := NewResultArray()
	found := false
	node := search.Root

	//scan leaves
	//move through tree until end of search term or not found
	for _, char := range term {

		//look for matching node
		found = false

		for _, child := range node.Children {
			thisChar := child.Value

			//move along
			if thisChar == rune(unicode.ToLower(char)) {
				found = true
				node = child
				result[0].Name = result[0].Name + string(unicode.ToLower(thisChar))
			}
		}

		if !found {
			//fmt.Println("Not Found")
			return result
		}
	}

	if found {
		//return results with node from end of term and prefix
		return getTree(node, result[0].Name)
	}

	return result
}

// getTree from end of term node find all branches
func getTree(node *node, str string) []result {
	result := make([]result, 0)

	if len(node.Location) > 0 {
		item := *NewResult(str, "", "", node.Location)
		result = append(result, item)
	}

	if len(node.Children) > 0 {
		for _, child := range node.Children {
			result = append(result, getTree(child, str+string(child.Value))...)
		}
	}

	return result
}

// func (search *search) FindNode(node *node, char rune) *node {

// }

// AddWord to Trie
func (search *search) AddWord(word string, location location) {
	//start at base node
	node := search.Root

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
			newNode := NewNode(wordChar)
			node.Children = append(node.Children, newNode)
			//traverse
			node = newNode
		}
	} // end char

	// word end, complete but id should
	// be array as there may be multiple articles with
	// the same words
	node.Location = append(node.Location, location)
}

// PopulateJSON Read JSON and add individual words
func (search *search) PopulateJSON(dataFilePath, siteMapPath string) {

	search.Pages = loadData(dataFilePath)
	search.SiteMap = loadSiteMap(siteMapPath)

	//now for each page
	for p := 0; p < len(search.Pages); p++ {

		//now add for each word of title type
		title := strings.Fields(search.Pages[p].Title)
		for index, word := range title {
			location := location{ID: search.Pages[p].ID, Position: index}
			search.AddWord(strings.ToLower(word), location)
		}

		//now add for each word of content type
		// content := strings.Fields(pages[p].Content)
		// for _, word := range content {
		// 	search.AddWord(word, pages[p].ID)
		// }
	}
}

// PopulateJSON Read JSON and add individual words stemmed
func (search *search) PopulateJSONStemmed(dataFilePath, siteMapPath string) {
	search.Pages = loadData(dataFilePath)
	search.SiteMap = loadSiteMap(siteMapPath)

	//now for each page
	for p := 0; p < len(search.Pages); p++ {

		//now add for each word of title type
		title := strings.Fields(search.Pages[p].Title)

		for index, word := range title {
			//now stem title
			wordStem := paicehusk.DefaultRules.Stem(word)

			location := location{ID: search.Pages[p].ID, Position: index}

			search.AddWord(strings.ToLower(wordStem), location)
		}

		//now add for each word of content type
		// content := strings.Fields(pages[p].Content)
		// for _, word := range content {
		// 	search.AddWord(word, pages[p].ID)
		// }
	}
}

// loadData, does what it says, loads json file returns array of 'pages'
func loadData(path string) []page {
	var dirtyPages []page
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println("Can't read file")
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &dirtyPages)
	var pages []page

	// get rid of trailing commas, full stops and colour codes
	regex := regexp.MustCompile("\\&#\\d*;|\\.^.|\\,|\\/^.|\\?|\\;|\\)|\\(|\\:")

	for _, page := range dirtyPages {
		title := strip.StripTags(page.Title)
		content := strip.StripTags(page.Content)

		title = regex.ReplaceAllString(title, "")
		content = regex.ReplaceAllString(content, "")

		cleanPage := *NewPage(page.ID, title, content)

		pages = append(pages, cleanPage)
	}

	fmt.Println("Page count: " + strconv.Itoa(len(pages)) + " (pages)")

	defer jsonFile.Close()

	//sort from low id to high
	sort.Slice(pages, func(i, j int) bool {
		return pages[i].ID < pages[j].ID
	})
	return pages
}

func loadSiteMap(path string) []site {
	var siteMap []site
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println("Can't read file")
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &siteMap)

	fmt.Println("Page count: " + strconv.Itoa(len(siteMap)) + " (sitemap)")

	defer jsonFile.Close()

	//sort from low id to high
	sort.Slice(siteMap, func(i, j int) bool {
		return siteMap[i].ID < siteMap[j].ID
	})
	return siteMap
}

//lazy binary search approx 1000pages
func (search *search) getArticleTitle(id int) string {
	//set index to reasonable value
	var max uint32 = uint32(len(search.Pages) - 1)
	var low uint32 = uint32(0)
	var index uint32 = max

	var countOut uint16 = SEARCHLIMIT

	for search.Pages[index].ID != id {
		//fmt.Printf("Search for %d @ index %d: %d (low: %d, max: %d) \n", id, index, search.Pages[index].ID, low, max)

		if search.Pages[index].ID > id {
			max = index
			index = ((max + low) >> 1)
		} else {
			low = index
			index = ((max + low) >> 1)
		}
	}

	countOut -= 1
	if countOut == 0 {
		return "-"
	}

	return search.Pages[index].Title
}

func (search *search) getArticleURL(id int) string {
	//set index to reasonable value
	var max uint32 = uint32(len(search.SiteMap) - 1)
	var low uint32 = uint32(0)
	var index uint32 = max

	var countOut uint16 = SEARCHLIMIT

	for search.SiteMap[index].ID != id {
		//fmt.Printf("Search for %d @ index %d: %d (low: %d, max: %d) \n", id, index, search.Pages[index].ID, low, max)

		if search.SiteMap[index].ID > id {
			max = index
			index = ((max + low) >> 1)
		} else {
			low = index
			index = ((max + low) >> 1)
		}
	}

	countOut -= 1
	if countOut == 0 {
		return "-"
	}

	return search.SiteMap[index].URL
}

func (search *search) getArticleContent(id int) string {
	//set index to reasonable value
	var max uint32 = uint32(len(search.Pages) - 1)
	var low uint32 = uint32(0)
	var index uint32 = max

	var countOut uint16 = SEARCHLIMIT

	for search.Pages[index].ID != id {
		//fmt.Printf("Search for %d @ index %d: %d (low: %d, max: %d) \n", id, index, search.Pages[index].ID, low, max)

		if search.Pages[index].ID > id {
			max = index
			index = ((max + low) >> 1)
		} else {
			low = index
			index = ((max + low) >> 1)
		}
	}

	countOut -= 1
	if countOut == 0 {
		return "-"
	}

	return search.Pages[index].Content
}
