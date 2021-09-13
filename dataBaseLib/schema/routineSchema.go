package schema

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

//@todo Add functionality to updated slot and represent data better

type SlotStruct struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Title     string             `json:"title"`
	Note      string             `json:"note"`
	UpdatedAt time.Time          `json:"updatedat"`
}
type NotificationStruct struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Message   string             `json:"message"`
	CreatedAt time.Time          `json:"createdAt"`
}
type RoutineModel struct {
	Id            primitive.ObjectID   `bson:"_id" json:"id"`
	Name          string               `json:"name"`
	Notification  []NotificationStruct `json:"notification"`
	FollowRequest []JWTModel           `json:"followrequest"`
	Owner         JWTModel             `json:"owner"`
	Slot          []SlotStruct         `json:"slot"`
	Access        []JWTModel           `json:"access"`
	Follower      []JWTModel           `json:"follower"`
	CreatedAt     time.Time            `json:"createdat"`
}
