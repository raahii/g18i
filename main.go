package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// models
type Recipe struct {
	ID          int    `json:"id"`
	Title       string `json:"title" validate:"required"`
	MakingTime  string `json:"making_time" validate:"required"`
	Serves      string `json:"serves" validate:"required"`
	Ingredients string `json:"ingredients" validate:"required"`
	Cost        int    `json:"cost" validate:"required"`
}

var validate *validator.Validate

// handlers
func GetRecipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var recipes []Recipe
	err := db.Find(&recipes).Error

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		response := map[string]interface{}{"message": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"recipes": recipes}
	json.NewEncoder(w).Encode(response)
}

func GetRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var recipe Recipe
	params := mux.Vars(r)

	err := db.First(&recipe, params["id"]).Error
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		response := map[string]interface{}{"message": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": "Recipe details by id",
		"recipe":  []Recipe{recipe}}
	json.NewEncoder(w).Encode(response)
}

func setRecipeParams(r *http.Request, recipe *Recipe) error {
	// title
	if v := r.FormValue("title"); v != "" {
		recipe.Title = v
	}

	// making_time
	if v := r.FormValue("making_time"); v != "" {
		recipe.MakingTime = v
	}

	// serves
	if v := r.FormValue("serves"); v != "" {
		recipe.Serves = v
	}

	// ingredients
	if v := r.FormValue("ingredients"); v != "" {
		recipe.Ingredients = v
	}

	// cost
	if v := r.FormValue("cost"); v != "" {
		cost, err := strconv.Atoi(v)
		if err != nil {
			return errors.New("Cost must be an integer.")
		}

		recipe.Cost = cost
	}

	return nil
}

func recipeJsonTag(fieldName string) string {
	var recipe Recipe
	t := reflect.TypeOf(recipe)
	field, _ := t.FieldByName(fieldName)
	return field.Tag.Get("json")
}

func CreateRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var recipe Recipe
	r.ParseForm()

	// set params
	if err := setRecipeParams(r, &recipe); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		response := map[string]string{
			"message": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// validation
	if err := validate.Struct(recipe); err != nil {
		var errorFields []string
		for _, e := range err.(validator.ValidationErrors) {
			errorFields = append(errorFields, recipeJsonTag(e.Field()))
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		response := map[string]string{
			"message":  "Recipe creation failed!",
			"required": strings.Join(errorFields, ", ")}
		json.NewEncoder(w).Encode(response)
		return
	}

	// create
	if err := db.Create(&recipe).Error; err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		response := map[string]interface{}{
			"message": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": "Recipe successfully created!",
		"recipe":  recipe}
	json.NewEncoder(w).Encode(response)
}

func UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var recipe Recipe
	r.ParseForm()

	// find
	params := mux.Vars(r)
	if err := db.First(&recipe, params["id"]).Error; err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		response := map[string]interface{}{"message": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// set params
	if err := setRecipeParams(r, &recipe); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		response := map[string]string{
			"message": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// update
	if err := db.Save(&recipe).Error; err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		response := map[string]interface{}{"message": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": "Recipe successfully updated!",
		"recipe":  []Recipe{recipe}}
	json.NewEncoder(w).Encode(response)
}

func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	// find
	var recipe Recipe
	params := mux.Vars(r)

	if err := db.First(&recipe, params["id"]).Error; err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		response := map[string]interface{}{"message": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := db.Delete(&recipe).Error; err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		response := map[string]interface{}{"message": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": "Recipe successfully removed!"}
	json.NewEncoder(w).Encode(response)
}

// db
var db *gorm.DB

func gormConnect() *gorm.DB {
	DBMS := "mysql"
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASSWORD")
	PROTOCOL := os.Getenv("DB_PROTOCOL")
	DBNAME := os.Getenv("DB_NAME")
	ENABLE_DB_LOG := os.Getenv("ENABLE_DB_LOG") == "1"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?parseTime=true"
	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		panic(err.Error())
	}

	db.LogMode(ENABLE_DB_LOG)

	return db
}

// main
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/recipes", GetRecipes).Methods("GET")
	router.HandleFunc("/recipes/{id}", GetRecipe).Methods("GET")
	router.HandleFunc("/recipes", CreateRecipe).Methods("POST")
	router.HandleFunc("/recipes/{id}", UpdateRecipe).Methods("PATCH")
	router.HandleFunc("/recipes/{id}", DeleteRecipe).Methods("DELETE")

	db = gormConnect()
	defer db.Close()

	validate = validator.New()

	log.Fatal(http.ListenAndServe(":8000", router))
}
