package userAPI

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	dataBase "webServer/dataBaseLib"
	"webServer/dataBaseLib/schema"
)

func viewUserProfile(ctx iris.Context) {
	jwtModelData := jwt.Get(ctx).(*schema.JWTModel)
	ctxDB, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result:=dataBase.UserDB.FindOne(ctxDB,bson.M{"_id":jwtModelData.Id})
	if result.Err()!=nil{
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	var UserDataDB struct {
		Id             primitive.ObjectID `bson:"_id" json:"id"`
		Email          string             `json:"email"`
		Name           string             `json:"name"`
		ProfilePicture string             `json:"profilePicture"`
		Routines       []schema.RoutineStruct    `json:"Routines"`
		JoinedAt       time.Time          `json:"joinedAt"`
	}
	err := result.Decode(&UserDataDB)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	//UserDataDB=delete(UserDataDB,"Password")
	ctx.JSON(UserDataDB)
}
//func editProfile(ctx iris.Context)  {
//	Name           string             `json:"name"`
//	ProfilePicture string             `json:"profilePicture"`
//}