package userAPI

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"net/http"
	"webServer/common"
	dataBase "webServer/dataBaseLib"
	"webServer/dataBaseLib/schema"
)

var (
	googleOauthConfig *oauth2.Config
	// TODO: randomize it
	oauthStateString = "pseudo-random"
)
var scopeGoogleOauth = []string{
	"https://www.googleapis.com/auth/plus.me",
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/userinfo.profile",
}
var redirectUrlGoogleOauth = common.LocalGetEnv("redirectUrlGoogleOauth")
var hostAddr = common.LocalGetEnv("hostAddr")
var hostPort = common.LocalGetEnv("hostPort")
func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  hostAddr+hostPort+"/auth" + redirectUrlGoogleOauth,
		ClientID:     common.LocalGetEnv("ClientID"),
		ClientSecret: common.LocalGetEnv("GoAuthSecret"),
		Scopes:       scopeGoogleOauth,
		Endpoint:     google.Endpoint,
	}
}

func handleGoogleLogin(ctx iris.Context) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	ctx.Redirect(url, iris.StatusTemporaryRedirect)
}

func handleGoogleCallback(ctx iris.Context) {
	//log.Println(ctx.Request())
	content, err := getUserInfo(ctx.FormValue("state"), ctx.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		ctx.Redirect("/", http.StatusTemporaryRedirect)
		return
	}
	var contentJson schema.OAuthDataModel
	err = json.Unmarshal(content, &contentJson)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	if !contentJson.VerifiedEmail {
		ctx.JSON("Please Verify your email first")
		return
	}
	//@TODO -> before creating the user,I must include _id in userInfo
	resultDB := dataBase.FindUserByEmail(contentJson)

	sentData := struct {
		Password string
	}{}
	if resultDB.Err() != nil {
		passwordData, err := dataBase.CreateUser(contentJson.Email, contentJson.Name, contentJson.Picture)
		if err != nil {

			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}
		sentData.Password = passwordData
	} else {
		err:=dataBase.SyncUser(contentJson.Email, contentJson.Name, contentJson.Picture)
		if err != nil {
			log.Println(err)
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}
		passwordData, err := dataBase.GenerateNewUserPassword(contentJson.Email, "")
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}
		sentData.Password = passwordData
	}
	ctx.JSON(sentData)
}

func getUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	ctx := context.TODO()
	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	//fmt.Println(string(contents))
	return contents, nil
}