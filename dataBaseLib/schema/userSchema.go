package schema

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

//@Todo connect with routine schema
//type User struct {
//	ID primitive.ObjectID
//	Email string
//	Password string
//	ProfilePicture string
//	FullName string
//	Active string
//	UserName string
//	JoinedAt string
//	Routines []string
//}
type RoutineStruct struct {
	Id primitive.ObjectID `bson:"_id" json:"id"`
}
type User struct {
	Id             primitive.ObjectID `bson:"_id" json:"id"`
	Email string                      `json:"email"`
	Name           string             `json:"name"`
	ProfilePicture string             `json:"profilePicture"`
	Routines       []RoutineStruct    `json:"Routines"`
	//Username string `json:"username"`
	Password string `json:"password"`
	JoinedAt time.Time `json:"joinedAt"`
}
type jwtData struct {
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
