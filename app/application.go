package app

import "github.com/gin-gonic/gin"

var (
	// router is only available in app package
	router = gin.Default()
)

func StartApplication() {
	mapUrls()
	router.Run(":8080")
}
