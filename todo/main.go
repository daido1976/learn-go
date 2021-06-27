package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const FILENAME = "todo.json"

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (t *Todo) retrive(id int) (Todo, error) {
	// read json file
	var todos []Todo
	f, err := os.Open(FILENAME)
	check(err)
	dec := json.NewDecoder(f)
	dec.Decode(&todos)

	// retrive by given id
	for _, todo := range todos {
		if todo.ID == id {
			return todo, nil
		}
	}
	return Todo{}, errors.New("todo is not found")
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
	// print access log.
	m := r.Method
	p := r.URL.Path
	log.Printf("[%s] %s\n", m, p)

	if m == "GET" && path.Base(p) == "todos" {
		handleIndex(w, r)
		return
	}

	switch m {
	case "GET":
		handleShow(w, r)
	case "POST":
		handleCreate(w, r)
		// case "PUT":
		// 	err = handlePut(w, r)
		// case "DELETE":
		// 	err = handleDelete(w, r)
	}
}

// GET /todos
func handleIndex(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadFile(FILENAME)
	check(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)

	// 以下でも同じ結果になる（一度 Go の Struct にデコード = Unmarshal してから、JSON にエンコード = Marshal し直して返す）
	// var todos []Todo
	// f, err := os.Open(FILENAME)
	// check(err)
	// dec := json.NewDecoder(f)
	// dec.Decode(&todos)

	// w.Header().Set("Content-Type", "application/json")
	// enc := json.NewEncoder(w)
	// enc.SetIndent("", "  ")
	// enc.Encode(todos)
}

// GET /todos/1
func handleShow(w http.ResponseWriter, r *http.Request) {
	// retrive todo by id
	var todo Todo
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	check(err)
	todo, err = todo.retrive(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	renderJson(w, todo)
}

// POST /todos/1
func handleCreate(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	dec := json.NewDecoder(r.Body)
	dec.Decode(&todo)

	renderJson(w, todo)
}

func renderJson(w http.ResponseWriter, t Todo) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(t)
}

func main() {
	http.HandleFunc("/todos", handleIndex)
	http.HandleFunc("/todos/", todoHandler)

	log.Println("[START] listen http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
