package domain

type User struct {
	Email   string `bson:"email" json:"email"`
	FreeeID int    `bson:"freeeID" json:"freeeID"`
}

type Users []User
