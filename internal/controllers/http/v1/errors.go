package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/v1adhope/flights/internal/entities"
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
				switch {
				case errors.Is(err, entities.ErrorNothingToChange),
					errors.Is(err, entities.ErrorNothingToDelete),
					errors.Is(err, entities.ErrorNothingFound):
					log.Debug(ginErr, "%s", "StatusNoContent")
					c.AbortWithStatus(http.StatusNoContent)
					return
				case errors.Is(err, entities.ErrorHasAlreadyExists),
					errors.Is(err, entities.ErrorPassengerDoesNotExists):
					log.Debug(ginErr, "%s", "StatusConflict")
					abortWithErrorMsg(c, http.StatusConflict, entities.ErrorHasAlreadyExists.Error())
					return
				}
			}

			log.Error(ginErr, "%s", "StatusInternalServerError")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}
