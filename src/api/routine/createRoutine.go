package API

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"webServer/dataBaseLib"
	"webServer/dataBaseLib/schema"
)

func createRoutine(ctx iris.Context) {
	claims:=jwt.Get(ctx).(*schema.JWTModel)
	fmt.Println(claims)
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := dataBase.UserDB.FindOne(ct, bson.M{
		"_id": claims.Id,
	})
	var valueFromUser schema.User
	err := result.Decode(valueFromUser)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(valueFromUser)
	ctx.JSON(valueFromUser)
}
func CheckingForRoutineLimit(jwtModel schema.JWTModel) (schema.User, int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := dataBase.UserDB.FindOne(ctx, bson.M{
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
