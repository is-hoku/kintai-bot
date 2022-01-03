package usecase

import (
	"kintai-bot/app/domain"
)

type UserInteractor struct {
	UserRepository UserRepository
}

func (interactor *UserInteractor) Add(u domain.User) (err error) {
	err = interactor.UserRepository.Store(u)
	return
}

func (interactor *UserInteractor) UserByEmail(email string) (user domain.User, err error) {
	user, err = interactor.UserRepository.FindByEmail(email)
	return
}
