# GoSearch
Golang Instant Search API from partial terms. Returns results as the client types. 
Very simple to use.

Given a list of pages with (id, title, content) builds an efficient search trie of every item in memory.
As partial terms are received returns most likely possibilities with id and position in the page.
As further words are typed i.e "blue ch" finds all occurrences of blue and then looks at preceding words to see if any start with "ch"

Input: "url/search/v1/"blue%20ch"

Result example:
ID: *article id
Text: "I like blue cheese"
URL: x/food/cheese

## Structure
Builds search trie/tree of each word in data with the article id and sentence and article position.
First strips all tags/html/colourcodes

Use stemming on all terms to allow for i.e (search, searches, searching)

Then returns url for result and highlights search text of article.

## Config
2020/01/16 Load data from JSON into memory keys(id, title, content)

## There are already searches why do I want to use this?
GoSearch can handle millions of requests with relatively low resources and very low latency.

## But why?
This is faster and better than everything.

## Run
go run main.go search.go app.go
or build and run ./search.exe

## Build
go build -o search.exe ./main.go ./app.go ./search.go

## Progress (TODO)
- 2020/01/16 Crawl sites
- 2020/02/05 Minimise payload size
-- Use stemming to return 

https://uxplanet.org/search-interface-20-things-to-consider-4b1466e98881
https://lucene.apache.org/solr/guide/7_4/phonetic-matching.html

