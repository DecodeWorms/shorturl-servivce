package services

import (
	"fmt"
	"shorturl-service/idgenerator"
	"shorturl-service/models"
	"shorturl-service/storage"
	"time"
)

type UserService struct {
	dataStore storage.DataStore
	idGne     idgenerator.IdGenerator
}

func NewUserService(dataStore storage.DataStore, idGen idgenerator.IdGenerator) UserService {
	return UserService{
		dataStore: dataStore,
		idGne:     idGen,
	}
}

func (u *UserService) SignUp(ID string, req *models.UserRequest) error {
	// Handle the edge case where user already exist
	_, err := u.dataStore.GetUser(ID)
	if err == nil {
		return fmt.Errorf("user already exist %v", err)
	}

	data := &models.Users{
		ID:          u.idGne.Generate(),
		UserName:    req.UserName,
		Email:       req.Email,
		DateCreated: time.Now(),
	}

	// Create new user
	if err := u.dataStore.CreateUser(data); err != nil {
		return err
	}
	return nil
}
