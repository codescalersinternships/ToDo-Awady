package main

import (
	"encoding/json"
	"log"

	"net/http"

	_ "restapi/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DBFILE = "./data/todo.db"

type ToDoData struct {
	Text string
}

type ToDo struct {
	ID   uint
	Text string
}

func (ToDo) TableName() string {
	return "ToDo"
}

type ErrorResponse struct {
	Error string
}

type Server struct {
	DB *gorm.DB
}

func (s *Server) GetAllToDosHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.Model(&ToDo{}).Rows()

	if err != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var res ToDo
		err := s.DB.ScanRows(rows, &res)
		if err != nil {
			errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
			http.Error(w, string(errJSON), http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(res)
		if err != nil {
			errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
			http.Error(w, string(errJSON), http.StatusInternalServerError)
			return
		}
		w.Write(data)
		w.Write([]byte("\n"))

	}
}
func (s *Server) AddToDoHandler(w http.ResponseWriter, r *http.Request) {
	var todoData ToDoData
	err := json.NewDecoder(r.Body).Decode(&todoData)
	if err != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
	s.DB.Create(&ToDo{Text: todoData.Text})
	var res ToDo
	s.DB.Last(&res)
	data, err := json.Marshal(res)
	if err != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
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
		http.Error(w, "{\"error\":\""+string(errJson)+"\"}", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(res)
	if err != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
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
	s.DB.Create(&ToDo{Text: todoData.Text})
	var res ToDo
	result := s.DB.Model(&ToDo{}).Where("ID = ?", id).Update("Text", todoData.Text)
	if result.Error != nil {
		errJson, _ := json.Marshal(ErrorResponse{Error: result.Error.Error()})
		http.Error(w, "{\"error\":\""+string(errJson)+"\"}", http.StatusInternalServerError)
		return
	}
	e := s.DB.First(&res, id)
	if e.Error != nil {
		errJson, _ := json.Marshal(ErrorResponse{Error: e.Error.Error()})
		http.Error(w, "{\"error\":\""+string(errJson)+"\"}", http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(res)
	if err != nil {
		errJSON, _ := json.Marshal(ErrorResponse{Error: err.Error()})
		http.Error(w, string(errJSON), http.StatusInternalServerError)
		return
	}
	w.Write(data)
	w.Write([]byte("\n"))
}

func (s *Server) DeleteToDoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	result := s.DB.Delete(&ToDo{}, id)
	if result.Error != nil {
		errJson, _ := json.Marshal(ErrorResponse{Error: result.Error.Error()})
		http.Error(w, "{\"error\":\""+string(errJson)+"\"}", http.StatusInternalServerError)
		return
	}
	s.GetAllToDosHandler(w, r)

}

func main() {
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

	log.Fatal(http.ListenAndServe(":5000", r))

}
