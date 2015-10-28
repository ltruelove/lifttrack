package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm"
	"net/http"
)

type Lift struct {
	Id         int64
	ProgramId  int64
	LiftType   LiftType
	LiftTypeId int64
	Weight     float64 `json:",string"`
	Sets       float64 `json:",string"`
	Reps       float64 `json:",string"`
}

func registerLiftRoutes(router *mux.Router) {
	db.AutoMigrate(&Lift{})

	router.HandleFunc("/lifts", liftList).Methods("GET")
	router.HandleFunc("/lift/{id}", liftFetch).Methods("GET")
	router.HandleFunc("/lift/", liftCreate).Methods("POST")

}

func liftList(writer http.ResponseWriter, request *http.Request) {
	lifts := make([]Lift, 0)
	db.Find(&lifts)

	marshalled, err := json.Marshal(lifts)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Error encoding lifts"))
	}

	writer.WriteHeader(200)
	writer.Write(marshalled)
}

func liftFetch(writer http.ResponseWriter, request *http.Request) {
	_, err := validateToken(writer, request)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte(err.Error()))
		return
	}

	var lift Lift
	vars := mux.Vars(request)

	//pass a 0 to get a blank record
	// it's probably a good idea to add the ProgramId on the client side
	if vars["id"] == "0" {
		var encodedLift []byte
		encodedLift, err = json.Marshal(lift)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte("Internal error"))
			return
		}

		writer.WriteHeader(200)
		writer.Write(encodedLift)
		return
	}

	db.Find(&lift, vars["id"])

	if &lift == nil {
		writer.WriteHeader(404)
		writer.Write([]byte("Lift record not found."))
		return
	}

	// turn the response into JSON
	var bytes []byte
	bytes, err = json.Marshal(lift)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Error encoding the lift"))
		return
	}

	writer.WriteHeader(200)
	writer.Write(bytes)
	return
}

func liftCreate(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var lift Lift

	err := decoder.Decode(&lift)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte("Could not decode the lift"))
		return
	}

	if !db.NewRecord(lift) {
		writer.WriteHeader(400)
		writer.Write([]byte("This lift record already exists"))
		return
	} else {
		db.Save(&lift)
		marshalled, err := json.Marshal(lift)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte("Error saving the lift"))
			return
		}
		writer.WriteHeader(200)
		writer.Write(marshalled)
		return
	}
}
