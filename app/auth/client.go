package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
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

func Auth() {
	//if err := godotenv.Load(".env"); err != nil {
	//	log.Fatal(err)
	//}
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		//Scopes:       []string{"SCOPE"},
		RedirectURL: "urn:ietf:wg:oauth:2.0:oob", // TODO: コールバックを指定して渡す
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.secure.freee.co.jp/public_api/authorize",
			TokenURL: "https://accounts.secure.freee.co.jp/public_api/token"},
	}
	state := fmt.Sprintf("st%d", time.Now().UnixNano())
	url := conf.AuthCodeURL(state) // 認可ページのURL
	fmt.Println(url)
	// Authorization Code を受け取る
	// ...
	// トークンエンドポイントにリクエスト、アクセストークン・リフレッシュトークンを取得
	oauth2Token, err := conf.Exchange(ctx, os.Getenv("AUTHORIZATION_CODE"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(oauth2Token)
	fmt.Println("AccessToken: ", oauth2Token.AccessToken)
	fmt.Println("RefreshToken: ", oauth2Token.RefreshToken)
	client := conf.Client(ctx, oauth2Token)
	resp, err := client.Get("https://api.freee.co.jp/hr/api/v1/users/me")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	byteArray, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray))
	user := &User{}
	if err = json.Unmarshal(byteArray, user); err != nil {
		log.Fatal(err)
	}
	i, err := contains(user.Companies, os.Getenv("COMPANY_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	env := map[string]string{
		"CLIENT_ID":          os.Getenv("CLIENT_ID"),
		"CLIENT_SECRET":      os.Getenv("CLIENT_SECRET"),
		"AUTHORIZATION_CODE": os.Getenv("AUTHORIZATION_CODE"),
		"COMPANY_ID":         strconv.Itoa(user.Companies[i].ID),
	}
	err = godotenv.Write(env, ".env")
}
