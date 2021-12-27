package usecase

import "github.com/is-hoku/kintai-bot/domain"

type UserRepository interface {
	Store(domain.User) error
	FindByEmail(string) (domain.User, error)
}
