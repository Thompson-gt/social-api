package main

import (
	"fmt"
	"net/http"
	"os"
	"social-api/database"
	"social-api/handlers"
	"social-api/helpers"
	"social-api/model"
	"social-api/types"
	"strings"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// make sure to add some logging later
const postEndpointLogPath string = "postLogFile.txt"

// auth and user enpoint will use this log file
const userEndpointLogPath string = "userLogFile.txt"

func main() {
	godotenv.Load(".env")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	uri := os.Getenv("MONGO_URL")
	databaseName := os.Getenv("DATABASE_NAME")
	dbClient := database.ConnectDatabase(uri, databaseName)
	userModel := model.NewUserModel(dbClient)

	AuthHandlers := handlers.NewAuthHandler(userModel, userEndpointLogPath)
	UserHandlers := handlers.NewUserHandler(userModel, userEndpointLogPath)
	PostsHandlers := handlers.NewPostHandler(model.NewPostModel(dbClient), postEndpointLogPath)

	http.HandleFunc("/timeline/", func(w http.ResponseWriter, r *http.Request) {
		paths := strings.Split(r.URL.Path, "/")
		requestUser, err := helpers.ParseBody(r.Body, types.AuthUserRequest{})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error when parsing the request"))
			return
		}
		dbUser, dbErr := userModel.GetEntry(bson.D{primitive.E{Key: "_id", Value: requestUser.UserId}})
		if dbErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error when getting given user"))
			return
		}
		switch len(paths) - 1 {
		case 2:
			switch paths[2] {
			case "all":
				PostsHandlers.GetTimeLine(w, r, dbUser)
				// this leaves the option to add filters to the timeline
				//(dont know how i would do that right now though)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("url does not match any timeline endpoint"))
		}
	})

	http.HandleFunc("/tester", PostsHandlers.Test)
	http.HandleFunc("/auth/", func(w http.ResponseWriter, r *http.Request) {
		paths := strings.Split(r.URL.Path, "/")
		switch len(paths) - 1 {
		case 2:
			// will switch on the string of to dermin which user endpoint handler
			// to use
			switch paths[2] {
			case "register":
				AuthHandlers.Register(w, r)
			case "login":
				AuthHandlers.Login(w, r)
			case "test":
				AuthHandlers.Test(w, r)
			}
		default:
			AuthHandlers.HandleNotFound(w, r, "no user endpoint for given url")
			return
		}
	})
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		paths := strings.Split(r.URL.Path, "/")
		fmt.Println(paths)
		switch len(paths) - 1 {
		case 3:
			id := paths[2]
			if len(id) <= 1 {
				UserHandlers.HandleNotFound(w, r, "no user id was given in path")
				return
			}
			if paths[3] == "follow" || paths[3] == "unfollow" {
				fmt.Println("follow/unfollow user hit")
				UserHandlers.FollowUnfollow(w, r, id)
			} else {
				UserHandlers.HandleNotFound(w, r, "invaild option was given for user id")
			}
		case 2:
			id := paths[2]
			// no id will be smaller than 2 chars
			if len(id) <= 1 {
				// if id not in the path then nothing can be done wtih the users handlers
				UserHandlers.HandleNotFound(w, r, "invalid path was given to the user route")
				return
			} else {
				if r.Method == "PUT" {
					fmt.Println("update user hit")
					UserHandlers.UpdateUser(w, r, id)
				} else if r.Method == "DELETE" {
					fmt.Println("delete usser hit")
					UserHandlers.DeleteUser(w, r, id)
				} else if r.Method == "GET" {
					fmt.Println("get user hit")
					UserHandlers.GetUser(w, r, id)
				} else {
					fmt.Println("")
					UserHandlers.HandleNotFound(w, r, "unsupported method  given to user route")
				}
			}
		default:
			// need to replce with logging then send a responce to the user
			AuthHandlers.HandleNotFound(w, r, "unexpected auth endpoint")
			return
		}
	})

	// will be the enpoint pertaining to all of the post handlers
	http.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		paths := strings.Split(r.URL.Path, "/")
		fmt.Println(paths)
		switch len(paths) - 1 {
		case 3:
			fmt.Println("like/dislike post hit")
			if paths[3] == "like" || paths[3] == "dislike" {
				PostsHandlers.HandleLikeDislike(w, r, paths[2])
			} else {
				PostsHandlers.HandleNotFound(w, r, "invalid endpoint for single post")
			}
		case 2:
			id := paths[2]
			if len(id) <= 1 {
				if len(paths[2]) > 1 {
					// if this hits means more than white space was passed
					PostsHandlers.HandleNotFound(w, r, "invalid url path was given")
					return
				}
				fmt.Println("create post was hit")
				PostsHandlers.CreatePost(w, r)
			} else {
				if r.Method == "PUT" {
					fmt.Println("update post hit")
					PostsHandlers.UpdatePost(w, r, id)
				} else if r.Method == "DELETE" {
					fmt.Println("delete post hit")
					PostsHandlers.DeletePost(w, r, id)
				} else if r.Method == "GET" {
					fmt.Println("get post hit")
					PostsHandlers.GetPost(w, r, id)
				} else {
					PostsHandlers.HandleNotFound(w, r, "unsupported method  given to post route")
				}
			}
		default:
			PostsHandlers.HandleNotFound(w, r, "no post endpoint for given url")
		}
	})
	http.ListenAndServe(host+":"+port, nil)
}
