package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "sync"

    "github.com/gorilla/mux"
)

type Task struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Done bool   `json:"done"`
}

var (
    tasks []Task
    idCounter int
    mu sync.Mutex
)

func getTasks(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
    var t Task
    json.NewDecoder(r.Body).Decode(&t)
    mu.Lock()
    idCounter++
    t.ID = idCounter
    tasks = append(tasks, t)
    mu.Unlock()
    json.NewEncoder(w).Encode(t)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, _ := strconv.Atoi(params["id"])
    mu.Lock()
    defer mu.Unlock()
    for i, t := range tasks {
        if t.ID == id {
            json.NewDecoder(r.Body).Decode(&tasks[i])
            tasks[i].ID = id
            json.NewEncoder(w).Encode(tasks[i])
            return
        }
    }
    http.Error(w, "Task not found", http.StatusNotFound)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, _ := strconv.Atoi(params["id"])
    mu.Lock()
    defer mu.Unlock()
    for i, t := range tasks {
        if t.ID == id {
            tasks = append(tasks[:i], tasks[i+1:]...)
            w.WriteHeader(http.StatusNoContent)
            return
        }
    }
    http.Error(w, "Task not found", http.StatusNotFound)
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/tasks", getTasks).Methods("GET")
    r.HandleFunc("/tasks", createTask).Methods("POST")
    r.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
    r.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

    log.Println("Server running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
