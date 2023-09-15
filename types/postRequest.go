package types

// stuct of the data sent when a new user is create or requested
type RequestPost struct {
	UserId string `json:"userId"`
	Desc   string `json:"desc"`
	Image  string `json:"img"`
}

// checkes if the given post in a request has the required values
// to create a post (returns true if post is valid)
func ValidReqestPost(post *RequestPost) bool {
	return (post.UserId != "" && post.Image != "")
}
