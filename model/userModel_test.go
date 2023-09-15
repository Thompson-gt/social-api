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

func TestUserGetEntry(t *testing.T) {
	testtable := []struct {
		userId   string
		input    bson.D
		expected types.Users
	}{
		{userId: "", input: bson.D{primitive.E{Key: "username", Value: "bob"}}, expected: types.Users{Username: "bob", Email: "bob@gmail.com", Password: "$2b$10$t4UkW8gp83Mmk2O8IXgKseOrvH8Eg2SaYaU4Az5rRWrAVq8B5KdfW", ProfilePic: "", CoverPic: "", Follwers: []string{"testUser"}, Follwings: []string{"testUser"}, IsAdmin: false}},
		{userId: "633356b45715fd08fc68798e", expected: types.Users{Username: "gabe", Email: "gabe@gmail.com", Password: "$2b$10$t4UkW8gp83Mmk2O8IXgKseOrvH8Eg2SaYaU4Az5rRWrAVq8B5KdfW", ProfilePic: "", CoverPic: "", Follwers: []string{"testUser"}, Follwings: []string{}, IsAdmin: false}},
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
		if tt.userId != "" {
			id, _ := primitive.ObjectIDFromHex(tt.userId)
			tt.input = bson.D{primitive.E{Key: "_id", Value: id}}
		}
		gotUser, modelError := userModel.GetEntry(tt.input)
		if modelError != nil {
			t.Fatalf("error when calling the GetEntry fucntion, :%v", modelError)
		}
		if gotUser.Username != tt.expected.Username {
			t.Errorf("wrong username, got=%s, want=%s", gotUser.Username, tt.expected.Username)
		}
		if gotUser.Email != tt.expected.Email {
			t.Errorf("wrong email, got=%s, want=%s", gotUser.Email, tt.expected.Email)
		}
		if len(gotUser.Follwers) != len(tt.expected.Follwers) {
			t.Errorf("wrong number of followers, got=%d, want=%d", len(gotUser.Follwers), len(tt.expected.Follwers))
		}
		if len(gotUser.Follwings) != len(tt.expected.Follwings) {
			t.Errorf("wrong number of following, got=%d, want=%d", len(gotUser.Follwings), len(tt.expected.Follwings))
		}
	}

}

// when used in real api pass will be hashed and all values will be added
func TestUserAddEntry(t *testing.T) {
	testtable := []struct {
		input    bson.D
		expected error
	}{
		{input: bson.D{primitive.E{Key: "_id", Value: primitive.NewObjectID()}, primitive.E{Key: "username", Value: "tommy"}, primitive.E{Key: "email", Value: "tommy@gmail.com"}, primitive.E{Key: "password", Value: "tommypassword"}}, expected: nil},
		{input: bson.D{primitive.E{Key: "_id", Value: primitive.NewObjectID()}, primitive.E{Key: "username", Value: "gilson"}, primitive.E{Key: "email", Value: "gilson@gmail.com"}, primitive.E{Key: "password", Value: "gilsonpassword"}}, expected: nil},
		{input: bson.D{primitive.E{Key: "password", Value: "failpassword"}}, expected: errors.New("not enough values given to add user")},
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
		err := userModel.AddEntry(tt.input)
		if err != nil {
			if err.Error() != tt.expected.Error() {
				t.Errorf("wrong error when adding to model, got=%s, want=%s", err.Error(), tt.expected.Error())
			}
		}
	}

}

func TestUserModifyEntry(t *testing.T) {
	testtable := []struct {
		inputFilter bson.D
		inputVal    bson.D
		expected    error
	}{
		{inputFilter: bson.D{primitive.E{Key: "username", Value: "tommy"}}, inputVal: bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "password", Value: "tommyPassword2"}}}}, expected: nil},
		{inputFilter: bson.D{primitive.E{Key: "username", Value: "gilson"}}, inputVal: bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "password", Value: "gilsonPassword2"}}}}, expected: nil},
		{inputFilter: bson.D{primitive.E{Key: "username", Value: "gilson"}}, inputVal: bson.D{{}}, expected: errors.New("no empty update value given")},
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
		err := userModel.ModifyEntry(tt.inputFilter, tt.inputVal)
		if err != nil {
			if err.Error() != tt.expected.Error() {
				t.Errorf("wrong error when updating to model, got=%s, want=%s", err.Error(), tt.expected.Error())
			}
		}
	}
}

func TestUserRemoveEntry(t *testing.T) {
	testtable := []struct {
		input    bson.D
		expected error
	}{
		{input: bson.D{primitive.E{Key: "username", Value: "gilson"}}, expected: nil},
		{input: bson.D{}, expected: errors.New("empty val value given")},
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
		err := userModel.RemoveEntry(tt.input)
		if err != nil {
			if err.Error() != tt.expected.Error() {
				t.Errorf("wrong error when removing from the model, got=%s, want=%s", err.Error(), tt.expected.Error())
			}
		}
	}

}
func TestUserGetEntryAdvanced(t *testing.T) {
	testtable := []struct {
		input    bson.D
		sort     bson.D
		expected int
	}{
		{input: bson.D{primitive.E{Key: "profilePicture", Value: ""}}, sort: bson.D{primitive.E{Key: "_id", Value: 1}}, expected: 3},
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
		usersArray, modelError := userModel.GetEntryAdvanced(tt.input, tt.sort)
		if modelError != nil {
			t.Fatalf("error when calling the GetEntry fucntion, :%v", modelError)
		}
		if len(usersArray) != tt.expected {
			t.Errorf("wrong number of values returned, got=%d, want=%d", len(usersArray), tt.expected)
		}
	}
}
