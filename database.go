package main

import (
	"encoding/json"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type database struct {
	db *gorm.DB
}

type ToDo struct {
	ID   uint   `json:"id"`
	Text string `json:"text"`
}

func (d *database) NewConnection(DBFILE string) error {
	var err error
	d.db, err = gorm.Open(sqlite.Open(DBFILE), &gorm.Config{})
	if err != nil {
		return err
	}
	err = d.db.AutoMigrate(&ToDo{})
	return err
}

func (d *database) GetAllToDos() ([]byte, error) {
	var todos []ToDo
	res := d.db.Find(&todos)

	if res.Error != nil {
		return nil, res.Error
	}
	data, err := json.Marshal(todos)
	return data, err

}

func (d *database) GetToDo(id string) ([]byte, error) {
	var res ToDo
	result := d.db.First(&res, id)
	if result.Error != nil {
		return nil, result.Error
	}

	data, err := json.Marshal(res)
	return data, err
}

func (d *database) AddToDo(text string) ([]byte, error) {
	var result *gorm.DB
	result = d.db.Create(&ToDo{Text: text})
	if result.Error != nil {
		return nil, result.Error
	}
	var res ToDo
	result = d.db.Last(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	data, err := json.Marshal(res)
	return data, err
}

func (d *database) UpdateToDo(id, text string) ([]byte, error) {
	var res ToDo
	var result *gorm.DB
	result = d.db.Model(&ToDo{}).Where("ID = ?", id).Update("Text", text)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected < 1 {
		return nil, gorm.ErrRecordNotFound

	}
	result = d.db.First(&res, id)
	if result.Error != nil {
		return nil, result.Error
	}
	data, err := json.Marshal(res)
	return data, err
}

func (d *database) DeleteToDo(id string) error {
	result := d.db.Delete(&ToDo{}, id)
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected < 1 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
