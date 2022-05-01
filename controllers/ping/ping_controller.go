package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	// c.String(http.StatusOK, "pong")
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
