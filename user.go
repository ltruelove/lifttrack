package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm"
	"net/http"
)

type User struct {
	Id       int64
	Username string
	Password string
}

func registerUserRoutes(router *mux.Router) {
	db.AutoMigrate(User{})

	router.HandleFunc("/users", userList).Methods("GET")
	router.HandleFunc("/user/{id}", userFetch).Methods("GET")
	router.HandleFunc("/user/", userCreate).Methods("POST")

}

func userList(writer http.ResponseWriter, request *http.Request) {
	users := make([]User, 0)
	db.Find(&users)

	marshalled, err := json.Marshal(users)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Error encoding users"))
	}

	writer.WriteHeader(200)
	writer.Write(marshalled)
}

func userFetch(writer http.ResponseWriter, request *http.Request) {
	var user User
	vars := mux.Vars(request)

	db.Find(&user, vars["id"])

	if &user == nil {
		writer.WriteHeader(404)
		writer.Write([]byte("User record not found."))
		return
	}

	// turn the response into JSON
	bytes, err := json.Marshal(user)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Error encoding the user"))
		return
	}

	writer.WriteHeader(200)
	writer.Write(bytes)
	return
}

func userCreate(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var user User

	err := decoder.Decode(&user)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte("Could not decode the user"))
		return
	}

	if !db.NewRecord(user) {
		writer.WriteHeader(400)
		writer.Write([]byte("This user record already exists"))
		return
	} else {
		db.Save(&user)
		marshalled, err := json.Marshal(user)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte("Error saving the user"))
			return
		}
		writer.WriteHeader(200)
		writer.Write(marshalled)
		return
	}
}
