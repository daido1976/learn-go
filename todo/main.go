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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// print access log.
	method := r.Method
	path := r.URL.Path
	log.Printf("[%s] %s\n", method, path)

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

func todoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleGet(w, r)
		// case "POST":
		// 	err = handlePost(w, r)
		// case "PUT":
		// 	err = handlePut(w, r)
		// case "DELETE":
		// 	err = handleDelete(w, r)
	}
}

// GET /todos/1
func handleGet(w http.ResponseWriter, r *http.Request) {
	// read json file
	var todos []Todo
	f, err := os.Open(FILENAME)
	check(err)
	dec := json.NewDecoder(f)
	dec.Decode(&todos)

	// retrive todo by id
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	check(err)
	todo, err := retrive(todos, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// render json
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(todo)
}

func retrive(todos []Todo, id int) (Todo, error) {
	for _, todo := range todos {
		if todo.ID == id {
			return todo, nil
		}
	}
	return Todo{}, errors.New("todo is not found")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/todos/", todoHandler)

	log.Println("[START] listen http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
