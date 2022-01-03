package usecase

import "kintai-bot/app/domain"

type TokenRepository interface {
	Update(domain.Token) error
	FindByCompanyID(int) (domain.Token, error)
}
