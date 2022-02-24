package controllers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"kintai-bot/app/common"
	"kintai-bot/app/domain"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type freeeBody struct {
	CompanyID int       `json:"company_id"`
	Type      string    `json:"type"`
	BaseDate  string    `json:"base_date"`
	Datetime  time.Time `json:"datetime"`
}

type AuthorizationCodeURL struct {
	URL string `bson:"url" json:"url"`
}

func (controller *TokenController) Auth(c echo.Context) error {
	// if アクセストークンがDBにあればNULLをかえす
	//filter, err := strconv.Atoi(os.Getenv("COMPANY_ID"))
	//if err != nil {
	//	return c.JSON(http.StatusBadRequest, common.NewErrorResponse(400, "No Parameters"))
	//}
	//token, err := controller.Interactor.TokenByCompanyID(filter)
	//if (token.AccessToken != "") && (token.RefreshToken != "") {
	//	return c.JSON(http.StatusOK, common.NewErrorResponse(200, "Access token is already set"))
	//}
	// else アクセストークンがなければURLをかえして認可、アクセストークンを取得
	conf := &oauth2.Config{
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
	conf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		//Scopes:       []string{"SCOPE"},
		RedirectURL: os.Getenv("REDIRECT_URL"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.secure.freee.co.jp/public_api/authorize",
			TokenURL: "https://accounts.secure.freee.co.jp/public_api/token"},
	}
	oauth2Token, err := conf.Exchange(ctx, c.QueryParam("code"))
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusBadRequest, common.NewErrorResponse(400, "Failed to get the access token."))
	}

	// AccessToken と RefreshToken を保存
	companyID, err := strconv.Atoi(os.Getenv("COMPANY_ID"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorResponse(500, "No Parameters"))
	}
	t := domain.Token{CompanyID: companyID, AccessToken: oauth2Token.AccessToken, RefreshToken: oauth2Token.RefreshToken, Expiry: oauth2Token.Expiry}
	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorResponse(500, "Invalid Request"))
	}
	if err := controller.Interactor.Update(t); err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorResponse(500, "Could not update record"))
	}
	return c.JSON(http.StatusCreated, "Authorization successful.")
}

func (controller *TokenController) Dakoku(c echo.Context) error {
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		//Scopes:       []string{"SCOPE"},
		RedirectURL: os.Getenv("REDIRECT_URL"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.secure.freee.co.jp/public_api/authorize",
			TokenURL: "https://accounts.secure.freee.co.jp/public_api/token"},
	}
	filter, err := strconv.Atoi(os.Getenv("COMPANY_ID"))
	tokenDomain, err := controller.Interactor.TokenByCompanyID(filter)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewErrorResponse(404, "Not Found"))
	}
	freee_id := c.Param("freee_id")
	if freee_id == "" {
		return c.JSON(http.StatusBadRequest, common.NewErrorResponse(400, "No Parameters"))
	}
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewErrorResponse(400, "Invalid Request"))
	}

	// アクセストークンの期限が切れていればリフレッシュ
	token := oauth2.Token{AccessToken: tokenDomain.AccessToken, TokenType: "Bearer", RefreshToken: tokenDomain.RefreshToken, Expiry: tokenDomain.Expiry}
	tokenSource := conf.TokenSource(ctx, &token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorResponse(500, "New Tokens are invalid"))
	}
	client := oauth2.NewClient(ctx, tokenSource)
	url := fmt.Sprintf("https://api.freee.co.jp/hr/api/v1/employees/%s/time_clocks", freee_id)
	body := bytes.NewBuffer([]byte(b))
	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	defer func() {
		if newToken.Valid() { // AccessToken と RefreshToken を保存
			t := domain.Token{CompanyID: filter, AccessToken: newToken.AccessToken, RefreshToken: newToken.RefreshToken, Expiry: newToken.Expiry}
			if t.AccessToken == "" || t.RefreshToken == "" {
				fmt.Println(t, "Returned AccessToken is empty")
			}
			if err := controller.Interactor.Update(t); err != nil {
				fmt.Println(t, "Cloud not update record")
			}
		}
	}()
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorResponse(500, "Failed to refresh the access token."))
	} else if resp.StatusCode != (200 | 201) {
		return c.JSON(http.StatusInternalServerError, common.NewErrorResponse(500, "Failed to refresh the access token."))
	}

	return c.JSON(http.StatusOK, token)
}
