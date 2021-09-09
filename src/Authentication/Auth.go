package Auth

import (
	"fmt"
	"github.com/kataras/iris/v12"
)

func Server(app *iris.Application) {
	fmt.Println(googleOauthConfig.RedirectURL)
	//app := iris.New()
	loginAPI := app.Party("/auth")
	{
		// loginAPI.Use(iris.Compression)
		

		// GET: http://localhost:8080/books
		loginAPI.Get("/", handleMain)
		loginAPI.Post("/", loginHandler)
		// POST: http://localhost:8080/books
		loginAPI.Get("/google", handleGoogleLogin)
		loginAPI.Get(redirectUrlGoogleOauth, handleGoogleCallback)
	}
}
