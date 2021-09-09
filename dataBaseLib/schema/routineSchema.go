package schema

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type RoutineModel struct {
	Id             primitive.ObjectID `bson:"_id" json:"id"`
	Notification  []string      `json:"notification"`
	FollowRequest []interface{} `json:"followRequest"`
	Name          string        `json:"name"`
	//Owner         struct {
	//	Id    primitive.ObjectID `bson:"_id" json:"id"`
	//	//Username string `json:"username"`
	//} `json:"owner"`
	Owner JWTModel `json:"owner"`
	Slot  []struct {
		Id    primitive.ObjectID `bson:"_id" json:"id"`
		Title string             `json:"title"`
		Note  string             `json:"note"`
	} `json:"slot"`
	Access JWTModel `json:"access"`
	Follower JWTModel `json:"follower"`
	CreatedAt struct {
		Date time.Time `json:"$date"`
	} `json:"createdAt"`
}
