package storage

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"shorturl-service/models"
	"testing"
)

func TestUserCreationAndRetrieval(t *testing.T) {
	client, store := SetupTestMongoDB(t)
	defer client.Disconnect(context.Background())

	user := &models.Users{
		ID:       "user123",
		UserName: "Jane Doe",
		Email:    "jane@example.com",
	}

	err := store.CreateUser(user)
	assert.NoError(t, err)

	result, err := store.GetUser("user123")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Jane Doe", result.UserName)
	assert.Equal(t, "jane@example.com", result.Email)
}

func TestUrlCreationAndRetrieval(t *testing.T) {
	client, store := SetupTestMongoDB(t)
	defer client.Disconnect(context.Background())

	url := &models.Urls{
		ID:       "url-123",
		LongUrl:  "https://example.com/page",
		ShortUrl: "abc1234",
		UserID:   "user123",
	}

	err := store.CreateUrl(url)
	assert.NoError(t, err)

	result, err := store.GetOriginalUrl("abc1234")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "https://example.com/page", result.LongUrl)
	assert.Equal(t, "abc1234", result.ShortUrl)
}
