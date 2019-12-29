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
}

func (a *App) Initialize() {
	fmt.Println("Seed Planted")
	search := NewSearch("data.json")
	fmt.Println("Tree Grown")

	a.Router = mux.NewRouter()
	a.Search = search
	a.initializeRoutes()
}

func (a *App) Run(port int) {
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/v1/{term}", a.searchHandler).Methods("GET")
}

func (a *App) searchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	result := a.Search.DoSearch(vars["term"])
	
	respondWithJSON(w, http.StatusOK, result)
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