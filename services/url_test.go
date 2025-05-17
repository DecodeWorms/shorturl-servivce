package services

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/net/context"
	"shorturl-service/mocks"
	"shorturl-service/models"
	"testing"
)

// Dummy ID generator
type mockIDGenn struct{}

func (m *mockIDGenn) Generate() string {
	return "mock-id-001"
}

func TestUrlService_CreateShortUrl_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockDataStore(ctrl)
	mockCache := mocks.NewMockRedisCache(ctrl)
	idGen := &mockIDGen{}

	svc := NewUrlService(mockStore, idGen, mockCache)

	req := &models.UrlRequest{
		LongUrl: "https://example.com/page",
	}

	// Setup
	mockStore.EXPECT().GetUser("user123").Return(&models.Users{ID: "user123"}, nil)
	mockStore.EXPECT().GetShortUrl(gomock.Any()).Return(nil, errors.New("not found")).Times(1)
	mockStore.EXPECT().CreateUrl(gomock.Any()).Return(nil)

	shortUrl, err := svc.CreateShortUrl(context.Background(), "user123", req)

	assert.NoError(t, err)
	assert.Len(t, shortUrl, 7)
}

func TestUrlService_GetLongUrl_CacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockDataStore(ctrl)
	mockCache := mocks.NewMockRedisCache(ctrl)
	svc := NewUrlService(mockStore, &mockIDGen{}, mockCache)

	// Simulate cache hit
	mockCache.EXPECT().Get(gomock.Any(), "abc123").Return(&models.UrlCache{
		LongUrl: "https://example.com/page",
	}, nil)

	url, err := svc.GetLongUrl(context.Background(), "abc123")

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/page", url)
}

func TestUrlService_GetLongUrl_CacheMiss_DBHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockDataStore(ctrl)
	mockCache := mocks.NewMockRedisCache(ctrl)
	svc := NewUrlService(mockStore, &mockIDGen{}, mockCache)

	// Cache miss
	mockCache.EXPECT().Get(gomock.Any(), "abc123").Return(nil, errors.New("cache miss"))

	// DB hit
	mockStore.EXPECT().GetShortUrl("abc123").Return(&models.Urls{
		LongUrl:  "https://example.com/page",
		ShortUrl: "abc123",
	}, nil)

	// Set cache
	mockCache.EXPECT().Set(gomock.Any(), "abc123", models.UrlCache{
		LongUrl:  "https://example.com/page",
		ShortUrl: "abc123",
	}).Return(nil)

	url, err := svc.GetLongUrl(context.Background(), "abc123")

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/page", url)
}
