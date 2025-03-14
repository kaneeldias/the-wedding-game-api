package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"the-wedding-game-api/middleware"
	"the-wedding-game-api/models"
	"the-wedding-game-api/types"
)

func GetGallery(c *gin.Context) {
	_, err := middleware.GetCurrentUser(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	images, err := models.GetGalleryImages()
	if err != nil {
		_ = c.Error(err)
		return
	}

	var response types.GalleryResponse
	response.Images = images

	c.IndentedJSON(http.StatusOK, response)
	return
}
