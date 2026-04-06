package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JSONData writes { "data": data } with the given status (backend_spec §8).
func JSONData(c *gin.Context, code int, data any) {
	c.JSON(code, gin.H{"data": data})
}

// JSONOK writes 200 and { "data": data }.
func JSONOK(c *gin.Context, data any) {
	JSONData(c, http.StatusOK, data)
}

// JSONError writes { "error": message }.
func JSONError(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
}
