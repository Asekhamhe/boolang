package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// BookResult model
type BookResult struct {
	ID string `json:"id,omitempty" bson:"_id"`
	// bson:inline flattens anonymous field for mongo driver
	Book `bson:"inline"`
}

// Book model for creating new book
type Book struct {
	ID     primitive.ObjectID `json:"-"`
	BookID string             `json:"book_id" bson:"book_id"`
	Title  string             `json:"title,omitempty" bson:"title"`
	Author string             `json:"author,omitempty" bson:"author"`
	Year   string             `json:"year,omitempty" bson:"year"`
	// One to One relationship
	// Permission Permission `json:"permission" bson:"permission"`
	// One to Many relationship
	// Permission []Permission `json:"permission" bson:"permission"`
}
