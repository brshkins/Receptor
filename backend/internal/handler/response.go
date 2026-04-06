package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONData(c *gin.Context, code int, data any) {
	c.JSON(code, gin.H{"data": data})
}

func JSONOK(c *gin.Context, data any) {
	JSONData(c, http.StatusOK, data)
}

func JSONError(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
}
