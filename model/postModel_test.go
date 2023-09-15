package model

import (
	"errors"
	"os"
	"social-api/database"
	"social-api/types"
	"testing"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPostGetEntry(t *testing.T) {
	testtable := []struct {
		input    bson.D
		expected types.Posts
	}{
		{input: bson.D{primitive.E{Key: "img", Value: "image.png"}}, expected: types.Posts{UserID: "633483d5d284eb292ef26363", Image: "image.png", Likes: []string{}}},
	}
	// (the dot env doesnt work with test files)
	// need to replace the with the actual URI when testing
	godotenv.Load(".env")
	uri := os.Getenv("MONGO_URL")
	db_name := os.Getenv("DATABASE_NAME")
	client := database.ConnectDatabase(
		uri,
		db_name)
	userModel := NewPostModel(client)
	for _, tt := range testtable {
		gotPost, modelError := userModel.GetEntry(tt.input)
		if modelError != nil {
			t.Fatalf("error when calling the GetEntry fucntion, :%v", modelError)
		}
		if gotPost.UserID != tt.expected.UserID {
			t.Errorf("wrong users post, got=%s, want=%s", gotPost.UserID, tt.expected.UserID)
		}
		if gotPost.Image != tt.expected.Image {
			t.Errorf("wrong post image, got=%s, want=%s", gotPost.Image, tt.expected.Image)
		}
		if len(gotPost.Likes) != len(tt.expected.Likes) {
			t.Errorf("wrong number of likes, got=%d, want=%d", len(gotPost.Likes), len(tt.expected.Likes))
		}
	}
}

func TestPostAddEntry(t *testing.T) {
	testtable := []struct {
		input    bson.D
		expected error
	}{
		{input: bson.D{primitive.E{Key: "_id", Value: primitive.NewObjectID()}, primitive.E{Key: "userId", Value: "63348350d284eb292ef2635f"}, primitive.E{Key: "img", Value: "image1.png"}, primitive.E{Key: "likes", Value: []string{}}}, expected: nil},
		{input: bson.D{primitive.E{Key: "_id", Value: primitive.NewObjectID()}, primitive.E{Key: "userId", Value: "63348350d284eb292ef2635f"}, primitive.E{Key: "img", Value: "image2.png"}, primitive.E{Key: "likes", Value: []string{}}}, expected: nil},
		{input: bson.D{primitive.E{Key: "img", Value: "image2.png"}, primitive.E{Key: "likes", Value: []string{}}}, expected: errors.New("not enough values given to add user")},
	}
	// (the dot env doesnt work with test files)
	// need to replace the with the actual URI when testing
	godotenv.Load(".env")
	uri := os.Getenv("MONGO_URL")
	db_name := os.Getenv("DATABASE_NAME")
	client := database.ConnectDatabase(
		uri,
		db_name)
	postModel := NewPostModel(client)
	for _, tt := range testtable {
		err := postModel.AddEntry(tt.input)
		if err != nil {
			if err.Error() != tt.expected.Error() {
				t.Errorf("wrong error when inserting, got=%s, want=%s", err.Error(), tt.expected.Error())
			}
		}
	}
}

func TestPostModifyEntry(t *testing.T) {
	testtable := []struct {
		idString string
		inputVal bson.D
		expected error
	}{
		{idString: "64dcfc7fe38b735c64135796", inputVal: bson.D{primitive.E{Key: "$set", Value: bson.D{{Key: "img", Value: "testimage23.png"}}}}, expected: nil},
	}
	// (the dot env doesnt work with test files)
	// need to replace the with the actual URI when testing
	godotenv.Load(".env")
	uri := os.Getenv("MONGO_URL")
	db_name := os.Getenv("DATABASE_NAME")
	client := database.ConnectDatabase(
		uri,
		db_name)
	userModel := NewUserModel(client)
	for _, tt := range testtable {
		id, hexError := primitive.ObjectIDFromHex(tt.idString)
		if hexError != nil {
			t.Fatalf("error converting string to object id, %v", hexError)
		}

		filter := bson.D{primitive.E{Key: "_id", Value: id}}
		err := userModel.ModifyEntry(filter, tt.inputVal)
		if err != nil {
			if err.Error() != tt.expected.Error() {
				t.Errorf("wrong error when updating to model, got=%s, want=%s", err.Error(), tt.expected.Error())
			}
		}
	}
}

func TestPostRemoveEntry(t *testing.T) {
	testtable := []struct {
		idString string
		input    bson.D
		expected error
	}{
		{idString: "64dcfc7fe38b735c64135796", expected: nil},
		{idString: "", input: bson.D{primitive.E{Key: "userId", Value: "63348350d284eb292ef2635f"}}, expected: nil},
		{idString: "", input: bson.D{}, expected: errors.New("empty val value given")},
	}
	// (the dot env doesnt work with test files)
	// need to replace the with the actual URI when testing
	godotenv.Load(".env")
	uri := os.Getenv("MONGO_URL")
	db_name := os.Getenv("DATABASE_NAME")
	client := database.ConnectDatabase(
		uri,
		db_name)
	postModel := NewPostModel(client)
	for _, tt := range testtable {
		if tt.idString != "" {
			id, _ := primitive.ObjectIDFromHex(tt.idString)
			tt.input = bson.D{primitive.E{Key: "_id", Value: id}}
		}
		err := postModel.RemoveEntry(tt.input)
		if err != nil {
			if err.Error() != tt.expected.Error() {
				t.Errorf("wrong error when removing from the model, got=%s, want=%s", err.Error(), tt.expected.Error())
			}
		}
	}

}

func TestPostGetEntryAdvanced(t *testing.T) {
	testtable := []struct {
		input    bson.D
		sort     bson.D
		expected int
	}{
		{input: bson.D{primitive.E{Key: "img", Value: "image1.png"}}, sort: bson.D{primitive.E{Key: "_id", Value: 1}}, expected: 5},
	}
	// (the dot env doesnt work with test files)
	// need to replace the with the actual URI when testing
	godotenv.Load(".env")
	uri := os.Getenv("MONGO_URL")
	db_name := os.Getenv("DATABASE_NAME")
	client := database.ConnectDatabase(
		uri,
		db_name)
	userModel := NewPostModel(client)
	for _, tt := range testtable {
		gotPostArray, modelError := userModel.GetEntryAdvanced(tt.input, tt.sort)
		if modelError != nil {
			t.Fatalf("error when calling the GetEntry fucntion, :%v", modelError)
		}
		if len(gotPostArray) != tt.expected {
			t.Errorf("wrong number of values returned, got=%d, want=%d", len(gotPostArray), tt.expected)
		}
	}
}
