    package main

    import (
        "encoding/gob"
        "code.google.com/p/go-uuid/uuid"
        "crypto/sha1"
        "encoding/json"
        "fmt"
        _ "github.com/go-sql-driver/mysql"
        "github.com/gorilla/mux"
        "github.com/gorilla/sessions"
        _ "github.com/jinzhu/gorm"
        "net/http"
        "time"
    )

    type User struct {
        Id       int64
        Username string `sql:"not null;unique"`
        Password string
        First    string
        Last     string
        Programs []Program
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

    var store = sessions.NewCookieStore([]byte("should-be-a-config"))

    func registerUserRoutes(router *mux.Router) {
        db.AutoMigrate(&User{})
        db.AutoMigrate(&AccessToken{})
        gob.Register(&AccessToken{})

        router.HandleFunc("/users", userList).Methods("GET")
        router.HandleFunc("/user/programs", userPrograms).Methods("GET")
        router.HandleFunc("/user/{id}", userFetch).Methods("GET")
        router.HandleFunc("/user", userCreate).Methods("POST")
        router.HandleFunc("/user/login", userLogin).Methods("POST")

    }

    func userList(writer http.ResponseWriter, request *http.Request) {
        _, err := validateToken(writer, request)
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

func userPrograms(writer http.ResponseWriter, request *http.Request) {

	token, err := validateToken(writer, request)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte(err.Error()))
		return
	}

    var user User
	db.Preload("Programs").Find(&user, token.UserId)

	if &user == nil {
		writer.WriteHeader(404)
		writer.Write([]byte("User record not found."))
		return
	}

	// turn the response into JSON
	var bytes []byte
	bytes, err = json.Marshal(user.Programs)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Error encoding the programs"))
		return
	}

	writer.WriteHeader(200)
	writer.Write(bytes)
	return
}

func userFetch(writer http.ResponseWriter, request *http.Request) {
	var user User
	vars := mux.Vars(request)

	//return a blank user
	if vars["id"] == "0" {
		encodedUser, err := json.Marshal(user)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte("Error encoding the user"))
			return
		}

		writer.WriteHeader(200)
		writer.Write(encodedUser)
		return
	}

	//validate after we test for fetching a blank record
	_, err := validateToken(writer, request)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte(err.Error()))
		return
	}

	db.Find(&user, vars["id"])

	if &user == nil {
		writer.WriteHeader(404)
		writer.Write([]byte("User record not found."))
		return
	}

	//blank out the password
	user.Password = ""
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

		var existing User
		db.Where("username = ?", user.Username).First(&existing)

		if !db.NewRecord(existing) {
			writer.WriteHeader(400)
			writer.Write([]byte("Username already in use"))
			return
		}

		user.EncryptPassword()

		db.Save(&user)
		var marshalled []byte
		//blank out the password
		//user.Password = ""
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

	db.Where("username = ?", user.Username).First(&result)

	if result.Id == 0 {
		writer.WriteHeader(401)
		writer.Write([]byte("User name and password combination does not exist"))
		return
	} else if user.Password != result.Password {
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

    session, _ := store.Get(request, "lifttrack-userToken")
    session.Values["token"] = accessToken
    session.Save(request, writer)

	writer.WriteHeader(200)
	writer.Write(encodedUser)

}

func (u *User) EncryptPassword() {
	//add encryption routine here
	salted := []byte(fmt.Sprintf("%s%s", PW_SALT, u.Password))
	u.Password = fmt.Sprintf("%x", sha1.Sum(salted))

	fmt.Println(u.Password)
}

func getToken(userId int64) AccessToken {
    loc, err := time.LoadLocation("UTC")
    if err != nil {
        fmt.Println("err: ", err.Error())
    }

	accessToken := AccessToken{"", userId, time.Now().In(loc), User{}}
	accessToken.Token = uuid.New()

	//remove other tokens for this user if they exist
	db.Where("user_id = ?", userId).Delete(AccessToken{})

	db.Save(&accessToken)

	return accessToken
}

func validateToken(writer http.ResponseWriter, req *http.Request) (*AccessToken, error) {
	var accessToken AccessToken
    var tokenText = ""
    session, _ := store.Get(req, "lifttrack-userToken")

    // Retrieve our struct and type-assert it
    val := session.Values["token"]
    var userToken = &AccessToken{}

    userToken, _ = val.(*AccessToken)
    fmt.Print("Token given: ")
    fmt.Println(userToken.Token)

    db.Where("token = ?", userToken.Token).First(&accessToken)

    if accessToken.Token == "" {
        return nil, fmt.Errorf("Token not found")
    }

    if accessToken.LastAccessed.Before(time.Now()) {
        dif := time.Since(accessToken.LastAccessed).Minutes()
        fmt.Print("minutes difference ")
        fmt.Println(dif)
        if dif > 15 {
            return nil, fmt.Errorf("Token has expired")
        } else {
            loc, err := time.LoadLocation("UTC")
            if err != nil {
                fmt.Println("err: ", err.Error())
            }
            
            accessToken.LastAccessed = time.Now().In(loc)
            db.Save(&accessToken)
        }
    } else {
        return nil, fmt.Errorf("Token has expired")
    }

    session.Values["token"] = tokenText
    session.Save(req, writer)

	return &accessToken, nil
}
