package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/raahii/g18i/handler"
	"log"
	"net/http"
	"os"
)

// db
var db *gorm.DB

func GormConnect() *gorm.DB {
	DBMS := os.Getenv("DB_KIND")
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASSWORD")
	PROTOCOL := os.Getenv("DB_PROTOCOL")
	DBNAME := os.Getenv("DB_NAME")
	ENABLE_DB_LOG := os.Getenv("ENABLE_DB_LOG") == "1"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME
	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		fmt.Println(">> Failed to connect to the database.")
		panic(err.Error())
	} else {
		fmt.Println(">> Connected to the database.")
	}

	db.LogMode(ENABLE_DB_LOG)

	return db
}

// main
func main() {
	db = GormConnect()
	defer db.Close()

	h := handler.NewHandler(db)

	router := mux.NewRouter()
	router.HandleFunc("/recipes", h.GetRecipes).Methods("GET")
	router.HandleFunc("/recipes/{id}", h.GetRecipe).Methods("GET")
	router.HandleFunc("/recipes", h.CreateRecipe).Methods("POST")
	router.HandleFunc("/recipes/{id}", h.UpdateRecipe).Methods("PATCH")
	router.HandleFunc("/recipes/{id}", h.DeleteRecipe).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Fatal(http.ListenAndServe(":"+port, router))
}
