package model

// models
type Recipe struct {
	ID          int    `json:"id"`
	Title       string `json:"title" validate:"required"`
	MakingTime  string `json:"making_time" validate:"required"`
	Serves      string `json:"serves" validate:"required"`
	Ingredients string `json:"ingredients" validate:"required"`
	Cost        int    `json:"cost" validate:"required"`
}
