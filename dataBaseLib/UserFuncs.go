package dataBase

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"webServer/common"
	"webServer/dataBaseLib/schema"
)

func FindUserByEmail(oAuthData schema.OAuthDataModel) *mongo.SingleResult {
	emailData := oAuthData.Email
	//fmt.Println(emailData)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := UserDB.FindOne(ctx, bson.D{
		{"email", emailData},
	})

	return result
}
func CreateUser(Email, Name, ProfilePicture string) (string, error) {
	var (
		emailData   = Email
		nameData    = Name
		pictureData = ProfilePicture
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	passwordData, err := common.GeneratePassword(30)
	if err != nil {
		return "", err
	}
	hashedPassword, err := common.HashPassword(passwordData)
	if err != nil {
		return "", err
	}
	_, err = UserDB.InsertOne(ctx, schema.User{
		Id:             primitive.NewObjectID(),
		Email:          emailData,
		Name:           nameData,
		ProfilePicture: pictureData,
		Password:       hashedPassword,
		JoinedAt:       time.Now(),
		Routines:       []schema.RoutineStruct{},
		Bio:            "Hey there! I am using Notio",
	})
	if err != nil {
		return "", err
	}
	return passwordData, nil
}
func SyncUser(Email, Name, ProfilePicture string) error {
	var (
		emailData   = Email
		nameData    = Name
		pictureData = ProfilePicture
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updatedResult := UserDB.FindOneAndUpdate(ctx, bson.M{"email": emailData},bson.M{"$set": bson.M{"name": nameData, "profilePicture": pictureData}})
	//updatedResult := UserDB.FindOneAndUpdate(ctx, bson.M{"email": emailData}, bson.M{"name": nameData, "profilePicture": pictureData})

	if updatedResult.Err() != nil {
		return updatedResult.Err()
	}
	return nil
}
func GenerateNewUserPassword(emailData, passwordData string) (string, error) {
	if passwordData == "" {
		generatedPasswordData, err := common.GeneratePassword(30)
		if err != nil {
			return "", err
		}
		passwordData = generatedPasswordData
	}
	hashedPasswordData, err := common.HashPassword(passwordData)
	if err != nil {
		return "", err
	}
	ctxDB, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = UserDB.UpdateOne(ctxDB, bson.M{"email": emailData}, bson.M{
		"$set": bson.M{"password": hashedPasswordData},
	})
	if err != nil {
		return "", err
	}
	return passwordData, nil
}

func AddRoutineToUser(userData *schema.JWTModel, routineData *schema.RoutineStruct) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	updatedResult := UserDB.FindOneAndUpdate(ctx, bson.M{
		"_id": userData.Id,
	}, bson.M{"$addToSet": bson.M{"Routines": schema.RoutineStruct{Id: routineData.Id, Name: routineData.Name}}})
	return updatedResult
}
func AddUserToRoutine(userData *schema.JWTModel, routineData *schema.RoutineStruct) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	updatedResult := RoutineDB.FindOneAndUpdate(ctx, bson.M{
		"_id": routineData.Id,
	}, bson.M{"$addToSet": bson.M{"Routines": schema.JWTModel{Id: userData.Id, Email: userData.Email, Name: userData.Name}}})
	return updatedResult
}
func RemoveRoutineFromUser(userId primitive.ObjectID, RoutineId primitive.ObjectID) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	updatedResult := UserDB.FindOneAndUpdate(ctx, bson.M{
		"_id":      userId,
		"Routines": bson.M{"$elemMatch": bson.M{"_id": RoutineId}},
	}, bson.M{"$pull": bson.M{"Routines": bson.M{"_id": RoutineId}}})
	return updatedResult
}
func RemoveUserFromRoutine(userId primitive.ObjectID, RoutineId primitive.ObjectID) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	updatedResult := RoutineDB.FindOneAndUpdate(ctx,
		bson.M{
			"_id":      RoutineId,
			"follower": bson.M{"$elemMatch": bson.M{"_id": userId}},
			"owner":    bson.M{"$not": bson.M{"_id": userId}},
		},
		bson.M{
			"$pull": bson.M{"follower": bson.M{"_id": RoutineId, "access": bson.M{"_id": RoutineId}}},
		},
	)
	return updatedResult
}
