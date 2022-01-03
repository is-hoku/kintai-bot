package database

import (
	"context"

	"kintai-bot/app/domain"

	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository struct {
	DBHandler
}

func (repo *UserRepository) Store(u domain.User) (err error) {
	_, err = repo.InsertOne(context.TODO(), u)
	return
}

func (repo *UserRepository) FindByEmail(email string) (user domain.User, err error) {
	//filter := []byte(fmt.Sprintf(`{"email": %s}`, email))
	filter := bson.D{{"email", email}}
	err = repo.FindOne(context.TODO(), filter).Decode(&user)
	return
}
