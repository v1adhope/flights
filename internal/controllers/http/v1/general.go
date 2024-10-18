package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type id struct {
	Id string `uri:"id" binding:"required,uuid"`
}

func setLocationHeader(c *gin.Context, url, id string) {
	c.Header("location", fmt.Sprintf("%s%s", url, id))
}
