package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"myapp/internal/domain"
	"myapp/pkg/db"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

// CreateNoteHandler handles the creation of a new note
func CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	var note domain.Note

	// Decode the incoming JSON to the note struct
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connect to the database
	conn, err := db.ConnectToDB()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	// Write the SQL insert statement
	sqlStatement := `INSERT INTO notes (title, content) VALUES ($1, $2) RETURNING id, created_at, updated_at`

	// Execute the statement
	err = conn.QueryRow(context.Background(), sqlStatement, note.Title, note.Content).Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute the statement: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

// GetAllNotesHandler retrieves all notes from the database
func GetAllNotesHandler(w http.ResponseWriter, r *http.Request) {
	// Connect to the database
	conn, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	// Query all notes
	rows, err := conn.Query(context.Background(), "SELECT id, title, content, created_at, updated_at FROM notes")
	if err != nil {
		http.Error(w, "Failed to query notes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notes []domain.Note
	for rows.Next() {
		var note domain.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
			http.Error(w, "Failed to scan note", http.StatusInternalServerError)
			return
		}
		notes = append(notes, note)
	}

	// Respond with all notes in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// GetNoteByIDHandler retrieves a single note by its ID
func GetNoteByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the note ID from the URL path
	//id := r.URL.Path[len("/notes/"):]
	vars := mux.Vars(r)
	id := vars["id"]

	// Connect to the database
	conn, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	// Query the note by ID
	row := conn.QueryRow(context.Background(), "SELECT id, title, content, created_at, updated_at FROM notes WHERE id = $1", id)

	var note domain.Note
	err = row.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if err == pgx.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Failed to scan note", http.StatusInternalServerError)
		return
	}

	// Respond with the note in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

// UpdateNoteHandler updates an existing note by its ID
func UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the note ID from the URL path
	//id := r.URL.Path[len("/notes/"):]
	vars := mux.Vars(r)
	id := vars["id"]

	// Decode the incoming JSON to the note struct
	var updatedNote domain.Note
	if err := json.NewDecoder(r.Body).Decode(&updatedNote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connect to the database
	conn, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	// Update the note
	commandTag, err := conn.Exec(context.Background(), "UPDATE notes SET title = $1, content = $2, updated_at = NOW() WHERE id = $3", updatedNote.Title, updatedNote.Content, id)
	if err != nil {
		http.Error(w, "Failed to update note", http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() != 1 {
		http.NotFound(w, r)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Note updated successfully"})
}

// DeleteNoteHandler removes an existing note by its ID
func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the note ID from the URL path
	//id := r.URL.Path[len("/notes/"):]
	vars := mux.Vars(r)
	id := vars["id"]

	// Connect to the database
	conn, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	// Delete the note
	commandTag, err := conn.Exec(context.Background(), "DELETE FROM notes WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() != 1 {
		http.NotFound(w, r)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Note deleted successfully"})
}
