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
	companyID := t.CompanyID
	filter := bson.D{{"company_id", companyID}}
	_, err = repo.UpdateOne(context.TODO(), filter, t)
	return
}

func (repo *TokenRepository) FindByCompanyID(companyID int) (token domain.Token, err error) {
	filter := bson.D{{"company_id", companyID}}
	err = repo.FindOne(context.TODO(), filter).Decode(&token)
	return
}
