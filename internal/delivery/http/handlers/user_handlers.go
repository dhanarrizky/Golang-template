package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/dhanarrizky/Golang-template/internal/usecase/user"
	"github.com/dhanarrizky/Golang-template/pkg/utils"
)

func CreateUserHandler(uc *user.CreateUserUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input user.CreateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		createdUser, err := uc.Execute(c.Request.Context(), input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": utils.UnwrapError(err)})
			return
		}

		c.JSON(http.StatusCreated, createdUser)
	}
}

func GetUserHandler(uc *user.GetUserUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		user, err := uc.Execute(c.Request.Context(), uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func DeleteUserHandler(uc *user.DeleteUserUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		if err := uc.Execute(c.Request.Context(), uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": utils.UnwrapError(err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
	}
}