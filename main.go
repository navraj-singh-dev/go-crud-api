package main

import (
	"fmt"
	"log"
	"strconv"
	"net/http"
	"math/rand"
	"encoding/json"
	"github.com/gorilla/mux"
)

type Movie struct {
	ID string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
}

func main() {
	// initializing gorilla mux
	r := mux.NewRouter()

	// Creating instances of Director
	director1 := &Director{Firstname: "Christopher", Lastname: "Nolan"}
	director2 := &Director{Firstname: "Quentin", Lastname: "Tarantino"}

	// Creating instances of Movie and filling the slice
	movies := []*Movie{
		{ID: "1", Isbn: "123456", Title: "Inception", Director: director1},
		{ID: "2", Isbn: "789012", Title: "Pulp Fiction", Director: director2},
	}

	// Define endpoints\routes
	r.HandleFunc("/movie", getAllMovies).Methods("GET")
	r.HandleFunc("/movie/{id}", getMovieById).Methods("GET")
	r.HandleFunc("/movie", createMovie).Methods("POST")
	r.HandleFunc("/movie/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movie/{id}", deleteMovie).Methods("DELETE")

	PORT := ":8080"
	fmt.Printf("Serve is running on PORT: %v", PORT)
	 if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal(err)
	 }
}
