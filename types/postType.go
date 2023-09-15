package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Posts struct {
	PostID    primitive.ObjectID `bson:"_id"`
	UserID    string             `bson:"userId"`
	Image     string             `bson:"img"`
	Desc      string             `bson:"desc"`
	Likes     []string           `bson:"likes"` //will be a array of userid of people who liked it
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"` // need to update this whenever changing data
}

func NewPost() *Posts {
	post := &Posts{
		PostID:    primitive.NewObjectID(),
		UserID:    "defaultUserID",
		Image:     "default.png",
		Likes:     []string{},
		CreatedAt: time.Now(),
	}
	return post
}

// a post is valid if a unique userid was given
// return true if post is valid
func ValidPost(post Posts) bool {
	if post.UserID == "defaultUserID" || post.Image == "default.png" {
		return false
	}
	return true
}
