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

func (UserInteractor *UserInteractor) UserByEmail(email string) (user domain.User, err error) {
	//filter := []byte(fmt.Sprintf(`{"email": %s}`, email))
	user, err = UserInteractor.UserRepository.FindByEmail(email)
	return
}
