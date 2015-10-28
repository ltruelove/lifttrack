package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm"
	"net/http"
	"time"
    "fmt"
)

type Program struct {
	Id        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User `json:"-"`
	UserId    int64
	Lifts     []Lift
}

func registerProgramRoutes(router *mux.Router) {
	db.AutoMigrate(&Program{})

	router.HandleFunc("/programs", programList).Methods("GET")
	router.HandleFunc("/program/{id}", programFetch).Methods("GET")
	router.HandleFunc("/program", programCreate).Methods("POST")
	router.HandleFunc("/program", programUpdate).Methods("PUT")
	router.HandleFunc("/user/programs", listByUser).Methods("GET")
}

func programList(writer http.ResponseWriter, request *http.Request) {
	_, err := validateToken(writer, request)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte(err.Error()))
		return
	}

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
	accessToken, err := validateToken(writer, request)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte(err.Error()))
		return
	}

	var program Program
	vars := mux.Vars(request)

	//pass a 0 in if you want a blank program record
	if vars["id"] == "0" {
		program.UserId = accessToken.UserId
		program.Lifts = make([]Lift, 1)

		blankProgram, err := json.Marshal(program)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte("Internal error"))
			return
		}

		writer.WriteHeader(200)
		writer.Write(blankProgram)
		return
	}

	db.Preload("LiftType").Preload("Lifts").Find(&program, vars["id"])

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
	_, err := validateToken(writer, request)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte(err.Error()))
		return
	}

	decoder := json.NewDecoder(request.Body)
	var program Program

	err = decoder.Decode(&program)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(fmt.Sprintf("%s: %s","Could not decode the program",err)))
		return
	}

	if !db.NewRecord(program) {
		writer.WriteHeader(400)
		writer.Write([]byte("This program record already exists"))
		return
	}

	db.Save(&program)
	var marshalled []byte
	marshalled, err = json.Marshal(program)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Error saving the program"))
		return
	}
	writer.WriteHeader(200)
	writer.Write(marshalled)
	return
}

func programUpdate(writer http.ResponseWriter, request *http.Request) {
	_, err := validateToken(writer, request)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte(err.Error()))
		return
	}

	decoder := json.NewDecoder(request.Body)
	var program Program

	err = decoder.Decode(&program)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte("Could not decode the program"))
		return
	}

	if db.NewRecord(program) {
		writer.WriteHeader(400)
		writer.Write([]byte("This program does not exist"))
		return
	}

	db.Save(&program)
	var marshalled []byte
	marshalled, err = json.Marshal(program)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Error saving the program"))
		return
	}
	writer.WriteHeader(200)
	writer.Write(marshalled)
	return
}
