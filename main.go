package main

import (
	"./database"
	"./model/message"
	"./model/post"
	"./model/user"
	"./route"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

/*
	@brief : initDatabase is initialising our data base
*/
func initDatabase() {
	var err error                                                    // create a variable named err to handle errors
	database.DBConn, err = gorm.Open("sqlite3", "database/books.db") // open our database and return an error if there is one
	if err != nil {                                                  // will inform that the database encountered an error
		panic("Failed to connect to database")
	}
	fmt.Println("Database connection successfully opened")
	database.DBConn.AutoMigrate(&user.User{}, &post.Post{}, &message.Message{}) // "automigrate" will sett our database structure (creating an user with the paremeter already in user)
	fmt.Println("Database Migrated")
}

/**
@brief: main will call our function to run the program
*/
func main() {
	app := mux.NewRouter()                       // create a new server
	initDatabase()                               //call our function who initialize the data base
	defer database.DBConn.Close()                //at the end defer(defer will be the last instruction in our program) will close our database
	route.IniRoute(app)                          // set our routes
	log.Fatal(http.ListenAndServe(":8080", app)) // start our server on port 8080
}
