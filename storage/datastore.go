package storage

import (
	"golang.org/x/net/context"
	"shorturl-service/models"
)

//go:generate mockgen -source=datastore.go -destination=../mocks/datastore_mock.go -package=mocks
type DataStore interface {
	CreateUser(data *models.Users) error
	GetUser(ID string) (*models.Users, error)
	CreateUrl(url *models.Urls) error
	GetOriginalUrl(url string) (*models.Urls, error)
	GetShortUrl(url string) (*models.Urls, error)
}

type RedisCache interface {
	Set(ctx context.Context, key string, value models.UrlCache) error
	Get(ctx context.Context, key string) (*models.UrlCache, error)
	Delete(key string) (string, error)
}
