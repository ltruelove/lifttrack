package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm"
	"net/http"
)

type LiftType struct {
	Id          int64
	Name        string
	Description string
}

func registerLiftTypeRoutes(router *mux.Router) {
	db.AutoMigrate(LiftType{})

	/*** Create the initial lift types if they don't already exist ***/
	squat := LiftType{Id: 0, Name: "Squat", Description: "Standard squat."}
	if db.NewRecord(squat) {
		db.Save(&squat)
	}

	bench := LiftType{Id: 0, Name: "Bench Press", Description: "Standard bench press."}
	if db.NewRecord(bench) {
		db.Save(&bench)
	}

	deadlift := LiftType{Id: 0, Name: "Dead Lift", Description: "Standard dead lift."}
	if db.NewRecord(deadlift) {
		db.Save(&deadlift)
	}

	press := LiftType{Id: 0, Name: "Press", Description: "Standing press."}
	if db.NewRecord(press) {
		db.Save(&press)
	}

	powerclean := LiftType{Id: 0, Name: "Power Clean", Description: "Standard power clean."}
	if db.NewRecord(powerclean) {
		db.Save(&powerclean)
	}

	router.HandleFunc("/liftTypes", liftTypeList).Methods("GET")
	router.HandleFunc("/liftType/{id}", liftTypeFetch).Methods("GET")
	router.HandleFunc("/liftType/", liftTypeCreate).Methods("POST")

}

func liftTypeList(writer http.ResponseWriter, request *http.Request) {
	liftTypes := make([]LiftType, 0)
	db.Find(&liftTypes)

	marshalled, err := json.Marshal(liftTypes)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Error encoding liftTypes"))
	}

	writer.WriteHeader(200)
	writer.Write(marshalled)
}

func liftTypeFetch(writer http.ResponseWriter, request *http.Request) {
	var liftType LiftType
	vars := mux.Vars(request)

	db.Find(&liftType, vars["id"])

	if &liftType == nil {
		writer.WriteHeader(404)
		writer.Write([]byte("LiftType record not found."))
		return
	}

	// turn the response into JSON
	bytes, err := json.Marshal(liftType)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Error encoding the liftType"))
		return
	}

	writer.WriteHeader(200)
	writer.Write(bytes)
	return
}

func liftTypeCreate(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var liftType LiftType

	err := decoder.Decode(&liftType)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte("Could not decode the liftType"))
		return
	}

	if !db.NewRecord(liftType) {
		writer.WriteHeader(400)
		writer.Write([]byte("This liftType record already exists"))
		return
	} else {
		db.Save(&liftType)
		marshalled, err := json.Marshal(liftType)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte("Error saving the liftType"))
			return
		}
		writer.WriteHeader(200)
		writer.Write(marshalled)
		return
	}
}
