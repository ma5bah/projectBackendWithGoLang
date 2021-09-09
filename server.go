package main

import (
	"github.com/kataras/iris/v12"
	"log"
	"webServer/src/Authentication"
	"webServer/src/api/routine"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}



// albums slice to seed record album data.
// func homeHandler(ctx *gin.Context) {
// 	ctx.IndentedJSON(http.StatusOK, albums)
// }

func Server() {
	app := iris.New()
	Auth.Server(app)
	API.RoutineEntryRoute(app)
	booksAPI := app.Party("/api")
	{
		booksAPI.Use(iris.Compression)

		// GET: http://localhost:8080/books
		booksAPI.Get("/", list)
		// POST: http://localhost:8080/books
		booksAPI.Post("/", create)
	}

	app.Listen(":8421")

}

// Book example.
type Book struct {
	Title string `json:"title"`
}

func list(ctx iris.Context) {
	books := []Book{
		{"Mastering Concurrency in Go"},
		{"Go Design Patterns"},
		{"Black Hat Go"},
	}

	_, err := ctx.JSON(books)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	// TIP: negotiate the response between server's prioritizes
	// and client's requirements, instead of ctx.JSON:
	// ctx.Negotiation().JSON().MsgPack().Protobuf()
	// ctx.Negotiate(books)
}

func create(ctx iris.Context) {
	var b Book
	err := ctx.ReadJSON(&b)
	// TIP: use ctx.ReadBody(&b) to bind
	// any type of incoming data instead.
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Book creation failure").DetailErr(err))
		// TIP: use ctx.StopWithError(code, err) when only
		// plain text responses are expected on errors.
		return
	}

	println("Received Book: " + b.Title)

	ctx.StatusCode(iris.StatusCreated)
}
