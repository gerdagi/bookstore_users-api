package users

import (
	"net/http"
	"strconv"

	"github.com/gerdagi/bookstore_oauth-go/oauth"
	"github.com/gerdagi/bookstore_users-api/domain/users"
	"github.com/gerdagi/bookstore_users-api/services"
	resterr "github.com/gerdagi/bookstore_utils-go/rest_errors"
	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {

	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		c.JSON(err.Status, err)
		return
	}

	if callerId := oauth.GetCallerId(c.Request); callerId == 0 {
		err := resterr.RestError{
			Status:  http.StatusUnauthorized,
			Message: "resource not available",
		}
		c.JSON(err.Status, err)
		return
	}

	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	user, getErr := services.UsersService.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	if oauth.GetCallerId(c.Request) == int64(user.Id) {
		c.JSON(http.StatusOK, user.Marshall(false))
		return
	}

	c.JSON(http.StatusOK, user.Marshall(oauth.IsPublic(c.Request)))
}

func getUserId(userIdParam string) (int64, *resterr.RestError) {
	userId, userErr := strconv.ParseInt(userIdParam, 10, 64)
	if userErr != nil {
		return 0, resterr.NewBadRequestError("user id should be a number")
	}

	return userId, nil
}

func Search(c *gin.Context) {
	status := c.Query("status")

	users, err := services.UsersService.SearchUser(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.Marshall(c.GetHeader("X-Public") == "true")
	}

	c.JSON(http.StatusOK, result)
}

func Delete(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	delErr := services.UsersService.DeleteUser(userId)
	if delErr != nil {
		c.JSON(delErr.Status, delErr)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func Update(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := resterr.NewBadRequestError("invalid JSON body")
		c.JSON(restErr.Status, restErr)
		return
	}

	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	user.Id = int32(userId)

	isPartial := c.Request.Method == http.MethodPatch

	result, updErr := services.UsersService.UpdateUser(isPartial, user)
	if updErr != nil {
		c.JSON(updErr.Status, updErr)
		return
	}

	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Create(c *gin.Context) {
	var user users.User

	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := resterr.NewBadRequestError("invalid JSON body")
		c.JSON(restErr.Status, restErr)
		return
	}

	// We use servise now
	// services/users_service.go
	result, saveErr := services.UsersService.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}

	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Login(c *gin.Context) {
	var request users.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := resterr.NewBadRequestError("invalid JSON data")
		c.JSON(restErr.Status, restErr)
		return
	}

	user, getErr := services.UsersService.LoginUser(request)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-Public") == "true"))
}
