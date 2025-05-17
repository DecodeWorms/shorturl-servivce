package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shorturl-service/models"
	"shorturl-service/services"
)

type UrlHandler struct {
	urlService services.UrlService
}

func NewUrlHandler(urlService services.UrlService) UrlHandler {
	return UrlHandler{
		urlService: urlService,
	}
}

func (u *UrlHandler) GenerateShortUrl() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get ID from query parameter
		ID := ctx.Query("id")
		if ID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
			return
		}

		var urlReq models.UrlRequest
		if err := ctx.ShouldBindJSON(&urlReq); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "error bad user request" + err.Error()})
			return
		}

		shortUrl, err := u.urlService.CreateShortUrl(ctx, ID, &urlReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error generating short url , pls try again" + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success":   "Short url generated",
			"short url": shortUrl,
		})
	}
}

func (u *UrlHandler) Redirect() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		shortUrl := ctx.Query("short_url")
		if shortUrl == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "error user bad request, pls provide correct short url"})
			return
		}
		longUrl, err := u.urlService.GetLongUrl(ctx, shortUrl)
		if err != nil || longUrl == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "long url is not available"})
			return
		}

		// Perform a 302 redirect to the long URL
		ctx.Redirect(http.StatusFound, longUrl)
	}
}
