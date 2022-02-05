package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

type Company struct {
	ID          int
	Name        string
	Role        string
	ExternalCID string
	EmployeeID  int
	DisplayName string
}

var conf *oauth2.Config

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
		return c.JSON(http.StatusBadRequest, common.NewErrorResponse(400, "Failed to get the access token."))
	}

	// AcessToken と RefreshToken を保存
	companyID, err := strconv.Atoi(os.Getenv("COMPANY_ID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewErrorResponse(400, "No Parameters"))
	}
	t := domain.Token{CompanyID: companyID, AccessToken: oauth2Token.AccessToken, RefreshToken: oauth2Token.RefreshToken, Expiry: oauth2Token.Expiry}
	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewErrorResponse(400, "Invalid Request"))
	}
	if err := controller.Interactor.Update(t); err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorResponse(500, "Could not update record"))
	}
	return c.JSON(http.StatusCreated, "Authorization successful.")
}

func (controller *TokenController) Refresh(c echo.Context) error {
	ctx := context.Background()
	conf = &oauth2.Config{
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

	// アクセストークンの期限が切れていなければ何もしない
	if tokenDomain.Expiry.After(time.Now().UTC()) {
		return c.JSON(http.StatusOK, tokenDomain)
	}

	// 期限切れならリフレッシュ
	token := oauth2.Token{AccessToken: tokenDomain.AccessToken, TokenType: "Bearer", RefreshToken: tokenDomain.RefreshToken, Expiry: tokenDomain.Expiry}
	tokenSource := conf.TokenSource(ctx, &token)
	client := oauth2.NewClient(ctx, tokenSource)
	url := fmt.Sprintf("https://accounts.secure.freee.co.jp/public_api/token?grant_type=%s&client_id=%s&client_secret=%s&refresh_token=%s", "refresh_token", os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), tokenDomain.RefreshToken)
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenDomain.AccessToken))
	expiry := time.Now().UTC().Add(24 * time.Hour)
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorResponse(500, "Failed to refresh the access token."))
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorResponse(500, "Faild to parse the response body"))
	}

	// AcessToken と RefreshToken を保存
	var resToken domain.Token
	json.Unmarshal(respBody, &resToken)
	t := domain.Token{CompanyID: filter, AccessToken: resToken.AccessToken, RefreshToken: resToken.RefreshToken, Expiry: expiry}
	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorResponse(500, "Invalid Token Structure"))
	}
	if t.AccessToken == "" || t.RefreshToken == "" {
		return c.JSON(http.StatusInternalServerError, common.NewErrorResponse(500, "Returned AcessToken is empty. Try it again!"))
	}
	if err := controller.Interactor.Update(t); err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorResponse(500, "Could not update record"))
	}

	return c.JSON(http.StatusOK, t)
}
