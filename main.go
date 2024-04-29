package main

import (
	"fmt"
	"log"

	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var movies []*Movie

// ----------------- structs -------------------
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// ---------------- controllers -----------------
func getAllMovies(w http.ResponseWriter, r *http.Request) {
	// set a header to let the client know data is in JSON format
	w.Header().Set("Content-Type", "application/json")

	// just give back the movies slice in JSON format
	json.NewEncoder(w).Encode(movies) // marshalling
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	// Set a header to let the client know data is in JSON format
	w.Header().Set("Content-Type", "application/json")

	// Get the movie ID from the URL
	params := mux.Vars(r)
	movieID := params["id"]

	// Initialize a variable to hold the deleted movie
	var deletedMovie *Movie

	// Find and delete the movie from the slice
	for index, movie := range movies {
		if movie.ID == movieID {
			// Save the deleted movie
			deletedMovie = movies[index]

			// Delete the movie from the slice
			movies = append(movies[:index], movies[index+1:]...)

			// Send response indicating movie was deleted
			fmt.Fprintf(w, "Movie successfully deleted:\n")

			// Break out of the loop once the movie is deleted
			break
		}
	}

	// If the movie was not found, return an error response
	if deletedMovie == nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	// Send the deleted movie and the updated movie slice as response
	response := struct {
		DeletedMovie *Movie   `json:"deleted_movie"`
		Movies       []*Movie `json:"movies"`
	}{
		DeletedMovie: deletedMovie,
		Movies:       movies,
	}

	// Marshal the response data to JSON and send it in the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getMovieById(w http.ResponseWriter, r *http.Request) {
	// set a header to let the client know data is in JSON format
	w.Header().Set("Content-Type", "application/json")

	// task: get the movie
	// get the id of movie from url
	params := mux.Vars(r)
	movieId := params["id"]
	// loop over slice and return the JSON encoded movie if exist
	for _, movie := range movies {
		if movie.ID == movieId {
			json.NewEncoder(w).Encode(movie) // marshalling
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	// set a header to let the client know data is in JSON format
	w.Header().Set("Content-Type", "application/json")

	// get the movie from req's body by conveerting it to a struct (un-marshalling)
	var newMovie Movie
	err := json.NewDecoder(r.Body).Decode(&newMovie) // un-marshalling
	if err != nil {
		http.Error(w, "JSON decode error", http.StatusBadRequest)
		return
	}

	// The incoming movie from request body will not have ID so we need to make it ourselves
	newMovie.ID = strconv.Itoa(rand.Intn(100000000000))
	// add new movie to slice
	movies = append(movies, &newMovie)

	// send the created movie back in response
	json.NewEncoder(w).Encode(newMovie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	/*
		The way i will update the movie is weird.
		I will delete the movie which is to be updated from slice first.
		Then i will create a new movie from the movie that came in request body.
		Then push that new movie to the slice.
		That's it.
	*/

	// set a header to let the client know data is in JSON format
	w.Header().Set("Content-Type", "application/json")

	// get the id
	params := mux.Vars(r)
	movieToUpdateId := params["id"]
	// loop over slice to get the movie
	for index, movie := range movies {
		if movie.ID == movieToUpdateId {
			// delete the movie
			movies = append(movies[:index], movies[index+1:]...)
			// un-marshall movie from req
			var newMovie Movie
			_ = json.NewDecoder(r.Body).Decode(&newMovie) // un-marshalling
			newMovie.ID = movieToUpdateId
			// append this update\new movie to slice
			movies = append(movies, &newMovie)
			// send back the response in JSON
			json.NewEncoder(w).Encode(newMovie) // marshalling
			return
		}
	}
}

// ----------------- MAIN ---------------------
func main() {
	// initializing gorilla mux
	r := mux.NewRouter()

	// Creating instances of Director
	director1 := &Director{Firstname: "Christopher", Lastname: "Nolan"}
	director2 := &Director{Firstname: "Quentin", Lastname: "Tarantino"}

	// Creating instances of Movie and filling the slice
	movies = []*Movie{
		{ID: "1", Isbn: "123456", Title: "Inception", Director: director1},
		{ID: "2", Isbn: "789012", Title: "Pulp Fiction", Director: director2},
	}

	// Define endpoints\routes
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("Key")
	}).Methods("GET")
	r.HandleFunc("/movie", getAllMovies).Methods("GET")
	r.HandleFunc("/movie/{id}", getMovieById).Methods("GET")
	r.HandleFunc("/movie", createMovie).Methods("POST")
	r.HandleFunc("/movie/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movie/{id}", deleteMovie).Methods("DELETE")

	PORT := ":8080"
	fmt.Printf("Server is running on PORT %v\n\n", PORT)
	if err := http.ListenAndServe(PORT, r); err != nil {
		log.Fatal(err)
	}
}
