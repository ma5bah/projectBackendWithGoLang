package routineAPI

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
	dataBase "webServer/dataBaseLib"
	"webServer/dataBaseLib/schema"
)

func updateOrDeleteSlot(ctx iris.Context) {
	//@TODO filter malicious input data

	inputData := struct {
		RoutineID primitive.ObjectID
		SlotID    primitive.ObjectID
		Title     string
		Message   string
	}{}
	err := ctx.ReadJSON(&inputData)
	//log.Println(inputData)
	if err != nil {
		log.Println(err)
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	if len(inputData.Message) > 50 || len(inputData.Message) > 200 {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	jwtModelData := jwt.Get(ctx).(*schema.JWTModel)
	ctxDB, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//@todo learn array operation in mongodb
	//@todo https://stackoverflow.com/questions/56306164/mongo-go-how-to-use-arrayfilter-to-find-elem-in-array-of-objects-inside-array-o
	updateResult := dataBase.RoutineDB.FindOneAndUpdate(ctxDB,
		bson.M{
			"_id":    inputData.RoutineID,
			"access": bson.M{"$elemMatch": bson.M{"_id": jwtModelData.Id}},
			//"access.$._id": jwtModelData.Id,
		},
		bson.M{
			"$set": bson.M{"slot.$[element].title": inputData.Title, "slot.$[element].note": inputData.Message,"slot.$[element].updatedat":time.Now()},
		},
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{Filters: []interface{}{bson.M{"element._id": inputData.SlotID}}}),
	)

	if updateResult.Err() != nil {
		//log.Println("ERROR occurred", updateResult.Err())
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Writef("Please give correct query or Ask owner to give permission to modify Routine")
		return
	}

	ctx.Writef("Successfully updated Slot")
}
