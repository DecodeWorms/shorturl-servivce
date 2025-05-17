package services

import (
	"context"
	"fmt"
	"log"
	"shorturl-service/idgenerator"
	"shorturl-service/models"
	"shorturl-service/storage"
	"shorturl-service/urlshortener"
)

type UrlService struct {
	store storage.DataStore
	idGen idgenerator.IdGenerator
	cache storage.RedisCache
}

func NewUrlService(store storage.DataStore, idGen idgenerator.IdGenerator, cache storage.RedisCache) UrlService {
	return UrlService{
		store: store,
		idGen: idGen,
		cache: cache,
	}

}

func (u *UrlService) CreateShortUrl(ctx context.Context, userID string, data *models.UrlRequest) (string, error) {
	// Handle the edge case where user is not registered
	_, err := u.store.GetUser(userID)
	if err != nil {
		return "", fmt.Errorf("user record not found %v", err)
	}

	// Avoid Collision by checking if the short url already exist
	var newShortUrl string
	const maxAttempt = 7

	// generate a new short url, try 7 times
	for i := 0; i < maxAttempt; i++ {
		shortUrl := urlshortener.GenerateShortCode()
		url, err := u.store.GetShortUrl(shortUrl)
		if err != nil && url == nil {
			newShortUrl = shortUrl
			break
		}
	}

	if newShortUrl == "" {
		return "", fmt.Errorf("could not generate a unique url after %d attempts", maxAttempt)
	}

	rec := &models.Urls{
		ID:       u.idGen.Generate(),
		LongUrl:  data.LongUrl,
		ShortUrl: newShortUrl,
		UserID:   userID,
	}
	if err := u.store.CreateUrl(rec); err != nil {
		return "", fmt.Errorf("error in persisting the record to DB %v", err)
	}
	return newShortUrl, nil
}

func (u *UrlService) GetLongUrl(ctx context.Context, url string) (string, error) {
	//Check the cache first
	lurl, err := u.cache.Get(ctx, url)
	if err == nil && lurl != nil {
		return lurl.LongUrl, nil
	}

	// Hit the DB if cache miss
	urlRec, err := u.store.GetShortUrl(url)
	if err != nil {
		return "", fmt.Errorf("urls are not available")
	}

	// Generate long url for redirection
	var longUrl string
	if urlRec != nil {
		longUrl = urlRec.LongUrl
	}
	cachedData := models.UrlCache{
		LongUrl:  longUrl,
		ShortUrl: url,
	}
	if err := u.cache.Set(ctx, url, cachedData); err != nil {
		log.Printf("error saving entry %v", err)
	}
	return longUrl, nil
}
