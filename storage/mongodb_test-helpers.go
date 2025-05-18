package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName         = "rescuesupport"
	userCollection = "users"
	aidCollection  = "aids"
	mongoPort      = "27017"
	mongoURI       = "mongodb://localhost:" + mongoPort
)

// setupTestMongoDB connects to MongoDB and initializes the data store
func SetupTestMongoDB(t *testing.T) (*mongo.Client, DataStore) {
	t.Helper()

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	assert.NoError(t, err)

	store, _, err := NewMongo(dbName, mongoURI)
	assert.NoError(t, err)

	return client, store
}
