package serverutil

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shorturl-service/config"
	"shorturl-service/handler"
	"shorturl-service/idgenerator"
	"shorturl-service/services"
	"shorturl-service/storage"
	"syscall"
	"time"
)

func SetUpDatabase(url, name string) (storage.DataStore, *mongo.Client) {
	repo, client, err := storage.NewMongo(name, url)
	if err != nil {
		log.Fatalf("Error failed to open MongoDB: %v", err)
	}
	return repo, client
}

func SetUpRedisCache(add, pass string, db int, ctx context.Context) (storage.RedisCache, error) {
	return storage.New(add, pass, db, ctx)
}

func SetUpUserService(store storage.DataStore) services.UserService {
	return services.NewUserService(store, idgenerator.New())
}

func SetUpUrlService(store storage.DataStore, redis storage.RedisCache) services.UrlService {
	return services.NewUrlService(store, idgenerator.New(), redis)
}

func SetUpUserHandler(userService services.UserService) handler.UserHandler {
	return handler.NewUserHandler(userService)
}

func SetUpUrlHandler(urlService services.UrlService) handler.UrlHandler {
	return handler.NewUrlHandler(urlService)
}

func SetUpRouter(user *handler.UserHandler, url *handler.UrlHandler) *gin.Engine {
	router := gin.Default()

	// Add Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all (use specific domain in prod)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api/v1")
	{
		// API endpoint(s) for accessing a user resources
		api.POST("/user", user.CreateUser())

		//API endpoints for accessing url resources
		api.POST("/url", url.GenerateShortUrl())
		api.GET("/url", url.Redirect())
	}

	return router
}

func StartServer(router *gin.Engine) {
	//var c config.Config
	var c = config.ImportConfig(config.OSSource{})
	interruptHandler := make(chan os.Signal, 1)
	signal.Notify(interruptHandler, syscall.SIGTERM, syscall.SIGINT)

	addr := fmt.Sprintf(":%s", c.ServicePort)
	go func(addr string) {
		log.Printf("ShortenUrl.sv API service running on %v. Environment=%s", addr, c.AppEnv)
		if err := http.ListenAndServe(addr, router); err != nil {
			log.Printf("Error starting server: %v", err)
		}
	}(addr)

	<-interruptHandler
}
