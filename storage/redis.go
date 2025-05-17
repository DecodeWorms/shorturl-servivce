package storage

import (
	"context"
	"encoding/json"
	goredis "github.com/redis/go-redis/v9" // Using goredis as the alias
	"log"
	"shorturl-service/models"
	"time"
)

var _ RedisCache = &RedisStore{}

type RedisStore struct {
	Client *goredis.Client
}

func New(address, password string, db int, ctx context.Context) (RedisCache, error) {
	//Connecting to redis store
	log.Println("Connecting to a redis store...")
	client := goredis.NewClient(&goredis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	//Connection was successful
	log.Println("Connected to redis store successfully..")
	return &RedisStore{
		Client: client,
	}, nil
}

func (r RedisStore) Set(ctx context.Context, key string, url models.UrlCache) error {
	v := models.UrlCache{
		LongUrl:  url.LongUrl,
		ShortUrl: url.ShortUrl,
	}

	u, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = r.Client.Set(ctx, key, string(u), time.Hour*24*7).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r RedisStore) Get(ctx context.Context, key string) (*models.UrlCache, error) {
	userJson, err := r.Client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var storedUser models.UrlCache
	err = json.Unmarshal(userJson, &storedUser)
	if err != nil {
		return nil, err
	}
	return &storedUser, nil
}

func (r RedisStore) Delete(key string) (string, error) {
	return "", nil
}
