package services

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"shorturl-service/mocks"
	"shorturl-service/models"
	"testing"
)

// Dummy ID generator
type mockIDGen struct{}

func (m *mockIDGen) Generate() string {
	return "test-id-001"
}

func TestUserService_SignUp_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockDataStore(ctrl)
	mockIDGen := &mockIDGen{}
	service := NewUserService(mockStore, mockIDGen)

	// Mock input
	userID := "abc123"
	req := &models.UserRequest{
		UserName: "john",
		Email:    "john@example.com",
	}

	// Expectations
	mockStore.EXPECT().GetUser(userID).Return(nil, errors.New("not found"))
	mockStore.EXPECT().CreateUser(gomock.Any()).Return(nil)

	// Execute
	err := service.SignUp(userID, req)

	// Assert
	assert.NoError(t, err)
}

func TestUserService_SignUp_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockDataStore(ctrl)
	mockIDGen := &mockIDGen{}
	service := NewUserService(mockStore, mockIDGen)

	// Mock input
	userID := "abc123"
	req := &models.UserRequest{
		UserName: "john",
		Email:    "john@example.com",
	}

	// Mock existing user
	existing := &models.Users{
		ID:       "abc123",
		UserName: "john",
		Email:    "john@example.com",
	}

	// Expectations
	mockStore.EXPECT().GetUser(userID).Return(existing, nil)

	// Execute
	err := service.SignUp(userID, req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user already exist")
}
