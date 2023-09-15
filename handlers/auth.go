package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"social-api/helpers"
	"social-api/logger"
	"social-api/model"
	"social-api/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func buildDataBaseType(user *types.Users) bson.D {
	return bson.D{
		primitive.E{Key: "_id", Value: user.UserID},
		primitive.E{Key: "username", Value: user.Username},
		primitive.E{Key: "email", Value: user.Email},
		primitive.E{Key: "password", Value: user.Password},
		primitive.E{Key: "profilePicture", Value: user.ProfilePic},
		primitive.E{Key: "coverPicture", Value: user.CoverPic},
		primitive.E{Key: "follwers", Value: user.Follwers},
		primitive.E{Key: "follwings", Value: user.Follwings},
		primitive.E{Key: "isAdmin", Value: user.IsAdmin},
		primitive.E{Key: "desc", Value: user.Desc},
		primitive.E{Key: "city", Value: user.City},
		primitive.E{Key: "from", Value: user.From},
		primitive.E{Key: "relationship", Value: user.Relationship},
		primitive.E{Key: "created_at", Value: user.CreatedAt},
		primitive.E{Key: "updated_at", Value: user.UpdatedAt},
	}
}

// checks if the given secret matches the admin password
func vailidAdminPassword(secret string) bool {
	pass := os.Getenv("AdminSecretPassword")
	return pass == secret

}

// this needs to be the type to handle all of the
// authorization for the users, will use the modeler interface
// to interact with the database
type AuthHandler struct {
	db  model.Modeler[*types.Users, bson.D]
	log logger.Logger
}

func NewAuthHandler(db model.Modeler[*types.Users, bson.D], logFilePath string) *AuthHandler {
	l := logger.NewLogger()
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("error when makeing the log file for user routes" + err.Error())
	}
	InfoLogger := log.New(file, "INFO: ", log.Ldate|log.Ltime)
	WarningLogger := log.New(file, "WARNING: ", log.Ldate|log.Ltime)
	ErrorLogger := log.New(file, "ERROR: ", log.Ldate|log.Ltime)
	FatalLogger := log.New(file, "FATAL: ", log.Ldate|log.Ltime)
	l.AddLogger(logger.INFO, InfoLogger)
	l.AddLogger(logger.WARNING, WarningLogger)
	l.AddLogger(logger.ERROR, ErrorLogger)
	l.AddLogger(logger.FATAL, FatalLogger)
	return &AuthHandler{
		db:  db,
		log: l,
	}
}

func (ah *AuthHandler) Test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello this is the auth handler test"))
	id, _ := primitive.ObjectIDFromHex("633356b45715fd08fc68798e")
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	user, err := ah.db.GetEntry(filter)
	ah.log.WriteToLogger(logger.INFO, "test endpoint was hit")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", user)
}

// handle to login of the user and send the user data to the client
// login only needs email or username and password
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	requestUser, parseError := helpers.ParseBody(r.Body, types.AuthUserRequest{})
	if parseError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error when parsing the request"))
		fmt.Println(parseError)
	}
	var searchKey string
	var searchParam string
	if requestUser.Email != "" && (requestUser.UserName == "" && requestUser.Password != "") {
		searchKey = "email"
		searchParam = requestUser.Email
	} else if requestUser.UserName != "" && (requestUser.Email == "" && requestUser.Password != "") {
		searchKey = "username"
		searchParam = requestUser.UserName
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("either username or email is required with password"))
	}
	fmt.Println(searchParam)
	key := bson.D{primitive.E{Key: searchKey, Value: searchParam}}
	dbUser, dbErr := ah.db.GetEntry(key)
	if dbErr != nil {
		if errors.Is(dbErr, mongo.ErrNoDocuments) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("user not found in database"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("unknow server error"))
			fmt.Println("unknown error when getting user from db", dbErr)
		}
		return
	}
	fmt.Println(dbUser)
	// returns nil if the passwords are the same
	correctUser := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(requestUser.Password))
	if correctUser != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("incorrect password given"))
		return
	} else {
		// dont send the password hash to the client
		dbUser.Password = "************"
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dbUser)
	}
}

// will handle the creation of the user in the database, will send user data back after creation
func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	requestUser, parseError := helpers.ParseBody(r.Body, types.AuthUserRequest{})
	if parseError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error when parsing the request"))
		fmt.Println(parseError)
	}
	if !types.ValidAuthUser(requestUser) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("invalid user given"))
		return
	}
	hashedPass, hashErr := bcrypt.GenerateFromPassword([]byte(requestUser.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		fmt.Printf("error when hashing the user password: %v", hashErr)
	}
	user := types.NewUser()
	// this only covers the basic fields to make the user
	// need to add all the other feild that get the default values from NewUser
	user.Email = requestUser.Email
	user.Username = requestUser.UserName
	user.Password = string(hashedPass)
	if requestUser.AdminSecret != "" {
		if !vailidAdminPassword(requestUser.AdminSecret) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("invalid atempt to make a admin account"))
			// dont create any user if failed admin attempt
			return
		}
		user.IsAdmin = true
	}
	dbUser := buildDataBaseType(user)
	ah.db.AddEntry(dbUser)
	fmt.Printf("%+v", dbUser)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user successfully registed"))
}

func (ah *AuthHandler) HandleNotFound(w http.ResponseWriter, r *http.Request, msg string) {
	ah.log.WriteToLogger(logger.WARNING, "invalid url was given to post handlers"+r.URL.Path)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(msg))
}
