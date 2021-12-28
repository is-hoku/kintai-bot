package usecase

import "kintai-bot/app/domain"

type UserRepository interface {
	Store(domain.User) error
	FindByEmail(string) (domain.User, error)
}
