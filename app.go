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

type App struct {
	Router *mux.Router
	Search *Node
	Pages  []Page
}

type SearchResult struct {
	ID       int
	Rendered string
}

func (a *App) Initialize() {
	fmt.Println("Seed Planted")

	search := NewSearch()
	search.PopulateJSON("data.json")

	fmt.Println("Tree Grown")

	a.Router = mux.NewRouter()
	a.Search = search
	a.Pages = loadData("data.json")
	a.initializeRoutes()
}

func (a *App) Run(port int) {
	//log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), a.Router))

	//Allow CORS
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(a.Router)))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/v1/search/{term}", a.searchHandler).Methods("GET")
	a.Router.HandleFunc("/v1/ping", a.ping).Methods("GET")
}

func (a *App) searchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rawResults := a.Search.DoSearch(strings.ToLower(vars["term"]))

	searchResults := make([]SearchResult, 0)

	for _, r := range rawResults {
		for _, id := range r.ID {
			title := a.getArticleTitle(id)
			searchResult := SearchResult{ID: id, Rendered: title}
			searchResults = append(searchResults, searchResult)
		}
	}

	respondWithJSON(w, http.StatusOK, searchResults)
}

func (a *App) ping(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, "What!")
}

func (a *App) getArticleTitle(id int) string {
	for index := 0; index < len(a.Pages); index++ {
		if id == a.Pages[index].ID {
			return a.Pages[index].Title
		}
	}
	return ""
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
