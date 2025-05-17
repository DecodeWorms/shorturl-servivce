package models

type UserRequest struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

type UrlRequest struct {
	LongUrl string `json:"long_url"`
}
