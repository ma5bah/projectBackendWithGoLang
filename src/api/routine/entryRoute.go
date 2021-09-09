package API

import (
	"github.com/kataras/iris/v12"
	"webServer/common"
)

func RoutineEntryRoute( app *iris.Application)  {
	routine:=app.Party("/routine")
	{
		routine.Use(common.VerifyJwtMiddleware())
		routine.Post("/create",createRoutine)
	}
}
