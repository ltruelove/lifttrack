package main

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.crypto/pbkdf2"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type User struct {
	Id       int64
	Username string
	Password string `json:"-"`
	First    string
	Last     string
}

type DisplayUser struct {
	Username string
	First    string
	Last     string
}

type AccessToken struct {
	Token        string
	UserId       int64
	LastAccessed time.Time
	User         User
}

func registerUserRoutes(router *mux.Router) {
	db.AutoMigrate(User{})
	db.AutoMigrate(AccessToken{})

	router.HandleFunc("/users", userList).Methods("GET")
	router.HandleFunc("/user/{id}", userFetch).Methods("GET")
	router.HandleFunc("/user/", userCreate).Methods("POST")
	router.HandleFunc("/user/login", userLogin).Methods("POST")

}

func userList(writer http.ResponseWriter, request *http.Request) {
	_, err := validateToken(request)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte(err.Error()))
		return
	}

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
	_, err := validateToken(request)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte(err.Error()))
		return
	}

	var user User
	vars := mux.Vars(request)

	//return a blank user
	if vars["id"] == "0" {
		var encodedUser []byte
		encodedUser, err = json.Marshal(user)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte("Error encoding the user"))
			return
		}

		writer.WriteHeader(200)
		writer.Write(encodedUser)
		return
	}

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
	_, err := validateToken(request)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte(err.Error()))
		return
	}

	decoder := json.NewDecoder(request.Body)
	var user User

	err = decoder.Decode(&user)
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
		var marshalled []byte
		marshalled, err = json.Marshal(user)
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

func userLogin(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var user User

	err := decoder.Decode(&user)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte("Could not decode the user"))
		return
	}

	user.EncryptPassword()
	var result User

	db.Where("username = ?", user.Username).
		Where("password = ?", user.Password).
		First(&result)

	if result.Id == 0 {
		writer.WriteHeader(401)
		writer.Write([]byte("User name and password combination does not exist"))
		return
	}

	displayUser := DisplayUser{Username: result.Username, First: result.First, Last: result.Last}
	encodedUser, err := json.Marshal(displayUser)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Internal error"))
		return
	}

	accessToken := getToken(result.Id)
	writer.Header().Set("Token", accessToken.Token)
	writer.WriteHeader(200)
	writer.Write(encodedUser)

}

func (u *User) EncryptPassword() {
	//add encryption routine here
	salt := []byte(PW_SALT)
	u.Password = string(HashPassword([]byte(u.Password), salt))
}

func HashPassword(password, salt []byte) []byte {
	defer clear(password)
	return pbkdf2.Key(password, salt, 4096, sha256.Size, sha256.New)
}

func clear(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
}

func getToken(userId int64) AccessToken {
	accessToken := AccessToken{"", userId, time.Now(), User{}}
	accessToken.Token = uuid.New()

	//remove other tokens for this user if they exist
	db.Where("user_id = ?", userId).Delete(AccessToken{})

	db.Save(&accessToken)

	return accessToken
}

func validateToken(req *http.Request) (*AccessToken, error) {
	var accessToken AccessToken
	tokenText := req.Header.Get("Token")

	db.Where("token = ?", tokenText).First(&accessToken)

	if accessToken.Token == "" {
		return nil, fmt.Errorf("Token not found")
	}

	current := time.Now()
	if accessToken.LastAccessed.Before(current) {
		dif := current.Sub(accessToken.LastAccessed).Minutes()
		if dif > 15 {
			return nil, fmt.Errorf("Token has expired")
		} else {
			accessToken.LastAccessed = time.Now()
			db.Save(&accessToken)
		}
	} else {
		return nil, fmt.Errorf("Token has expired")
	}

	return &accessToken, nil
}
