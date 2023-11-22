package models

type Product struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}
