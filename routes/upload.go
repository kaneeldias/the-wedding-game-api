package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"the-wedding-game-api/middleware"
	"the-wedding-game-api/types"
	"the-wedding-game-api/utils"
)

func HandleImageUpload(c *gin.Context) {
	err := middleware.CheckIsAdmin(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		_ = c.Error(err)
		return
	}

	url, err := utils.UploadFile(file)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := types.UploadResponse{
		Url: url,
	}
	c.IndentedJSON(http.StatusOK, response)
}
