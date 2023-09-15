package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"social-api/helpers"
	"social-api/logger"
	"social-api/model"
	"social-api/types"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// will make the new user with the given data from the client
// and the existing user from the database
func updateUserData(dbUser *types.Users, rUser *types.RequestUser) *types.Users {
	// do this so the client doesnt need to resend data in the database
	finalUser := types.NewUser()
	if rUser.City != "" {
		finalUser.City = rUser.City
	} else {
		finalUser.City = dbUser.City
	}
	if rUser.Email != "" {
		finalUser.Email = rUser.Email
	} else {
		finalUser.Email = dbUser.Email
	}
	if rUser.CoverPic != "" {
		finalUser.CoverPic = rUser.CoverPic
	} else {
		finalUser.CoverPic = dbUser.CoverPic
	}
	if rUser.Desc != "" {
		finalUser.Desc = rUser.Desc
	} else {
		finalUser.Desc = dbUser.Desc
	}
	if rUser.From != "" {
		finalUser.From = rUser.From
	} else {
		finalUser.From = dbUser.From
	}
	if rUser.Password != "" {
		hashedPass, _ := bcrypt.GenerateFromPassword([]byte(rUser.Password), bcrypt.DefaultCost)
		finalUser.Password = string(hashedPass)
	} else {
		finalUser.Password = dbUser.Password
	}
	if rUser.ProfilePic != "" {
		finalUser.ProfilePic = rUser.ProfilePic
	} else {
		finalUser.ProfilePic = dbUser.ProfilePic
	}
	if rUser.Username != "" {
		finalUser.Username = rUser.Username
	} else {
		finalUser.Username = dbUser.Username
	}
	if rUser.Relationship != dbUser.Relationship {
		finalUser.Relationship = rUser.Relationship
	} else {
		finalUser.Relationship = dbUser.Relationship
	}
	// updated at will be set when newUser is called
	finalUser.CreatedAt = dbUser.CreatedAt
	// makes sure the userid stays the same
	finalUser.UserID = dbUser.UserID

	return finalUser
}

type UserHandler struct {
	db  model.Modeler[*types.Users, bson.D]
	log logger.Logger
}

func NewUserHandler(db model.Modeler[*types.Users, bson.D], logFilePath string) *UserHandler {
	l := logger.NewLogger()
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("error when making the log file for user routes" + err.Error())
	}
	InfoLogger := log.New(file, "INFO: ", log.Ldate|log.Ltime)
	WarningLogger := log.New(file, "WARNING: ", log.Ldate|log.Ltime)
	ErrorLogger := log.New(file, "ERROR: ", log.Ldate|log.Ltime)
	FatalLogger := log.New(file, "FATAL: ", log.Ldate|log.Ltime)
	l.AddLogger(logger.INFO, InfoLogger)
	l.AddLogger(logger.WARNING, WarningLogger)
	l.AddLogger(logger.ERROR, ErrorLogger)
	l.AddLogger(logger.FATAL, FatalLogger)
	return &UserHandler{
		db:  db,
		log: l,
	}
}

func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, id string) {
	key := bson.D{primitive.E{Key: "_id", Value: id}}
	user, dbError := uh.db.GetEntry(key)
	if dbError != nil {
		helpers.HandleDbError(dbError, w, uh.log, fmt.Sprintf("error when getting user with id %s", id))
		return
	}
	// censer the password before sending data to client
	user.Password = "********"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, id string) {
	//	fmt.Println("have not make the update user handler yet", id)
	//	w.WriteHeader(http.StatusNotImplemented)
	//	w.Write([]byte("have not make the update user handler yet"))
	key := bson.D{primitive.E{Key: "_id", Value: id}}
	dbuser, dbError := uh.db.GetEntry(key)
	if dbError != nil {
		helpers.HandleDbError(dbError, w, uh.log, fmt.Sprintf("error when getting user with id %s", id))
		return
	}
	rUser, parseError := helpers.ParseBody(r.Body, types.RequestUser{})
	if parseError != nil {
		helpers.HandleParserError(parseError, w, uh.log)
		return
	}
	if rUser.Username == "" || rUser.Username != dbuser.Username {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid username given, username is required to update account"))
		return
	}
	correctUser := bcrypt.CompareHashAndPassword([]byte(dbuser.Password), []byte(rUser.Password))
	if correctUser != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("incorrect password given"))
		return
	}
	if rUser.UserID == dbuser.UserID.Hex() || rUser.IsAdmin {
		newUser := updateUserData(dbuser, rUser)
		val := bson.D{
			primitive.E{Key: "username", Value: newUser.Username},
			primitive.E{Key: "email", Value: newUser.Email},
			primitive.E{Key: "password", Value: newUser.Password},
			primitive.E{Key: "profilePicture", Value: newUser.ProfilePic},
			primitive.E{Key: "coverPicture", Value: newUser.CoverPic},
			primitive.E{Key: "desc", Value: newUser.Desc},
			primitive.E{Key: "city", Value: newUser.City},
			primitive.E{Key: "from", Value: newUser.From},
			primitive.E{Key: "relationship", Value: newUser.Relationship},
			primitive.E{Key: "created_at", Value: newUser.CreatedAt},
			primitive.E{Key: "updated_at", Value: newUser.UpdatedAt},
		}
		if err := uh.db.ModifyEntry(key, val); err != nil {
			helpers.HandleDbError(err, w, uh.log, "error when updating the user")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("user has been updated"))
		return
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("you are not authorized to modify this users account"))
		return
	}

}

