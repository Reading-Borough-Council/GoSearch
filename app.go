package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const RESULTLIMIT = 16
const CONTENTLENGTH = 1023

type App struct {
	Router *mux.Router
	Search *search
}

type SearchResult struct {
	ID       int
	Rendered string
	URL      string
}

type FullSearchResult struct {
	ID      int
	Title   string
	Content string
	URL     string
}

func (a *App) Initialize(dataFile, siteMapFile string) {
	a.Search = NewSearch()
	a.Router = mux.NewRouter()

	a.Search.PopulateJSON(dataFile, siteMapFile)
	//a.Search.PopulateJSONStemmed(dataFile, siteMapFile)

	a.initializeRoutes()
}

func (a *App) Run(port int) {
	//Allow CORS
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port),
		handlers.CORS(handlers.AllowedHeaders(
			[]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(a.Router)))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/v1/search/short/{query}", a.searchHandler).Methods("GET")
	a.Router.HandleFunc("/v1/search/full/{query}", a.fullSearchHandler).Methods("GET")
	a.Router.HandleFunc("/v1/ping", a.ping).Methods("GET")
}

func (a *App) searchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	searchResults := make([]SearchResult, 0)
	fmt.Println(vars["query"])

	if vars["query"] != "" {
		if len(vars["query"]) < MINTERM {
			respondWithJSON(w, http.StatusOK, searchResults)
			return
		}
	}

	rawResults := a.Search.DoSimpleConcurrentSearch(strings.ToLower(vars["query"]), RESULTLIMIT)

	//array of possible words with ids
	//[{apple: [12, 43, 62]}, {application: [1, 43, 52]}]

	for _, r := range rawResults { //for each word
		for _, loc := range r.Location { //for each id
			title := a.Search.getArticleTitle(loc.ID)
			url := a.Search.getArticleURL(loc.ID)
			searchResult := SearchResult{ID: loc.ID, Rendered: title, URL: url}
			searchResults = append(searchResults, searchResult)
		}
	}

	respondWithJSON(w, http.StatusOK, searchResults)
}

func (a *App) fullSearchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	searchResults := make([]FullSearchResult, 0)
	fmt.Println(vars["query"])

	if vars["query"] != "" {
		if len(vars["query"]) < MINTERM {
			respondWithJSON(w, http.StatusOK, searchResults)
			return
		}
	}

	rawResults := a.Search.DoSimpleConcurrentSearch(strings.ToLower(vars["query"]), RESULTLIMIT)

	//array of possible words with ids
	//[{apple: [12, 43, 62]}, {application: [1, 43, 52]}]

	for _, r := range rawResults { //for each word
		for _, loc := range r.Location { //for each id
			title := a.Search.getArticleTitle(loc.ID)
			url := a.Search.getArticleURL(loc.ID)
			content := a.Search.getArticleContent(loc.ID)

			if len(content) > CONTENTLENGTH {
				content = content[:CONTENTLENGTH] + "..."
			}
			searchResult := FullSearchResult{ID: loc.ID, Title: title, Content: content, URL: url}
			searchResults = append(searchResults, searchResult)
		}
	}

	respondWithJSON(w, http.StatusOK, searchResults)
}

func (a *App) ping(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, "What!")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
