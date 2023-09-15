package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	UserID       primitive.ObjectID `bson:"_id"`
	Username     string             `bson:"username"`
	Email        string             `bson:"email"`
	Password     string             `bson:"password"`
	ProfilePic   string             `bson:"profilePicture"`
	CoverPic     string             `bson:"coverPicture"`
	Follwers     []string           `bson:"follwers"`
	Follwings    []string           `bson:"follwings"`
	IsAdmin      bool               `bson:"isAdmin"`
	Desc         string             `bson:"desc"`
	City         string             `bson:"city"`
	From         string             `bson:"from"`
	Relationship int                `bson:"relationship"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"` // need to update this whenever changing data
}

func NewUser() *Users {
	user := &Users{
		UserID:       primitive.NewObjectID(),
		Username:     "default",
		Email:        "default@default.com",
		Password:     "defaultPassword",
		ProfilePic:   "",
		CoverPic:     "",
		Follwers:     []string{},
		Follwings:    []string{},
		IsAdmin:      false,
		Desc:         "",
		City:         "",
		From:         "",
		Relationship: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return user
}

// user is valid if username, password, email is given (other fields have default values)
// returns true if the data of the user is valid
func ValidUser(user Users) bool {
	if user.Username == "default" || user.Email == "default@default.com" || user.Password == "defaultPassword" {
		return false
	}
	return true
}
