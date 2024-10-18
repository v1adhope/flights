package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setBindError(c *gin.Context, err error) {
	c.Error(err).SetType(gin.ErrorTypeBind)
}

func setAnyError(c *gin.Context, err error) {
	c.Error(err).SetType(gin.ErrorTypeAny)
}

func abortWithErrorMsg(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, gin.H{
		"errMsg": msg,
	})
}

func errorsHandler(log Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, ginErr := range c.Errors {
			err, errType := errors.Unwrap(ginErr.Err), ginErr.Type
			_ = err

			switch errType {
			case gin.ErrorTypeBind:
				log.Debug(ginErr, "%s", "StatusUnprocessableEntity")
				abortWithErrorMsg(c, http.StatusUnprocessableEntity, ginErr.Error())
				return
			case gin.ErrorTypeAny:
			}

			log.Error(ginErr, "%s", "StatusInternalServerError")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}
