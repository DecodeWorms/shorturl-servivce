package models

import "time"

type Users struct {
	ID          string    `json:"id" bson:"id"`
	UserName    string    `json:"user_name" bson:"user_name"`
	Email       string    `json:"email" bson:"email"`
	DateCreated time.Time `json:"date_created" bson:"date_created"`
}
