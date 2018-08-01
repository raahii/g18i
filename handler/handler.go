package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	. "github.com/raahii/g18i/model"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type Handler struct {
	db  *gorm.DB
	val *validator.Validate
}

func NewHandler(db *gorm.DB) Handler {
	val := validator.New()
	return Handler{db: db, val: val}
}

// handlers
func (c Handler) GetRecipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var recipes []Recipe
	err := c.db.Find(&recipes).Error

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

func (c Handler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var recipe Recipe
	params := mux.Vars(r)

	err := c.db.First(&recipe, params["id"]).Error
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

func recipeJsonTag(fieldName string) string {
	var recipe Recipe
	t := reflect.TypeOf(recipe)
	field, _ := t.FieldByName(fieldName)
	return field.Tag.Get("json")
}

func (c Handler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var recipe Recipe

	// parse post params
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response := map[string]string{
			"message": "Recipe creation failed!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// decode json, set params
	err = json.Unmarshal(body, &recipe)
	fmt.Println("recipe:", recipe.Title, recipe.Cost)
	if err != nil {
		response := map[string]string{
			"message": "Recipe creation failed!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// validation
	if err := c.val.Struct(recipe); err != nil {
		var errorFields []string
		for _, e := range err.(validator.ValidationErrors) {
			errorFields = append(errorFields, recipeJsonTag(e.Field()))
		}

		response := map[string]string{
			"message":  "Recipe creation failed!",
			"required": strings.Join(errorFields, ", ")}
		json.NewEncoder(w).Encode(response)
		return
	}

	// create
	if err := c.db.Create(&recipe).Error; err != nil {
		response := map[string]interface{}{
			"message": "Recipe creation failed!",
			"error":   err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"message": "Recipe successfully created!",
		"recipe":  recipe}
	json.NewEncoder(w).Encode(response)
}

func (c Handler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var recipe Recipe
	r.ParseForm()

	// find
	params := mux.Vars(r)
	if err := c.db.First(&recipe, params["id"]).Error; err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		response := map[string]interface{}{"message": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// parse post params
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response := map[string]string{
			"message": "Recipe creation failed!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// decode json, set params
	err = json.Unmarshal(body, &recipe)
	fmt.Println("recipe:", recipe.Title, recipe.Cost)
	if err != nil {
		response := map[string]string{
			"message": "Recipe creation failed!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// validation
	if err := c.val.Struct(recipe); err != nil {
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

	// update
	if err := c.db.Save(&recipe).Error; err != nil {
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

func (c Handler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	// find
	var recipe Recipe
	params := mux.Vars(r)

	if err := c.db.First(&recipe, params["id"]).Error; err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		response := map[string]interface{}{"message": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// delete
	if err := c.db.Delete(&recipe).Error; err != nil {
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
