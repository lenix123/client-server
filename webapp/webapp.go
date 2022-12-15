package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Note struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	NoteText string `json:"note_text"`
}

var NoteStorage = []Note{}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintf(w, "Hi %s", r.URL.Query().Get("name"))
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func saveNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Bad Request.", http.StatusBadRequest)
		return
	}

	newNote := Note{}
	if json.Unmarshal(body, &newNote) != nil {
		http.Error(w, "Bad Request.", http.StatusBadRequest)
		return
	}
	NoteStorage = append(NoteStorage, newNote)
	fmt.Printf("Введённые данные: \n  имя: %s\n  фамилия: %s\n  заметка: %s\n", newNote.Name, newNote.Surname, newNote.NoteText)
}

func listAllNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	jsonResp, err := json.Marshal(NoteStorage)
	if err != nil {
		http.Error(w, "error happened", http.StatusInternalServerError)
	}

	w.Write(jsonResp)
	return
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	noteId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	
	numNotes := len(NoteStorage)
	if noteId < 0 || noteId >= numNotes {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if numNotes == 1 {
		NoteStorage = []Note{}
		return
	}
	
	NoteStorage[noteId] = NoteStorage[numNotes-1]
    NoteStorage = NoteStorage[:numNotes-1]
	return
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/save_note", saveNote)
	router.HandleFunc("/list_all", listAllNotes)
	router.HandleFunc("/delete_note/{id:[0-9]+}", deleteNote)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
