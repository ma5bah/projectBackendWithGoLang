package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

type JWTModel struct {
	Id    primitive.ObjectID `bson:"_id" json:"id"`
	Email string             `json:"email"`
	Name  string             `json:"name"`
}
