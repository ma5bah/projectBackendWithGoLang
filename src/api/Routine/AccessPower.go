package routineAPI

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

func upgradeToAccess(ctx iris.Context) {
	inputData := struct {
		UserInfo  schema.JWTModel
		RoutineId primitive.ObjectID
	}{}
	err := ctx.ReadJSON(&inputData)
	log.Println(inputData)
	if err != nil {
		log.Println(err)
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	jwtModelData := jwt.Get(ctx).(*schema.JWTModel)
	ctxDB, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	updatedResult := dataBase.RoutineDB.FindOneAndUpdate(ctxDB, bson.M{
		"_id":   inputData.RoutineId,
		"owner": jwtModelData,
		//"$elemMatch": bson.M{"follower": inputData.UserInfo},
		"follower": bson.M{"$elemMatch": bson.M{"_id": inputData.UserInfo.Id, "email": inputData.UserInfo.Email, "name": inputData.UserInfo.Name}},
		"access":   bson.M{"$nin": bson.A{bson.M{"_id": inputData.UserInfo.Id, "email": inputData.UserInfo.Email, "name": inputData.UserInfo.Name}}},
	}, bson.M{
		"$addToSet": bson.M{"access": inputData.UserInfo},
	})
	if updatedResult.Err() != nil {
		log.Println(updatedResult.Err())
		ctx.StatusCode(iris.StatusBadRequest)

		ctx.Writef("Please give correct query or Ask owner to give permission to modify Routine")
		return
	}
	ctx.Writef("Successfully updated accesser")
}
func degradeFromAccess(ctx iris.Context) {
	inputData := struct {
		UserInfo  schema.JWTModel
		RoutineId primitive.ObjectID
	}{}
	err := ctx.ReadJSON(&inputData)

	if err != nil {
		log.Println(err)
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	jwtModelData := jwt.Get(ctx).(*schema.JWTModel)
	ctxDB, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	updatedResult := dataBase.RoutineDB.FindOneAndUpdate(ctxDB, bson.M{
		"_id":    inputData.RoutineId,
		"owner":  jwtModelData,
		"access": bson.M{"$elemMatch": bson.M{"_id": inputData.UserInfo.Id, "email": inputData.UserInfo.Email, "name": inputData.UserInfo.Name}},
	}, bson.M{
		//"$pull": bson.M{"$elemMatch": bson.M{"_id": inputData.UserInfo.Id}},
		"$pull": bson.M{"access": inputData.UserInfo},
	})
	if updatedResult.Err() != nil {
		log.Println(updatedResult.Err())
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Writef("Please give correct query or Ask owner to give permission to modify Routine")
		return
	}
	x, _ := updatedResult.DecodeBytes()
	ctx.Writef("Successfully updated accesser %s", x)
}
func removeUserFromRoutine(ctx iris.Context) {
	inputData := struct {
		UserInfo  schema.JWTModel
		RoutineId primitive.ObjectID
	}{}
	err := ctx.ReadJSON(&inputData)
	log.Println(inputData)
	if err != nil {
		log.Println(err)
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	jwtModelData := jwt.Get(ctx).(*schema.JWTModel)
	ctxDB, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	updatedResult := dataBase.RoutineDB.FindOneAndUpdate(ctxDB, bson.M{
		"_id": inputData.RoutineId,
		"access": bson.M{
			"$elemMatch": bson.M{
				"_id": jwtModelData.Id,
			},
			"$nin": bson.A{bson.M{
				"_id": inputData.UserInfo.Id,
			}},
		},
	}, bson.M{
		"$pull": bson.M{"follower": inputData.UserInfo},
	})
	if updatedResult.Err() != nil {
		log.Println(updatedResult.Err())
		ctx.StatusCode(iris.StatusBadRequest)

		ctx.Writef("Please give correct query or Ask owner to give permission to modify Routine")
		return
	}
	ctx.Writef("Successfully updated Follower")
}
