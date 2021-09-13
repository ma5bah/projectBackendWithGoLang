package userAPI

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
	dataBase "webServer/dataBaseLib"
	"webServer/dataBaseLib/schema"
)

func viewUserProfile(ctx iris.Context) {
	jwtModelData := jwt.Get(ctx).(*schema.JWTModel)
	ctxDB, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := dataBase.UserDB.FindOne(ctxDB, bson.M{"_id": jwtModelData.Id})
	if result.Err() != nil {
		log.Println(result.Err())
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	var UserDataDB struct {
		Id             primitive.ObjectID     `bson:"_id" json:"id"`
		Email          string                 `json:"email"`
		Name           string                 `json:"name"`
		ProfilePicture string                 `json:"profilePicture"`
		Routines       []schema.RoutineStruct `json:"Routines"`
		JoinedAt       time.Time              `json:"joinedAt"`
		Bio            string                 `json:"bio"`
	}
	err := result.Decode(&UserDataDB)
	if err != nil {
		log.Println(err)
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	ctx.JSON(UserDataDB)
}

func updateBioHandler(ctx iris.Context) {
	inputData := struct {
		Bio string
	}{}
	err := ctx.ReadJSON(&inputData)
	log.Println(inputData)
	if err != nil {
		log.Println(err)
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	if len(inputData.Bio) > 256 {
		ctx.Writef("bio can't exceed 256 character")
		return
	}
	jwtModelData := jwt.Get(ctx).(*schema.JWTModel)
	ctxDB, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	updatedResult := dataBase.UserDB.FindOneAndUpdate(ctxDB, bson.M{"_id": jwtModelData.Id}, bson.M{"$set": bson.M{"bio": inputData.Bio}})
	if updatedResult.Err() != nil {
		log.Println(updatedResult.Err())
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	ctx.Writef("Successfully updated Bio")
}
