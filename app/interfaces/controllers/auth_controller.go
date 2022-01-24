package controllers

import (
	"errors"
	"fmt"
	"kintai-bot/app/common"
	"kintai-bot/app/domain"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type Company struct {
	ID          int
	Name        string
	Role        string
	ExternalCID string
	EmployeeID  int
	DisplayName string
}

type User struct {
	ID        int
	Companies []Company
}

func contains(companies []Company, company_name string) (index int, err error) {
	for i, v := range companies {
		if v.Name == company_name {
			return i, nil
		}
	}
	return 0, errors.New("Not Found The Target Company")
}

type AuthorizationCodeURL struct {
	URL string `bson:"url" json:"url"`
}

var conf *oauth2.Config
var ctx context.Context

func (controller *TokenController) Auth(c echo.Context) error {
	// if アクセストークンがDBにあればNULLをかえす
	//filter, err := strconv.Atoi(os.Getenv("COMPANY_ID"))
	//if err != nil {
	//	return c.JSON(http.StatusBadRequest, common.NewErrorResponse(400, "No Parameters"))
	//}
	//token, err := controller.Interactor.TokenByCompanyID(filter)
	//if (token.AccessToken != "") && (token.RefreshToken != "") {
	//	fmt.Println("Access Token is already set")
	//	return c.JSON(http.StatusOK, common.NewErrorResponse(200, "Access token is already set"))
	//}
	// else アクセストークンがなければURLをかえして認可、アクセストークンを取得
	conf = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		//Scopes:       []string{"SCOPE"},
		RedirectURL: os.Getenv("REDIRECT_URL"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.secure.freee.co.jp/public_api/authorize",
			TokenURL: "https://accounts.secure.freee.co.jp/public_api/token"},
	}
	//state := fmt.Sprintf("st%d", time.Now().UnixNano())
	url := conf.AuthCodeURL(os.Getenv("OAUTH_STATE")) // 認可ページのURL
	// Authorization Code を受け取る
	return c.JSON(http.StatusOK, AuthorizationCodeURL{URL: url})
}

func (controller *TokenController) AuthCallback(c echo.Context) error {
	// トークンエンドポイントにリクエスト、アクセストークン・リフレッシュトークンを取得
	ctx := context.Background()
	oauth2Token, err := conf.Exchange(ctx, c.QueryParam("code"))
	if err != nil {
		log.Fatal(err)
		return err
	}

	// AcessToken と RefreshToken を保存
	companyID, err := strconv.Atoi(os.Getenv("COMPANY_ID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewErrorResponse(400, "No Parameters"))
	}
	t := domain.Token{CompanyID: companyID, AccessToken: oauth2Token.AccessToken, RefreshToken: oauth2Token.RefreshToken}
	fmt.Println(t)
	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewErrorResponse(400, "Invalid Request"))
	}
	if err := controller.Interactor.Update(t); err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorResponse(500, "Could not update record"))
	}
	//return c.JSON(http.StatusCreated, t)
	return c.JSON(http.StatusCreated, "Authorization successful.")
}
