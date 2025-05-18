package main

import (
	"context"
	"log"
	"shorturl-service/config"
	"shorturl-service/serverutil"
)

var c config.Config

func main() {
	c = config.ImportConfig(config.OSSource{})
	// Setup DB connection
	store, _ := serverutil.SetUpDatabase(c.DatabaseURL, c.DatabaseName)
	//Setup redis connection
	var ctx = context.Background()
	redisCache, err := serverutil.SetUpRedisCache(c.RedisAddress, c.RedisPassword, c.RedisDB, ctx)
	if err != nil {
		log.Printf("error connecting to redis store %v", err)
	}
	// Setup service
	userService := serverutil.SetUpUserService(store)
	urlService := serverutil.SetUpUrlService(store, redisCache)
	// Setup server
	userHandler := serverutil.SetUpUserHandler(userService)
	urlHandler := serverutil.SetUpUrlHandler(urlService)
	serverEng := serverutil.SetUpRouter(&userHandler, &urlHandler)
	serverutil.StartServer(serverEng)
}
