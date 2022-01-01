package domain

type User struct {
	Email   string `bson:"email" json:"email"`
	FreeeID int    `bson:"freee_id" json:"freee_id"`
}

type Users []User
