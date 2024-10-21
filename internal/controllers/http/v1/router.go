package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("names", names)
		v.RegisterStructValidation(ticketCreateReqStructLevelValidation, ticketCreateReq{})
		v.RegisterStructValidation(reportByPassengerIdForPeriodQueryStructLevelValidation, reportByPassengerIdForPeriodQuery{})
	}

	rg := r.Handler.Group("/v1")
	rg.Use(errorsHandler(r.Log))
	{
		rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

		registerTicketGroup(&ticketGroup{rg, r.Usecases})
		registgerPassengerGroup(&passengerGroup{rg, r.Usecases})
		registerDocumentGroup(&documentGroup{rg, r.Usecases})
		registerReportGroup(&reportGroup{rg, r.Usecases})
	}
}

func SetMode(m string) {
	if m == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
}
