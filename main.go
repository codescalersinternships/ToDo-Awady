package main

import (
	"encoding/json"
	"log"
	"os"

	"net/http"

	_ "restapi/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DBFILE = "todo.db"
var LISTENURL = ":5000"

type ToDoData struct {
	Text string
}

type ToDo struct {
	ID   uint   `json:"id"`
	Text string `json:"text"`
}

type ErrorResponse struct {
	Error string
}

type Server struct {
	DB *gorm.DB
}

func (s *Server) GetAllToDosHandler(w http.ResponseWriter, r *http.Request) {
	var rows []ToDo
	res := s.DB.Find(&rows)

	if res.Error != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: res.Error.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(rows)
	if err != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
	w.Write([]byte("\n"))
}
func (s *Server) AddToDoHandler(w http.ResponseWriter, r *http.Request) {
	var todoData ToDoData
	err := json.NewDecoder(r.Body).Decode(&todoData)
	if err != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
	var result *gorm.DB
	result = s.DB.Create(&ToDo{Text: todoData.Text})
	if result.Error != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: result.Error.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
	var res ToDo
	result = s.DB.Last(&res)
	if result.Error != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: result.Error.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(res)
	if err != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(data)
	w.Write([]byte("\n"))
}

func (s *Server) GetToDoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var res ToDo
	result := s.DB.First(&res, id)
	if result.Error != nil {
		errJson, _ := json.Marshal(ErrorResponse{Error: result.Error.Error()})
		http.Error(w, string(errJson), http.StatusNotFound)
		return
	}

	data, err := json.Marshal(res)
	if err != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
	w.Write([]byte("\n"))
}

func (s *Server) UpdateToDoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var todoData ToDoData
	err := json.NewDecoder(r.Body).Decode(&todoData)
	if err != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
	var res ToDo
	var result *gorm.DB
	result = s.DB.Model(&ToDo{}).Where("ID = ?", id).Update("Text", todoData.Text)
	if result.Error != nil {
		errJson, _ := json.Marshal(ErrorResponse{Error: result.Error.Error()})
		http.Error(w, string(errJson), http.StatusInternalServerError)
		return
	} else if result.RowsAffected < 1 {
		errJson, _ := json.Marshal(ErrorResponse{Error: "record not found"})
		http.Error(w, string(errJson), http.StatusNotFound)
		return

	}
	result = s.DB.First(&res, id)
	if result.Error != nil {
		errJson, _ := json.Marshal(ErrorResponse{Error: result.Error.Error()})
		http.Error(w, string(errJson), http.StatusNotFound)
		return
	}
	data, err := json.Marshal(res)
	if err != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(data)
	w.Write([]byte("\n"))
}

func (s *Server) DeleteToDoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	result := s.DB.Delete(&ToDo{}, id)
	if result.Error != nil {
		errJson, _ := json.Marshal(ErrorResponse{Error: result.Error.Error()})
		http.Error(w, string(errJson), http.StatusInternalServerError)
		return
	} else if result.RowsAffected < 1 {
		errJson, _ := json.Marshal(ErrorResponse{Error: "record not found"})
		http.Error(w, string(errJson), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	if f, ok := os.LookupEnv("DBFILE"); ok {
		DBFILE = f
	}
	if u, ok := os.LookupEnv("LISTENURL"); ok {
		LISTENURL = u
	}
	s := Server{}
	var err error
	s.DB, err = gorm.Open(sqlite.Open(DBFILE), &gorm.Config{})

	if err != nil {
		panic("couldn't connect")
	}
	s.DB.AutoMigrate(&ToDo{})

	r := mux.NewRouter()

	r.HandleFunc("/todo", s.GetAllToDosHandler).Methods("GET")
	r.HandleFunc("/todo", s.AddToDoHandler).Methods("POST")
	r.HandleFunc("/todo/{id}", s.GetToDoHandler).Methods("GET")
	r.HandleFunc("/todo/{id}", s.UpdateToDoHandler).Methods("PUT")
	r.HandleFunc("/todo/{id}", s.DeleteToDoHandler).Methods("DELETE")
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(LISTENURL, r))

}
