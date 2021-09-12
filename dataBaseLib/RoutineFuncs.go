package dataBase

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
	"webServer/dataBaseLib/schema"
)
func CheckingForRoutineLimit(userId primitive.ObjectID) int {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := UserDB.FindOne(ctx, bson.M{
		"_id": userId,
	})
	var valueFromUser schema.User
	err := result.Decode(&valueFromUser)
	if err != nil {
		log.Println(err)
		return -1
	}
	return len(valueFromUser.Routines)
}

