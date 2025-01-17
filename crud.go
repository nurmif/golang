package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type Item struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price float64 `json:"price"`
}

var items = make(map[string]Item)

func getItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	values := []Item{}
	for _, item := range items {
		values = append(values, item)
	}
	json.NewEncoder(w).Encode(values)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	item, exists := items[params["id"]]
	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(item)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newItem Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if _, exists := items[newItem.ID]; exists {
		http.Error(w, "Item already exists", http.StatusBadRequest)
		return
	}
	items[newItem.ID] = newItem
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	if _, exists := items[params["id"]]; !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	var updatedItem Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	updatedItem.ID = params["id"] // Ensure ID consistency
	items[params["id"]] = updatedItem
	json.NewEncoder(w).Encode(updatedItem)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	if _, exists := items[params["id"]]; !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	delete(items, params["id"])
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/items", getItems).Methods("GET")
	r.HandleFunc("/items/{id}", getItem).Methods("GET")
	r.HandleFunc("/items", createItem).Methods("POST")
	r.HandleFunc("/items/{id}", updateItem).Methods("PUT")
	r.HandleFunc("/items/{id}", deleteItem).Methods("DELETE")

	fmt.Println("Server started at :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
