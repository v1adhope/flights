package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/v1adhope/flights/internal/entities"
)

type passengerGroup struct {
	rg         *gin.RouterGroup
	passengerG PassengerUsecaser
}

func registgerPassengerGroup(group *passengerGroup) {
	passengerG := group.rg.Group("passengers")
	{
		passengerG.POST("/", group.create)
	}
}

type passengerCreateReq struct {
	FirstName  string `json:"firstName" example:"Riley" binding:"required,max=255"`
	LastName   string `json:"lastName" example:"Scott" binding:"required,max=255"`
	MiddleName string `json:"middleName" example:"Reed" binding:"required,max=255"`
}

// @tags Passengers
// @accept json
// @param passanger body passengerCreateReq true "Passanger request entity"
// @response 201
// @header 201 {string} location "Return /v/passenger/"
// @response 204
// @response 422
// @response 500
// @router /passengers/ [POST]
func (g *passengerGroup) create(c *gin.Context) {
	req := passengerCreateReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		setBindError(c, err)
		return
	}

	err := g.passengerG.CreatePassenger(
		c.Request.Context(),
		entities.Passenger{
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			MiddleName: req.MiddleName,
		},
	)
	if err != nil {
		setAnyError(c, err)
		return
	}

	setLocationHeader(c, "/v1/passangers/", "")

	c.Status(http.StatusCreated)
}
