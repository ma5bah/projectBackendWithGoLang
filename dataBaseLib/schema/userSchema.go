package schema

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

//@Todo connect with Routine schema

type RoutineStruct struct {
	Id primitive.ObjectID `bson:"_id" json:"id"`
	Name          string               `json:"name"`
}
type User struct {
	Id             primitive.ObjectID `bson:"_id" json:"id"`
	Email          string             `json:"email"`
	Name           string             `json:"name"`
	ProfilePicture string             `json:"profilePicture"`
	Routines       []RoutineStruct    `json:"Routines"`
	Password       string             `json:"password"`
	JoinedAt       time.Time          `json:"joinedAt"`
	Bio string `json:"username"`
}
type OAuthDataModel struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Picture       string `json:"picture"`
	Id            string `json:"id"`
	VerifiedEmail bool   `json:"verified_email"`
}

/*
email
password
profilePicture
fullname
active
username
bio
joinedAt
Routines
*/
