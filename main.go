package main

import (
	"encoding/json"
	"fmt"
	"log"

	"net/http"

	_ "restapi/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DBFILE = "todo.db"

var db, err = gorm.Open(sqlite.Open(DBFILE), &gorm.Config{})

type ToDoData struct {
	ToDoText string
}

type ToDo struct {
	ID   uint
	Text string
}

func (ToDo) TableName() string {
	return "ToDo"
}

func GetAllToDosHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Model(&ToDo{}).Where("ID > ?", "0").Rows()

	if err != nil {
		panic("error parsing data")
	}
	defer rows.Close()
	for rows.Next() {
		var res ToDo
		db.ScanRows(rows, &res)
		data, err := json.Marshal(res)
		if err != nil {
			panic("error parsing data")
		}
		w.Write(data)
		w.Write([]byte("\n"))

	}
}
func AddToDoHandler(w http.ResponseWriter, r *http.Request) {
	var todoData ToDoData
	err := json.NewDecoder(r.Body).Decode(&todoData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db.Create(&ToDo{Text: todoData.ToDoText})
	var res ToDo
	db.Last(&res)
	data, err := json.Marshal(res)
	if err != nil {
		panic("error parsing data")
	}
	w.Write(data)
	w.Write([]byte("\n"))
}

func GetToDoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var res ToDo
	db.First(&res, id)
	data, err := json.Marshal(res)
	if err != nil {
		panic("error parsing data")
	}
	w.Write(data)
	w.Write([]byte("\n"))
}

func UpdateToDoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var todoData ToDoData
	err := json.NewDecoder(r.Body).Decode(&todoData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db.Create(&ToDo{Text: todoData.ToDoText})
	var res ToDo
	db.Model(&ToDo{}).Where("ID = ?", id).Update("Text", todoData.ToDoText)
	db.First(&res, id)
	data, err := json.Marshal(res)
	if err != nil {
		panic("error parsing data")
	}
	w.Write(data)
	w.Write([]byte("\n"))
}

func DeleteToDoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	result := db.Delete(&ToDo{}, id)
	if result.Error != nil {
		panic("error deleting data")
	}
	fmt.Fprintf(w, "%q Deleted", id)

}

func main() {

	if err != nil {
		panic("couldn't connect")
	}
	db.AutoMigrate(&ToDo{})

	r := mux.NewRouter()

	r.HandleFunc("/todo", GetAllToDosHandler).Methods("GET")
	r.HandleFunc("/todo", AddToDoHandler).Methods("POST")
	r.HandleFunc("/todo/{id}", GetToDoHandler).Methods("GET")
	r.HandleFunc("/todo/{id}", UpdateToDoHandler).Methods("PUT")
	r.HandleFunc("/todo/{id}", DeleteToDoHandler).Methods("DELETE")
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(":5000", r))

}
