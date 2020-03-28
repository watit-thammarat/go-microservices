package users

import (
	"microservices/bookstore_users-api/domain/users"
	"microservices/bookstore_users-api/services"
	"microservices/bookstore_users-api/utils/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func isPublic(c *gin.Context) bool {
	return c.GetHeader("X-Public") == "true"
}

func getUserID(id string) (int64, *errors.RestErr) {
	userID, userErr := strconv.ParseInt(id, 10, 64)
	if userErr != nil {
		restErr := errors.NewBadRequestError("invalid user id")
		return 0, restErr
	}
	return userID, nil
}

func GetUser(c *gin.Context) {
	userID, userErr := getUserID(c.Param("user_id"))
	if userErr != nil {
		c.JSON(userErr.Status, userErr)
		return
	}
	user, getErr := services.UsersService.GetUser(userID)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusCreated, user.Marshall(isPublic(c)))
}

func CreateUser(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	result, saveErr := services.UsersService.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshall(isPublic(c)))
}

func UpdateUser(c *gin.Context) {
	userID, userErr := getUserID(c.Param("user_id"))
	if userErr != nil {
		c.JSON(userErr.Status, userErr)
		return
	}
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	user.Id = userID
	isPartial := c.Request.Method == http.MethodPatch
	result, updateErr := services.UsersService.UpdateUser(user, isPartial)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(isPublic(c)))
}

func DeleteUser(c *gin.Context) {
	userID, userErr := getUserID(c.Param("user_id"))
	if userErr != nil {
		c.JSON(userErr.Status, userErr)
		return
	}
	if err := services.UsersService.DeleteUser(userID); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func Search(c *gin.Context) {
	status := c.Query("status")
	users, err := services.UsersService.SearchUser(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, users.Marshall(isPublic(c)))
}
