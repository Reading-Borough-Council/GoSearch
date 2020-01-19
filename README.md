# GoSearch
Golang Instant Search API from partial terms. Returns results as the client types.
Very simple to use.

## Structure
Builds search trie/tree of each word in data with the article id and sentence and article position.
First strips all tags/html/colourcodes

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

## Progress
- 2020/01/16 Crawl sites, 
