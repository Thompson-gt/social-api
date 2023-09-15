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
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostHandler struct {
	db  model.Modeler[*types.Posts, bson.D]
	log logger.Logger
}

func NewPostHandler(db model.Modeler[*types.Posts, bson.D], logFilePath string) *PostHandler {
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
	return &PostHandler{
		db:  db,
		log: l,
	}
}

func (ph *PostHandler) Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello this is the post handler test"))
	//	filter := bson.D{primitive.E{Key: "img", Value: "image1.png"}}
	//	sort := bson.D{primitive.E{Key: "_id", Value: -1}}
	//	posts, _ := ph.db.GetEntryAdvanced(filter, sort)
	ph.log.WriteToLogger(logger.INFO, "test endpoint was hit")
	//fmt.Println(len(posts))
	//for _, post := range posts {
	//	fmt.Println(post)
	//}
}

func (ph *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	requestPost, parseError := helpers.ParseBody(r.Body, types.RequestPost{})
	if parseError != nil {
		helpers.HandleParserError(parseError, w, ph.log)
		return
	}
	if !types.ValidReqestPost(requestPost) {
		ph.log.WriteToLogger(logger.WARNING, "incomplete data given to create post handler")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("not enough data given to create new post"))
		return
	}
	key := bson.D{
		primitive.E{Key: "userId", Value: requestPost.UserId},
		primitive.E{Key: "desc", Value: requestPost.Desc},
		primitive.E{Key: "img", Value: requestPost.Image},
	}
	dberr := ph.db.AddEntry(key)
	if dberr != nil {
		helpers.HandleDbError(dberr, w, ph.log, "failed to add post to database")
		return
	} else {
		ph.log.WriteToLogger(logger.INFO, "post created in db")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("post created"))
	}

}

func (ph *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request, id string) {
	requestPost, parseError := helpers.ParseBody(r.Body, types.RequestPost{})
	if parseError != nil {
		helpers.HandleParserError(parseError, w, ph.log, "error when parsing post for UpdatePost")
		return
	}
	// the id needs to be a stirng when quering the database(i think)
	key := bson.D{primitive.E{Key: "_id", Value: id}}
	dbPost, dbError := ph.db.GetEntry(key)
	if dbError != nil {
		helpers.HandleDbError(dbError, w, ph.log)
		return
	}
	if dbPost.UserID != requestPost.UserId {
		ph.log.WriteToLogger(logger.WARNING, "user attempted to modify someone elses post")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("not allowed to update other peoples post"))
		return
	}
	var newImg string
	var newDesc string
	// this allows the client to update one or the other without needing
	// so send redundent data to endpoint
	if requestPost.Image == "" {
		newImg = dbPost.Image
	} else {
		newImg = requestPost.Image
	}
	if requestPost.Desc == "" {
		newDesc = dbPost.Desc
	} else {
		newDesc = requestPost.Desc
	}
	val := bson.D{
		primitive.E{Key: "img", Value: newImg},
		primitive.E{Key: "desc", Value: newDesc},
		primitive.E{Key: "updated_at", Value: time.Now()},
	}
	if updateError := ph.db.ModifyEntry(key, val); updateError != nil {
		helpers.HandleDbError(updateError, w, ph.log, fmt.Sprintf("error when updatin post with id of : %s", id))
		return
	}
	ph.log.WriteToLogger(logger.INFO, "post in database was updated")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("post was successfully updated"))

}

func (ph *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request, id string) {
	requestPost, parseError := helpers.ParseBody(r.Body, types.RequestPost{})
	if parseError != nil {
		helpers.HandleParserError(parseError, w, ph.log, "error when parsing in DeletePost")
		return
	}
	// the id needs to be a stirng when quering the database(i think)
	key := bson.D{primitive.E{Key: "_id", Value: id}}
	dbPost, dbError := ph.db.GetEntry(key)
	if dbError != nil {
		helpers.HandleDbError(dbError, w, ph.log, fmt.Sprintf("error when getting post with id of %s", id))
		return
	}
	if dbPost.UserID != requestPost.UserId {
		ph.log.WriteToLogger(logger.WARNING, "attempt to delete someones else post")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("not allowed to update other peoples post"))
		return
	}
	if removeErr := ph.db.RemoveEntry(key); removeErr != nil {
		helpers.HandleDbError(removeErr, w, ph.log, "error when removing the post from database")
		return
	} else {
		ph.log.WriteToLogger(logger.INFO, "post has been deleted form the database")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("post has been deleted"))
		return
	}
}

