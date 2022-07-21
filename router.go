package main

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DBFILE = "todo.db"

type ToDo struct {
	ID   uint
	Text string
}

func (ToDo) TableName() string {
	return "ToDo"
}

func main() {
	db, err := gorm.Open(sqlite.Open(DBFILE), &gorm.Config{})
	if err != nil {
		panic("couldn't connect")
	}
	db.AutoMigrate(&ToDo{})

	r := mux.NewRouter()

	r.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
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

		}
	}).Methods("GET")

	r.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
		txt := r.FormValue("todo_text")
		db.Create(&ToDo{Text: txt})
		var res ToDo
		db.Last(&res)
		data, err := json.Marshal(res)
		if err != nil {
			panic("error parsing data")
		}
		w.Write(data)
	}).Methods("POST")

	r.HandleFunc("/todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		var res ToDo
		db.First(&res, id)
		data, err := json.Marshal(res)
		if err != nil {
			panic("error parsing data")
		}
		w.Write(data)
	}).Methods("GET")

	r.HandleFunc("/todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		txt := r.FormValue("todo_text")
		var res ToDo
		db.Model(&ToDo{}).Where("ID = ?", id).Update("Text", txt)
		db.First(&res, id)
		data, err := json.Marshal(res)
		if err != nil {
			panic("error parsing data")
		}
		w.Write(data)
	}).Methods("PUT")

	r.HandleFunc("/todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		result := db.Delete(&ToDo{}, id)
		if result.Error != nil {
			panic("error deleting data")
		}
		fmt.Fprintf(w, "%q Deleted", id)

	}).Methods("DELETE")
	http.ListenAndServe(":5000", r)

}
