package dataBase

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"webServer/common"
	"webServer/dataBaseLib/schema"
)

func FindUserByEmail(emailData string) *mongo.SingleResult {
	//fmt.Println(emailData)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := UserDB.FindOne(ctx, bson.D{
		{"email", emailData},
	})
	fmt.Println("err",result.Err())
	return result
}
func CreateUser(Email, Name, ProfilePicture string) (*mongo.InsertOneResult,string, error) {
	var (
		emailData   = Email
		nameData    = Name
		pictureData = ProfilePicture
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	passwordData, err := common.PasswordGenerator(30)
	if err != nil {
		return nil,"", err
	}
	result, err := UserDB.InsertOne(ctx, schema.User{
		Id:             primitive.NewObjectID(),
		Email:          emailData,
		Name:           nameData,
		ProfilePicture: pictureData,
		Password:       passwordData,
		JoinedAt:       time.Now(),
		Routines:       []schema.RoutineStruct{},
	})
	if err != nil {
		return nil,"", err
	}
	return result,passwordData, nil
}
func CheckingForRoutineLimit(jwtModel schema.JWTModel) (schema.User, int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := UserDB.FindOne(ctx, bson.M{
		"_id": jwtModel.Id,
	})
	var valueFromUser schema.User
	err := result.Decode(valueFromUser)
	if err != nil {
		return schema.User{}, 500
	}
	if !(len(valueFromUser.Routines) < 5) {
		return schema.User{}, 400
	}
	return valueFromUser, 0
}
func AddRoutineToUser(userData *schema.User, routineData *schema.RoutineModel) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//updatedRoutineData := append(userData.Routines, schema.RoutineStruct{Id: routineData.Id})
	result, _ :=UserDB.UpdateOne(ctx, bson.M{
		"_id": userData.Id,
	}, bson.M{"$push": bson.M{"Routines": schema.RoutineStruct{Id: routineData.Id}}})
	fmt.Println(result)
	//bson.M{"$push":bson.M{"user.$.sales":bson.M{"$each":addsales}}}

	//result := UserDB.FindOne(ctx, bson.M{
	//	"_id": userID,
	//})
	//return result
}
