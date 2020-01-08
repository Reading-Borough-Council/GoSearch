package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	Search *Node
	Pages  []Page
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
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/v1/{term}", a.searchHandler).Methods("GET")
	a.Router.HandleFunc("/v1/ping", a.ping).Methods("GET")
}

func (a *App) searchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	result := a.Search.DoSearch(vars["term"])

	// for index := 0; index < len(result); index++ {
	// 	result[index].Title = a.getArticleTitle(result[index].ID)
	// }

	respondWithJSON(w, http.StatusOK, result)
}

func (a *App) ping(w http.ResponseWriter, r *http.Request) {

	respondWithJSON(w, http.StatusOK, "Success")
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
