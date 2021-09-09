//https://docs.iris-go.com/iris/security/jwt

package common

import (
	"github.com/kataras/iris/v12/context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"webServer/dataBaseLib/schema"
)

var (
	jwtSecret = LocalGetEnv("jwtSecret")
	jwtEncDecSecret=LocalGetEnv("jwtEncDecSecret")[:32]
)

func GenerateToken(userId primitive.ObjectID,userName,userEmail string) ([]byte,error)  {
	jwtModelData:=schema.JWTModel{
		Id: userId,
		Name: userName,
		Email: userEmail,
	}
	genSigner := jwt.NewSigner(jwt.HS256, jwtSecret, 60*time.Minute)
	genSigner.WithEncryption([]byte(jwtEncDecSecret), nil)
	return genSigner.Sign(jwtModelData)
}
func VerifyJwtMiddleware() context.Handler {
	verifier := jwt.NewVerifier(jwt.HS256, jwtSecret)
	// Enable server-side token block feature (even before its expiration time):
	verifier.WithDefaultBlocklist()
	// Enable payload decryption with:
	verifier.WithDecryption([]byte(jwtEncDecSecret), nil)
	verifyMiddleware := verifier.Verify(func() interface{} {
		return new(schema.JWTModel)
	})
	return verifyMiddleware
}

func Logout(ctx iris.Context) {
	err := ctx.Logout()
	if err != nil {
		ctx.WriteString(err.Error())
	} else {
		ctx.Writef("Logged Out")
	}
}


//func protected(ctx iris.Context) {
//
//	ctx.Write([]byte("masbah"))
//	Get the verified and decoded claims.
//	claims := jwt.Get(ctx).(*schema.JWTModel)
//
//	// Optionally, get token information if you want to work with them.
//	// Just an example on how you can retrieve
//	// all the standard claims (set by signer's max age, "exp").
//	//res:=jwt.GetVerifiedToken(ctx)
//	//fmt.Println(claims.Foo,res)
//	//ctx.Writef("%s,",claims.Foo,res)
//	standardClaims := jwt.GetVerifiedToken(ctx).StandardClaims
//	fmt.Println("Passed")
//	expiresAtString := standardClaims.ExpiresAt().
//		Format(ctx.Application().ConfigurationReadOnly().GetTimeFormat())
//	fmt.Println("Passed")
//
//	timeLeft := standardClaims.Timeleft()
//	fmt.Println("Passed")
//
//	fmt.Println(claims)
//	ctx.Writef("ID=%s\nexpires at: %s\ntime left: %s\n", claims.Id, expiresAtString, timeLeft)
//}


//func InitJwt(app *iris.Application) {
//
//	//signer := jwt.NewSigner(jwt.HS256, secret, 60*time.Minute)
//	// Enable payload encryption with:
//	//signer.WithEncryption([]byte("itsa16bytesecret"), nil)
//	//app.Get("/auth/jwt", generateToken(signer))
//
//	//verifier := jwt.NewVerifier(jwt.HS256, secret)
//	//// Enable server-side token block feature (even before its expiration time):
//	//verifier.WithDefaultBlocklist()
//	//// Enable payload decryption with:
//	//// verifier.WithDecryption(encKey, nil)
//	//verifyMiddleware := verifier.Verify(func() interface{} {
//	//	return new(schema.JWTModel)
//	//})
//
//	protectedAPI := app.Party("/protected")
//	// Register the verify middleware to allow access only to authorized clients.
//	// protectedAPI.Use(verifyMiddleware)
//	// ^ or UseRouter(verifyMiddleware) to disallow unauthorized http error handlers too.
//	protectedAPI.Get("/jn", VerifyJwtMiddleware(), dummyHandler)
//	protectedAPI.Get("/", protected)
//	// Invalidate the token through server-side, even if it's not expired yet.
//	protectedAPI.Get("/logout", logout)
//
//	app.Listen(":8421")
//}