func (ph *PostHandler) GetPost(w http.ResponseWriter, r *http.Request, id string) {
	key := bson.D{primitive.E{Key: "_id", Value: id}}
	dbPost, dbError := ph.db.GetEntry(key)
	if dbError != nil {
		helpers.HandleDbError(dbError, w, ph.log, fmt.Sprintf("error when getting post with id of %s", id))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dbPost)

}

func (ph *PostHandler) HandleLikeDislike(w http.ResponseWriter, r *http.Request, postId string) {
	requestUser, err := helpers.ParseBody(r.Body, types.AuthUserRequest{})
	if err != nil {
		helpers.HandleParserError(err, w, ph.log)
		return
	}
	if requestUser.UserId == "" {
		ph.log.WriteToLogger(logger.WARNING, "incomplete data sent to like/dislike handler")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("need to provide userId of user liking the post"))
		return
	}
	postKey := bson.D{primitive.E{Key: "_id", Value: postId}}
	dbPost, dbError := ph.db.GetEntry(postKey)
	if dbError != nil {
		helpers.HandleDbError(dbError, w, ph.log, fmt.Sprintf("error when getting post with id of %s", postId))
		return
	}
	// attempt to remove the userId from the Likes array(dislike post),
	// if error is return then append the userId to the likes array(like post)
	newLikesArray, likeError := helpers.RemoveElement(dbPost.Likes, requestUser.UserId)
	if likeError != nil && errors.Is(likeError, errors.New("element not in array")) {
		dbPost.Likes = append(dbPost.Likes, requestUser.UserId)
		val := bson.D{
			primitive.E{Key: "likes", Value: dbPost.Likes},
			primitive.E{Key: "updated_at", Value: time.Now()},
		}
		if err := ph.db.ModifyEntry(postKey, val); err != nil {
			helpers.HandleDbError(err, w, ph.log, fmt.Sprintf("error when liking the post with id of: %s", postId))
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("post has been liked"))
			return
		}
	} else {
		val := bson.D{
			primitive.E{Key: "likes", Value: newLikesArray},
			primitive.E{Key: "updated_at", Value: time.Now()},
		}
		if err := ph.db.ModifyEntry(postKey, val); err != nil {
			helpers.HandleDbError(err, w, ph.log, fmt.Sprintf("error when liking the post with id of: %s", postId))
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("post has been unliked"))
			return
		}
	}
}

// will return the timeline for the user whos id was provided into the request
func (ph *PostHandler) GetTimeLine(w http.ResponseWriter, r *http.Request, requestUser *types.Users) {
	var friendPosts []*types.Posts
	if len(requestUser.Follwings) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("you are not following anyone yet"))
		return
	}
	// can ignore this error because its ok if the current user doesnt have any posts
	mypost, _ := ph.db.GetEntry(bson.D{primitive.E{Key: "userId", Value: requestUser.UserID.Hex()}})
	for _, user := range requestUser.Follwings {
		post, err := ph.db.GetEntry(bson.D{primitive.E{Key: "userId", Value: user}})
		if err != nil {
			helpers.HandleDbError(err, w, ph.log, "error when getting friends posts")
			return
		}
		friendPosts = append(friendPosts, post)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(append(friendPosts, mypost))
}

func (ph *PostHandler) HandleNotFound(w http.ResponseWriter, r *http.Request, msg string) {
	ph.log.WriteToLogger(logger.WARNING, "invalid url was given to post handlers"+r.URL.Path)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(msg))
}
