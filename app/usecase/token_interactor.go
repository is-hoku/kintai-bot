package usecase

import (
	"kintai-bot/app/domain"
)

type TokenInteractor struct {
	TokenRepository TokenRepository
}

func (interactor *TokenInteractor) Update(u domain.Token) (err error) {
	err = interactor.TokenRepository.Update(u)
	return
}

func (interactor *TokenInteractor) TokenByCompanyID(company_id int) (token domain.Token, err error) {
	token, err = interactor.TokenRepository.FindByCompanyID(company_id)
	return
}
