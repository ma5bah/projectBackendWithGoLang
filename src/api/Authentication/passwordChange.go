package userAPI

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
	"webServer/common"
	dataBase "webServer/dataBaseLib"
	"webServer/dataBaseLib/schema"
)

func changePasswordHandler(ctx iris.Context)  {
	inputData:= struct {
		Password string
		ChangedPassword string
	}{}
	err:= ctx.ReadJSON(&inputData)
	if err!=nil{
		log.Println(err)
		ctx.StatusCode(500)
		return
	}
	if inputData.Password==""||inputData.ChangedPassword==""||len(inputData.Password)>100||len(inputData.ChangedPassword)>100{
		ctx.StatusCode(400)
		return
	}
	jwtModelData := jwt.Get(ctx).(*schema.JWTModel)
	ctxDB, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//dataBase.UserDB.UpdateOne(ctxDB,bson.M{"_id":jwtModelData.Id},bson.M{"password":})
	result:=dataBase.UserDB.FindOne(ctxDB,bson.M{"email":jwtModelData.Email})
	if result.Err()!=nil{
		log.Println(result.Err())
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	var UserModelData schema.User
	err = result.Decode(&UserModelData)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	if !common.CheckPasswordHash(inputData.Password,UserModelData.Password){
		//	log.Println(common.CheckPasswordHash(inputData.Password,UserModelData.Password))
		//log.Println(UserModelData.Password)
		//log.Println(inputData.Password)
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Writef("Give correct password or generate new one by clicking google login")
		return
	}
	_, err = dataBase.GenerateNewUserPassword(UserModelData.Email,inputData.ChangedPassword)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	ctx.Writef("Successfully changed password!")
}