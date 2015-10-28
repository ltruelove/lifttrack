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
	db.AutoMigrate(&LiftType{})

	/*** Create the initial lift types if they don't already exist ***/

	lt := LiftType{}
	db.Where("name = ?", "Squat").First(&lt)

	if db.NewRecord(lt) {
		lt.Name = "Squat"
		lt.Description = "Standard squat."
		db.Save(&lt)
	}

	lt = LiftType{}
	db.Where("name = ?", "Bench Press").First(&lt)

	if db.NewRecord(lt) {
		lt.Name = "Bench Press"
		lt.Description = "Standard bench press."
		db.Save(&lt)
	}

	lt = LiftType{}
	db.Where("name = ?", "Dead Lift").First(&lt)

	if db.NewRecord(lt) {
		lt.Name = "Dead Lift"
		lt.Description = "Standard dead lift."
		db.Save(&lt)
	}

	lt = LiftType{}
	db.Where("name = ?", "Press").First(&lt)

	if db.NewRecord(lt) {
		lt.Name = "Press"
		lt.Description = "Standing press."
		db.Save(&lt)
	}

	lt = LiftType{}
	db.Where("name = ?", "Power Clean").First(&lt)

	if db.NewRecord(lt) {
		lt.Name = "Power Clean"
		lt.Description = "Standard power clean."
		db.Save(&lt)
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
