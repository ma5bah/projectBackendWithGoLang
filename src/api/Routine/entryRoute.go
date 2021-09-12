package routineAPI

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"webServer/common"
	dataBase "webServer/dataBaseLib"
	"webServer/dataBaseLib/schema"
)

//@todo user routine unsubscribe or owner delete
//@todo if owner delete routine then remove the routine from every follower
//@todo user routine unsubscribe when he is accessor....discession pending


func RoutineEntryRoute(app *iris.Application) {
	Routine := app.Party("/routine")
	{
		Routine.Use(common.VerifyJwtMiddleware())
		Routine.Post("/", viewRoutine)
		Routine.Post("/create", createRoutine)
		Routine.Post("/removeroutinefromuser", removeRoutineFromUser)

		Routine.Post("/upgradetoaccess", upgradeToAccess)
		Routine.Post("/degradefromaccess", degradeFromAccess)
		Routine.Post("/updateNotification", updateNotification)
		Routine.Post("/deleteNotification", deleteNotification)
		Routine.Post("/updateslot", updateOrDeleteSlot)
		Routine.Post("/followrequest", followRequest)
		Routine.Post("/removeuser", removeUserFromRoutine)
		Routine.Post("/acceptfollowrequest", acceptFollowRequest)
		Routine.Post("/followrequest", followRequest)
		Routine.Post("/deletefollowrequestfromuser", deleteFollowRequestFromUser)
		Routine.Post("/deletefollowrequestfromroutine", deleteFollowRequestFromRoutine)
	}
}

func removeRoutineFromUser(ctx iris.Context)  {
	inputData := struct {
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
	updatedResult:=	dataBase.RemoveRoutineFromUser(jwtModelData.Id,inputData.RoutineId)
	if updatedResult.Err()!=nil{
		log.Println(updatedResult.Err())
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	ctx.Writef("Success")
}