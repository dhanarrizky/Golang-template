package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/dhanarrizky/Golang-template/pkg/responses"
)

func ResponseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// ERROR RESPONSE
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if appErr, ok := err.(*errors.AppError); ok {
				c.JSON(appErr.Status, response.APIResponse{
					Success: false,
					Code:    appErr.Code,
					Message: appErr.Message,
					Errors:  appErr.Errors,
				})
				return
			}

			// fallback 500
			c.JSON(http.StatusInternalServerError, response.APIResponse{
				Success: false,
				Code:    "INTERNAL_ERROR",
				Message: "Internal server error",
			})
			return
		}

		// SUCCESS RESPONSE
		if val, exists := c.Get("response"); exists {
			res := val.(*response.Success)

			c.JSON(res.Status, response.APIResponse{
				Success: true,
				Code:    res.Code,
				Message: res.Message,
				Data:    res.Data,
			})
		}
	}
}


// example how to use
// func GetUser(c *gin.Context) {
// 	user, err := service.FindUser(c.Param("id"))
// 	if err != nil {
// 		c.Error(errors.NotFound("USER_NOT_FOUND", "User tidak ditemukan"))
// 		return
// 	}

// 	c.Set("response", response.OK(
// 		"USER_FOUND",
// 		"OK",
// 		user,
// 	))
// }
