package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

type TestDataType struct {
	Id    primitive.ObjectID `bson:"_id" json:"id"`
	Title string             `json:"title"`
	Array []string           `json:"array"`
}

