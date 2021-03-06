package domain

import "time"

type Token struct {
	CompanyID    int       `bson:"company_id" json:"company_id"`
	AccessToken  string    `bson:"access_token" json:"access_token"`
	RefreshToken string    `bson:"refresh_token" json:"refresh_token"`
	Expiry       time.Time `bson:"expiry" json:"expiry"`
}

type Tokens []Token