func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, id string) {
	//	fmt.Println("have not make the delete user handler yet", id)
	//	w.WriteHeader(http.StatusNotImplemented)
	//	w.Write([]byte("have not make the update user handler yet"))
	key := bson.D{primitive.E{Key: "_id", Value: id}}
	dbuser, dbError := uh.db.GetEntry(key)
	if dbError != nil {
		helpers.HandleDbError(dbError, w, uh.log, fmt.Sprintf("error when getting user with id %s", id))
		return
	}
	rUser, parseError := helpers.ParseBody(r.Body, types.RequestUser{})
	if parseError != nil {
		helpers.HandleParserError(parseError, w, uh.log)
		return
	}
	if rUser.Username == "" || rUser.Username != dbuser.Username {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid username given, username is required to update account"))
		return
	}
	correctUser := bcrypt.CompareHashAndPassword([]byte(dbuser.Password), []byte(rUser.Password))
	if correctUser != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("incorrect password given"))
		return
	}
	if rUser.UserID == dbuser.UserID.Hex() || rUser.IsAdmin {
		if err := uh.db.RemoveEntry(key); err != nil {
			helpers.HandleDbError(err, w, uh.log, "error when deleteing user: "+rUser.UserID)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("user has been deleted"))
		return
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("you are not authorized to delete this users account"))
		return
	}
}

func (uh *UserHandler) FollowUnfollow(w http.ResponseWriter, r *http.Request, followId string) {
	requestUser, err := helpers.ParseBody(r.Body, types.AuthUserRequest{})
	if err != nil {
		helpers.HandleParserError(err, w, uh.log)
		return
	}
	if requestUser.UserId == followId {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("cannot follow yourself"))
		return
	}
	currUserKey := bson.D{primitive.E{Key: "_id", Value: requestUser.UserId}}
	currentUser, cErr := uh.db.GetEntry(currUserKey)
	if cErr != nil {
		helpers.HandleDbError(cErr, w, uh.log)
	}
	followUserKey := bson.D{primitive.E{Key: "_id", Value: followId}}
	user, uErr := uh.db.GetEntry(followUserKey)
	if uErr != nil {
		helpers.HandleDbError(uErr, w, uh.log)
	}
	var updatedFollowerArray []string
	var updatedFollowingArray []string
	var option bool
	currentTime := time.Now()

	if !helpers.Includes(user.Follwers, currentUser.UserID.Hex()) {
		updatedFollowerArray = append(user.Follwers, currentUser.UserID.Hex())
		updatedFollowingArray = append(currentUser.Follwings, user.UserID.Hex())
		option = true
	} else {
		newUserArray, userErr := helpers.RemoveElement(user.Follwers, currentUser.UserID.Hex())
		newCurrentUserArray, currentUserErr := helpers.RemoveElement(currentUser.Follwings, user.UserID.Hex())
		if userErr != nil || currentUserErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("unknow error when unfollowing user"))
			return
		}
		updatedFollowerArray = newUserArray
		updatedFollowingArray = newCurrentUserArray
		option = false
	}

	currUserVal := bson.D{
		primitive.E{Key: "follwings", Value: updatedFollowingArray},
		primitive.E{Key: "updated_at", Value: currentTime},
	}
	if err := uh.db.ModifyEntry(currUserKey, currUserVal); err != nil {
		helpers.HandleDbError(err, w, uh.log, "error when updating current ussers follings")
		return
	}
	followUserVal := bson.D{
		primitive.E{Key: "follwers", Value: updatedFollowerArray},
		primitive.E{Key: "updated_at", Value: currentTime},
	}
	if err := uh.db.ModifyEntry(followUserKey, followUserVal); err != nil {
		helpers.HandleDbError(err, w, uh.log, "error when updating users followers")
		return
	}
	if option {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("user has been followed"))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("user has been unfollowed"))
		return
	}
}

func (uh *UserHandler) HandleNotFound(w http.ResponseWriter, r *http.Request, msg string) {
	uh.log.WriteToLogger(logger.WARNING, "invalid url was given to post handlers"+r.URL.Path)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(msg))
}
