package userAPI

import (
	"github.com/kataras/iris/v12"
	"webServer/common"
)

func UserEntryRoute( app *iris.Application)  {
	User:=app.Party("/user")
	{
		User.Use(common.VerifyJwtMiddleware())

		User.Post("/",viewUserProfile)
		User.Post("/changepassword",changePasswordHandler)
		User.Post("/updatebio",updateBioHandler)
		//User.Post("/syncprofile",syncProfile)

		//User.Post("/updateslot", updateOrDeleteSlot)


	}


	Auth := app.Party("/auth")
	{
		Auth.Use(iris.Compression)
		Auth.Post("/", loginHandler)
		Auth.Get("/google", handleGoogleLogin)
		Auth.Get(redirectUrlGoogleOauth, handleGoogleCallback)
	}
}
