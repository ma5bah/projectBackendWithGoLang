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

//@TODO add request in user profile also
func followRequest(ctx iris.Context) {
	inputData := struct {
		RoutineID primitive.ObjectID
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

	updatedResult := dataBase.RoutineDB.FindOneAndUpdate(ctxDB,
		bson.M{
			"_id": inputData.RoutineID,
		},
		bson.M{
			"$addToSet": bson.M{
				"followrequest": jwtModelData,
			},
		},
	)
	if updatedResult.Err() != nil {
		log.Println(updatedResult.Err())
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Writef("Please give correct query.")
		return
	}
	var routineData schema.RoutineModel
	err = updatedResult.Decode(&routineData)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	updatedResult=dataBase.AddRoutineToUser(jwtModelData, &schema.RoutineStruct{Id: routineData.Id, Name: routineData.Name})
	if updatedResult.Err() != nil {
		log.Println(updatedResult.Err())
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	ctx.Writef("Successfully handled follow request")
}

//@TODO delete follow request from routine accesser
//@todo delete follow request from user
func acceptFollowRequest(ctx iris.Context) {
	inputData := struct {
		UserInfo  schema.JWTModel
		RoutineID primitive.ObjectID
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
		"_id":           inputData.RoutineID,
		"access":        bson.M{"$elemMatch": bson.M{"_id": jwtModelData.Id}},
		"followrequest": bson.M{"$elemMatch": bson.M{"_id": inputData.UserInfo.Id}},
	}, bson.M{
		"$pull":     bson.M{"followrequest": inputData.UserInfo},
		"$addToSet": bson.M{"follower": inputData.UserInfo},
	})
	log.Println(updatedResult)
	if updatedResult.Err() != nil {
		log.Println(updatedResult.Err())
		ctx.StatusCode(iris.StatusBadRequest)

		ctx.Writef("Please give correct query or Ask owner to give permission to modify Routine")
		return
	}
	ctx.Writef("Successfully updated follower")
}
func deleteFollowRequestFromRoutine(ctx iris.Context)  {
	inputData := struct {
		UserInfo  schema.JWTModel
		RoutineID primitive.ObjectID
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
		"_id":           inputData.RoutineID,
		"access":        bson.M{"$elemMatch": bson.M{"_id": jwtModelData.Id}},
		"followrequest": bson.M{"$elemMatch": bson.M{"_id": inputData.UserInfo.Id}},
	}, bson.M{
		"$pull":     bson.M{"followrequest": inputData.UserInfo},
	})
	if updatedResult.Err()!=nil{
		log.Println(updatedResult.Err())
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	updatedResult=dataBase.RemoveRoutineFromUser(inputData.UserInfo.Id,inputData.RoutineID)
	if updatedResult.Err()!=nil{
		log.Println(updatedResult.Err())
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
}
func deleteFollowRequestFromUser(ctx iris.Context)  {
	inputData := struct {
		RoutineInfo schema.RoutineStruct
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
	updatedResult:=dataBase.RemoveRoutineFromUser(jwtModelData.Id,inputData.RoutineInfo.Id)
	if updatedResult.Err()!=nil{
		log.Println(updatedResult.Err())
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	updatedResult = dataBase.RoutineDB.FindOneAndUpdate(ctxDB, bson.M{
		"_id":           inputData.RoutineInfo.Id,
		"followrequest": bson.M{"$elemMatch": bson.M{"_id": jwtModelData.Id}},
	}, bson.M{
		"$pull":     bson.M{"followrequest": jwtModelData},
	})
	if updatedResult.Err()!=nil{
		log.Println(updatedResult.Err())
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
}
