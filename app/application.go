package app

import (
	"github.com/gerdagi/bookstore_users-api/logger"
	"github.com/gin-gonic/gin"
)

var (
	// router is only available in app package
	router = gin.Default()
)

func StartApplication() {
	mapUrls()
	logger.Info("about to start the application...")
	router.Run(":8081")
}
