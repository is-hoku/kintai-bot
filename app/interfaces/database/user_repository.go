package database

import (
	"context"
	"fmt"

	"github.com/is-hoku/kintai-bot/domain"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository struct {
	DBHandler
}

func (repo *UserRepository) Store(u domain.User) (err error) {
	result, err := repo.InsertOne(context.TODO(), u)
	fmt.Printf("%s", result)
	return
}

func (repo *UserRepository) FindByEmail(email string) (user domain.User, err error) {
	//filter := []byte(fmt.Sprintf(`{"email": %s}`, email))
	filter := bson.D{{"email", email}}
	err = repo.FindOne(context.TODO(), filter).Decode(&user)
	return
}
