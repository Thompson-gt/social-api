package types

// stuct of the data sent when a new user is create or requested
type AuthUserRequest struct {
	UserId      string `json:"userId"`
	UserName    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	AdminSecret string `json:"adminSecret"`
}

// checks if the given requestUser is valid
func ValidAuthUser(user *AuthUserRequest) bool {
	if user.Email == "" || user.Password == "" || user.UserName == "" {
		return false
	}
	return true
}
