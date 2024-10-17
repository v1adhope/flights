package v1

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "github.com/v1adhope/flights/docs"
	"github.com/v1adhope/flights/internal/usecases"
	"github.com/v1adhope/flights/pkg/logger"
)

type Router struct {
	Handler  *gin.Engine
	Usecases *usecases.Usecases
	Log      *logger.Log
}

func Register(r *Router) {
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Title = "Flights API"

	rg := r.Handler.Group("/v1")
	rg.Use(errorsHandler(r.Log))
	{
		rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

		registerTicketGroup(&ticketGroup{rg, r.Usecases})
	}
}

func SetMode(m string) {
	if m == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
}
