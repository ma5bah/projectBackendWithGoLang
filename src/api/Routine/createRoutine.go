package routineAPI

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"webServer/dataBaseLib"
	"webServer/dataBaseLib/schema"
)

func createRoutine(ctx iris.Context) {
	inputData := struct {
		RoutineName string
	}{}

	err := ctx.ReadJSON(&inputData)

	if inputData.RoutineName == "" || len(inputData.RoutineName) > 50 {
		log.Println(err)
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	jwtModelData := jwt.Get(ctx).(*schema.JWTModel)
	routines := dataBase.CheckingForRoutineLimit(jwtModelData.Id)
	if routines == -1 {

		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	if routines >= 5 {

		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	ctxDB, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	firstUserArray := make([]schema.JWTModel, 1)
	firstUserArray[0] = *jwtModelData
	slotArray := make([]schema.SlotStruct, 70)
	for idx := 0; idx < 70; idx++ {
		slotArray[idx].Id = primitive.NewObjectID()
		slotArray[idx].Title = ""
		slotArray[idx].Note = ""
	}
	//@todo take routine and also add it in user profile
	RoutineInsertedID := primitive.NewObjectID()
	_, err = dataBase.RoutineDB.InsertOne(ctxDB, schema.RoutineModel{
		Id:            RoutineInsertedID,
		Notification:  make([]schema.NotificationStruct, 0, 100),
		FollowRequest: make([]schema.JWTModel, 0, 100),
		Name:          inputData.RoutineName,
		Owner:         *jwtModelData,
		Slot:          slotArray,
		Access:        firstUserArray,
		Follower:      firstUserArray,
		CreatedAt:     time.Now(),
	})
	if err != nil {
		log.Println(err)
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	//@TODO Reasigned ctxDB, See if any problem occerred
	ctxDB, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dataBase.UserDB.FindOneAndUpdate(ctxDB,
		bson.M{
			"_id": jwtModelData.Id,
		},
		bson.M{"$push": bson.M{
			"Routines": schema.RoutineStruct{
				Id:   RoutineInsertedID,
				Name: inputData.RoutineName,
			},
		}},
	)
	ctx.Writef("Successfully created a Routine")
}
func viewRoutine(ctx iris.Context) {
	inputData := struct {
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
	result := dataBase.RoutineDB.FindOne(ctxDB, bson.M{
		"_id": inputData.RoutineId,
		"follower": bson.M{
			"$elemMatch": bson.M{
				"_id": jwtModelData.Id,
			},
		},
	})
	if result.Err() != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Writef("Please Provide a correct Query or Ask permission from Owner")
		return
	}
	log.Println(inputData)
	var RoutineData schema.RoutineModel
	err = result.Decode(&RoutineData)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	ctx.JSON(RoutineData)
}
