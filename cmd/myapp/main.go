package main

import (
	"fmt"
	"log"
	"myapp/internal/app/handlers" // Import the package where your handlers are
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("hello world")
	fmt.Println(os.Getenv("DATABASE_URL"))

	r := mux.NewRouter()

	// Setup routes
	r.HandleFunc("/notes", handlers.CreateNoteHandler).Methods("POST")
	r.HandleFunc("/notes", handlers.GetAllNotesHandler).Methods("GET")
	r.HandleFunc("/notes/{id}", handlers.GetNoteByIDHandler).Methods("GET")
	r.HandleFunc("/notes/{id}", handlers.UpdateNoteHandler).Methods("PUT")
	r.HandleFunc("/notes/{id}", handlers.DeleteNoteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
