package main

import (
  "fmt"
  "log"
  "os"
  "net/http"
  "github.com/gorilla/mux"
)

func GetRecipes(w http.ResponseWriter, r *http.Request) {}
func GetRecipe(w http.ResponseWriter, r *http.Request) {}
func CreateRecipe(w http.ResponseWriter, r *http.Request) {}
func UpdateRecipe(w http.ResponseWriter, r *http.Request) {}
func DeleteRecipe(w http.ResponseWriter, r *http.Request) {}

func main() {
  router := mux.NewRouter()
  router.HandleFunc("/recipes", GetRecipes).Methods("GET")
  router.HandleFunc("/recipes/{id}", GetRecipe).Methods("GET")
  router.HandleFunc("/recipes", CreateRecipe).Methods("POST")
  router.HandleFunc("/recipes/{id}", UpdateRecipe).Methods("POST")
  router.HandleFunc("/recipes/{id}", DeleteRecipe).Methods("DELETE")

  port := os.Getenv("PORT")
  if port == "" {
    log.Fatal("$PORT must be set")
  }
  log.Fatal(http.ListenAndServe(":8000", router))
}

func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello, %q", r.URL.Path[1:])
}
