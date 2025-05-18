package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"shorturl-service/models"
	"time"
)

const (
	UserCollection = "users"
	UrlCollection  = "urls"
)

type Mongodb struct {
	Client       *mongo.Client
	databaseName string
}

func (m *Mongodb) GetShortUrl(url string) (*models.Urls, error) {
	filter := bson.M{"short_url": url}
	var urlRec *models.Urls
	err := m.col(UrlCollection).FindOne(context.Background(), filter).Decode(&urlRec)
	return urlRec, err
}

func NewMongo(address, url string) (DataStore, *mongo.Client, error) {
	log.Println("Connecting to Mongodb store")

	//Config the datastore environment
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//Connect to the mongodb client
	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, nil, err
	}

	//Ping the database to know if it was connected
	if err := cli.Ping(ctx, readpref.Primary()); err != nil {
		return nil, nil, err
	}

	log.Println("Connected to Mongodb successfully")
	return &Mongodb{Client: cli, databaseName: address}, cli, nil
}

func (m *Mongodb) CreateUser(data *models.Users) error {
	if _, err := m.col(UserCollection).InsertOne(context.Background(), data); err != nil {
		return err
	}
	return nil
}

func (m *Mongodb) GetUser(ID string) (*models.Users, error) {
	filter := bson.M{"id": ID}
	var user *models.Users
	err := m.col(UserCollection).FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *Mongodb) col(collectionName string) *mongo.Collection {
	return m.Client.Database(m.databaseName).Collection(collectionName)
}

func (m *Mongodb) CreateUrl(url *models.Urls) error {
	_, err := m.col(UrlCollection).InsertOne(context.Background(), url)
	return err
}

func (m *Mongodb) GetOriginalUrl(url string) (*models.Urls, error) {
	filter := bson.M{"short_url": url}
	var longUrls *models.Urls
	err := m.col(UrlCollection).FindOne(context.Background(), filter).Decode(&longUrls)
	if err != nil {
		return nil, err
	}
	return longUrls, nil
}

var _ DataStore = &Mongodb{}
