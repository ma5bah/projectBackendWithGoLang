package userAPI

import (
	"context"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
	"webServer/common"
	dataBase "webServer/dataBaseLib"
	"webServer/dataBaseLib/schema"
)

func loginHandler(ctx iris.Context)  {
	inputData:= struct {
		Email string
		Password string
	}{}

	err:= ctx.ReadJSON(&inputData)
	if err!=nil{
		log.Println(err)
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	if !(inputData.Email!=""&&common.EmailVerifier(inputData.Email)){
		ctx.StatusCode(400)
		return
	}
	if inputData.Password==""{
		log.Println(err)
		ctx.StatusCode(400)
		return
	}
	ctxDB, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result:=dataBase.UserDB.FindOne(ctxDB,bson.M{"email":inputData.Email})
	if result.Err()!=nil{
		log.Println(err)
		ctx.StatusCode(401)
		return
	}
	var inputDataDB schema.User
	err = result.Decode(&inputDataDB)
	if err != nil {
		log.Println(err)
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	if !common.CheckPasswordHash(inputData.Password,inputDataDB.Password){
		log.Println(err)
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Writef("Please provide correct info.")
		return
	}
	token, err := common.GenerateToken(inputDataDB.Id, inputDataDB.Name, inputDataDB.Email)
	if err != nil {
		log.Println(err)
		ctx.StatusCode(500)
		return
	}
	sentData := struct {
		Token string
	}{}
	sentData.Token= string(token)
	ctx.JSON(sentData)
}
