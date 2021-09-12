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

func updateNotification(ctx iris.Context) {
	//@TODO check if object id needs to be changed by input

	inputData := struct {
		RoutineID primitive.ObjectID
		Message   string
	}{}
	err := ctx.ReadJSON(&inputData)
	log.Println(inputData.RoutineID)
	if err != nil {
		log.Println(err)
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	if inputData.Message == "" || len(inputData.Message) > 200 {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	jwtModelData := jwt.Get(ctx).(*schema.JWTModel)
	ctxDB, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := dataBase.RoutineDB.FindOne(ctxDB, bson.M{
		"_id":    inputData.RoutineID,
		"access": bson.M{"$elemMatch": bson.M{"_id": jwtModelData.Id}},
	})
	if result.Err() != nil {
		log.Println(result.Err())
		ctx.StatusCode(iris.StatusBadRequest)

		ctx.Writef("Please give correct query or Ask owner to give permission to modify Routine")
		return
	}
	var RoutineModelData schema.RoutineModel
	err = result.Decode(&RoutineModelData)
	if err != nil {
		log.Println(err)
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	if len(RoutineModelData.Notification) >= 100 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Writef("Your limit has been exceeded.Please Delete Some Notification")
		return
	}
	newCtxDB, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = dataBase.RoutineDB.UpdateOne(newCtxDB,
		bson.M{
			"_id":    inputData.RoutineID,
			"access": bson.M{"$elemMatch": bson.M{"_id": jwtModelData.Id}},
		},
		bson.M{
			"$push": bson.M{
				"notification": schema.NotificationStruct{Id: primitive.NewObjectID(), Message: inputData.Message, CreatedAt: time.Now()},
			},
		},
	)
	if err != nil {
		log.Println(err)
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Writef("Please give correct query or Ask owner to give permission to modify Routine")
		return
	}
	ctx.Writef("Successfully updated Notification")
}

//@TODO delete notification
func deleteNotification(ctx iris.Context) {
	inputData := struct {
		RoutineID      primitive.ObjectID
		NotificationID primitive.ObjectID
	}{}
	err := ctx.ReadJSON(&inputData)
	log.Println(inputData.RoutineID)
	if err != nil {
		log.Println(err)
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	jwtModelData := jwt.Get(ctx).(*schema.JWTModel)

	ctxDB, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = dataBase.RoutineDB.UpdateOne(ctxDB, bson.M{
		"_id":    inputData.RoutineID,
		"access": bson.M{"$elemMatch": bson.M{"_id": jwtModelData.Id}},
	}, bson.M{
		"$pull": bson.M{"notification": bson.M{"_id": inputData.NotificationID}},
	})
	if err != nil {
		log.Println(err)
		ctx.StatusCode(iris.StatusBadRequest)

		ctx.Writef("Please give correct query or Ask owner to give permission to modify Routine")
		return
	}
	ctx.Writef("Successfully updated the notification.")
}
