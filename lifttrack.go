package main

import (
	//"bytes"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

var db gorm.DB

const PW_SALT string = "StuffAndThings"

func main() {
	var err error
	db, err = gorm.Open("mysql", "lifttrack:123@/lifttrack?charset=utf8&parseTime=True")

	if err != nil {
		fmt.Printf("%s\r\n", err)
	}

	db.DB()

	//configure the database

	port := flag.Int("port", 9081, "port to serve on")
	dir := flag.String("directory", "web/", "directory of web files")
	flag.Parse()

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)

	// setup routes
	router := mux.NewRouter()

	registerUserRoutes(router)
	registerLiftTypeRoutes(router)
	registerLiftRoutes(router)
	registerProgramRoutes(router)

	router.PathPrefix("/").Handler(fileHandler)
	http.Handle("/", router)

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Listening on %d\r\n", *port)

	// this call blocks -- the progam runs here forever
	err = http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}
