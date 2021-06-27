package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const TODO_FILE_NAME = "todo.json"

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
	todos := fetchCurrentTodos()

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
	body, err := ioutil.ReadFile(TODO_FILE_NAME)
	check(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// GET /todos/1
func handleShow(w http.ResponseWriter, r *http.Request) {
	id := extractIdFrom(r.URL)

	// retrive todo by id
	var todo Todo
	todo, err := todo.retrive(id)
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

	todos := fetchCurrentTodos()

	// build newTodos
	newId := len(todos) + 1
	newTodo := Todo{
		ID:    newId,
		Title: paramT.Title,
		Body:  paramT.Body,
	}
	newTodos := append(todos, newTodo)

	persist(newTodos)
	renderJson(w, newTodo)
}

// PUT /todos/1
func handleUpdate(w http.ResponseWriter, r *http.Request) {
	// decode params json to ParamTodo struct
	var paramT ParamTodo
	dec := json.NewDecoder(r.Body)
	dec.Decode(&paramT)

	id := extractIdFrom(r.URL)
	todos := fetchCurrentTodos()

	// update todos
	var newTodos []Todo
	var newTodo Todo
	for _, todo := range todos {
		if todo.ID == id {
			todo.Title = paramT.Title
			todo.Body = paramT.Body
			newTodo = todo
		}
		newTodos = append(newTodos, todo)
	}

	persist(newTodos)
	renderJson(w, newTodo)
}

// DELETE /todos/1
func handleDelete(w http.ResponseWriter, r *http.Request) {
	id := extractIdFrom(r.URL)
	todos := fetchCurrentTodos()

	// delete todo
	var newTodos []Todo
	for _, todo := range todos {
		if todo.ID != id {
			newTodos = append(newTodos, todo)
		}
	}

	persist(newTodos)
	w.WriteHeader(200)
}

// extract params id
func extractIdFrom(url *url.URL) int {
	id, err := strconv.Atoi(path.Base(url.Path))
	check(err)
	return id
}

// fetch current todos from json file
func fetchCurrentTodos() []Todo {
	var todos []Todo
	f, err := os.Open(TODO_FILE_NAME)
	check(err)
	dec := json.NewDecoder(f)
	dec.Decode(&todos)
	return todos
}

// persist by writing todos to json file
func persist(todos []Todo) {
	f, err := os.Create(TODO_FILE_NAME)
	check(err)
	enc := json.NewEncoder(f)
	enc.Encode(todos)
	defer f.Close()
}

// render json as http response
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
