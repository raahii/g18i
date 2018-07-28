package main

import (
  "fmt"
  "log"
  "os"
  "time"
  "encoding/json"
  "net/http"
  "github.com/gorilla/mux"
)

type Recipe struct {
  ID int               `json:"id,omitempty"`
  Title string         `json:"title,omitempty"`
  MakingTime string    `json:"making_time,omitempty"`
  Serves string        `json:"serves,omitempty"`
  Ingredients string   `json:"ingredients,omitempty"`
  Cost int             `json:"cost,omitempty"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time  `json:"updated_at"`
}

func toTime(str string) (t time.Time){
  format := "2006-01-02 15:04:05"
  t, err := time.Parse(format, str)
  
  if err != nil {
    log.Fatal(err)
  }
  
  return
}

var recipes []Recipe

func GetRecipes(w http.ResponseWriter, r *http.Request) {
  json.NewEncoder(w).Encode(recipes)
}

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

  recipes = append(recipes, Recipe{ ID: 1, Title: "チキンカレー", MakingTime: "45分", Serves: "4人", Ingredients: "玉ねぎ,肉,スパイス", Cost: 1000, CreatedAt: toTime("2016-01-10 12:10:12"), UpdatedAt: toTime("2016-01-10 12:10:12") })
  recipes = append(recipes, Recipe{ ID: 2, Title: "オムライス", MakingTime: "30分", Serves: "2人", Ingredients: "玉ねぎ,卵,スパイス,醤油", Cost: 700, CreatedAt: toTime("2016-01-11 13:10:12"), UpdatedAt: toTime("2016-01-11 13:10:12") })

  port := os.Getenv("PORT")
  if port == "" {
    log.Fatal("$PORT must be set")
  }
  log.Fatal(http.ListenAndServe(":8000", router))
}

func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello, %q", r.URL.Path[1:])
}
