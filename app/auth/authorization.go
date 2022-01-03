package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
)

type DBHandler struct {
	Coll *mongo.Collection
}

type AuthorizationCodeURL struct {
	URL string `bson:"url" json:"url"`
}

func Auth(c echo.Context) error {
	// if アクセストークンがDBにあればNULLをかえす
	// else アクセストークンがなければURLをかえして認可、アクセストークンを取得
	fmt.Println("REDIRECT_URL=", os.Getenv("REDIRECT_URL"))
	conf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		//Scopes:       []string{"SCOPE"},
		RedirectURL: os.Getenv("REDIRECT_URL"), // TODO: コールバックを指定して渡す
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.secure.freee.co.jp/public_api/authorize",
			TokenURL: "https://accounts.secure.freee.co.jp/public_api/token"},
	}
	//state := fmt.Sprintf("st%d", time.Now().UnixNano())
	//fmt.Println(state)
	//url := conf.AuthCodeURL(state) // 認可ページのURL
	url := conf.AuthCodeURL(os.Getenv("OAUTH_STATE")) // 認可ページのURL
	fmt.Println(url)
	// Authorization Code を受け取る
	return c.JSON(http.StatusOK, AuthorizationCodeURL{URL: url})
}
