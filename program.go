package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type Program struct {
	Id        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User
	UserId    int64
	Lifts     []Lift
}

func registerProgramRoutes(router *mux.Router) {
	db.AutoMigrate(Program{})

	router.HandleFunc("/programs", programList).Methods("GET")
	router.HandleFunc("/program/{id}", programFetch).Methods("GET")
	router.HandleFunc("/program/", programCreate).Methods("POST")

}

func programList(writer http.ResponseWriter, request *http.Request) {
	programs := make([]Program, 0)
	db.Find(&programs)

	marshalled, err := json.Marshal(programs)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Error encoding programs"))
	}

	writer.WriteHeader(200)
	writer.Write(marshalled)
}

func programFetch(writer http.ResponseWriter, request *http.Request) {
	var program Program
	vars := mux.Vars(request)

	db.Find(&program, vars["id"])

	if &program == nil {
		writer.WriteHeader(404)
		writer.Write([]byte("Program record not found."))
		return
	}

	// turn the response into JSON
	bytes, err := json.Marshal(program)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Error encoding the program"))
		return
	}

	writer.WriteHeader(200)
	writer.Write(bytes)
	return
}

func programCreate(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var program Program

	err := decoder.Decode(&program)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte("Could not decode the program"))
		return
	}

	if !db.NewRecord(program) {
		writer.WriteHeader(400)
		writer.Write([]byte("This program record already exists"))
		return
	} else {
		db.Save(&program)
		marshalled, err := json.Marshal(program)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte("Error saving the program"))
			return
		}
		writer.WriteHeader(200)
		writer.Write(marshalled)
		return
	}
}
