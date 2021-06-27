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

type ParamTodo struct {
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
	case "PUT":
		handleUpdate(w, r)
	case "DELETE":
		handleDelete(w, r)
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
	// enc.Encode(todos)
}

// GET /todos/1
func handleShow(w http.ResponseWriter, r *http.Request) {
	// extract params id
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	check(err)

	// retrive todo by id
	var todo Todo
	todo, err = todo.retrive(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	renderJson(w, todo)
}

// POST /todos
func handleCreate(w http.ResponseWriter, r *http.Request) {
	// decode params json to ParamTodo struct
	var paramT ParamTodo
	dec := json.NewDecoder(r.Body)
	dec.Decode(&paramT)

	// read json file
	var todos []Todo
	f, err := os.Open(FILENAME)
	check(err)
	dec = json.NewDecoder(f)
	dec.Decode(&todos)

	// build newTodos
	newId := len(todos) + 1
	newTodo := Todo{
		ID:    newId,
		Title: paramT.Title,
		Body:  paramT.Body,
	}
	newTodos := append(todos, newTodo)

	// write newTodos to json file
	f, err = os.Create(FILENAME)
	check(err)
	enc := json.NewEncoder(f)
	enc.Encode(newTodos)
	defer f.Close()

	renderJson(w, newTodo)
}

// PUT /todos/1
func handleUpdate(w http.ResponseWriter, r *http.Request) {}

// DELETE /todos/1
func handleDelete(w http.ResponseWriter, r *http.Request) {
	// extract params id
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	check(err)

	// read json file
	var todos []Todo
	f, err := os.Open(FILENAME)
	check(err)
	dec := json.NewDecoder(f)
	dec.Decode(&todos)

	// delete todo
	var newTodos []Todo
	for _, todo := range todos {
		if todo.ID != id {
			newTodos = append(newTodos, todo)
		}
	}

	// write newTodos to json file
	f, err = os.Create(FILENAME)
	check(err)
	enc := json.NewEncoder(f)
	enc.Encode(newTodos)
	defer f.Close()

	w.WriteHeader(200)
}

func renderJson(w http.ResponseWriter, t Todo) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(t)
}

func main() {
	http.HandleFunc("/todos", handleIndex)
	http.HandleFunc("/todos/", todoHandler)

	log.Println("[START] listen http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
