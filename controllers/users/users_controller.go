package users

import (
	"net/http"
	"strconv"

	"github.com/gerdagi/bookstore_users-api/domain/users"
	"github.com/gerdagi/bookstore_users-api/services"
	resterr "github.com/gerdagi/bookstore_users-api/utils/errors"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user users.User

	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := resterr.NewBadRequestError("invalid JSON body")
		c.JSON(restErr.Status, restErr)
		return
	}

	// We use servise now
	// services/users_service.go
	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func GetUser(c *gin.Context) {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		restErr := resterr.NewBadRequestError("user id should be a number")
		c.JSON(restErr.Status, restErr)
	}

	user, getErr := services.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	c.JSON(http.StatusOK, user)
}
