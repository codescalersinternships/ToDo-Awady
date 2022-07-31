package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type App struct {
	srv *http.Server
	db  *database
}

type ErrorResponse struct {
	Error string
}

type ToDoData struct {
	Text string
}

func NewApp(DBFILE, LISTENURL string, router http.Handler) (App, error) {
	a := App{}
	a.db = &database{}
	a.srv = &http.Server{}
	err := a.db.NewConnection(DBFILE)
	if err != nil {
		return a, err
	}
	a.srv.Addr = LISTENURL
	a.srv.Handler = router
	return a, nil
}

func (a *App) GetAllToDosHandler(w http.ResponseWriter, r *http.Request) {
	todos, err := a.db.GetAllToDos()
	if err != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(todos)
	w.Write([]byte("\n"))
}
func (a *App) AddToDoHandler(w http.ResponseWriter, r *http.Request) {
	var todoData ToDoData
	err := json.NewDecoder(r.Body).Decode(&todoData)
	if err != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
	todo, err := a.db.AddToDo(todoData.Text)
	if err != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(todo)
	w.Write([]byte("\n"))
}

func (a *App) GetToDoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	todo, err := a.db.GetToDo(id)
	if err == gorm.ErrRecordNotFound {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusNotFound)
		return
	}
	if err != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(todo)
	w.Write([]byte("\n"))
}

func (a *App) UpdateToDoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var todoData ToDoData
	err := json.NewDecoder(r.Body).Decode(&todoData)
	if err != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
	todo, err := a.db.UpdateToDo(id, todoData.Text)
	if err == gorm.ErrRecordNotFound {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusNotFound)
		return
	}
	if err != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(todo)
	w.Write([]byte("\n"))
}

func (a *App) DeleteToDoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := a.db.DeleteToDo(id)
	if err == gorm.ErrRecordNotFound {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusNotFound)
		fmt.Println(w.Header())
		return
	}
	if err != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		fmt.Println(w.Header())
		return
	}

	// w.WriteHeader(http.StatusNoContent)
}
