package database

import (
	"context"

	"kintai-bot/app/domain"

	"go.mongodb.org/mongo-driver/bson"
)

type TokenRepository struct {
	TokenDBHandler
}

func (repo *TokenRepository) Update(t domain.Token) (err error) {
	filter := bson.D{{"company_id", t.CompanyID}}
	update := bson.D{{"$set", bson.D{{"access_token", t.AccessToken}, {"refresh_token", t.RefreshToken}, {"expiry", t.Expiry}}}}
	_, err = repo.UpdateOne(context.TODO(), filter, update)
	return
}

func (repo *TokenRepository) FindByCompanyID(companyID int) (token domain.Token, err error) {
	filter := bson.D{{"company_id", companyID}}
	err = repo.FindOne(context.TODO(), filter).Decode(&token)
	return
}
