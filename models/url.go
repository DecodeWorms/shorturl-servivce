package models

type Urls struct {
	ID       string `json:"id" bson:"id"`
	ShortUrl string `json:"short_url" bson:"short_url"`
	LongUrl  string `json:"long_url" bson:"long_url"`
	UserID   string `json:"user_id" bson:"user_id"`
}
