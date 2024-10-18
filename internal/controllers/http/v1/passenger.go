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
		passengerG.PUT("/:id", group.replace)
		passengerG.DELETE("/:id", group.delete)
		passengerG.GET("/", group.all)
	}
}

type passengerCreateReq struct {
	FirstName  string `json:"firstName" example:"Riley" binding:"required,max=255"`
	LastName   string `json:"lastName" example:"Scott" binding:"required,max=255"`
	MiddleName string `json:"middleName" example:"Reed" binding:"required,max=255"`
}

// @tags Passengers
// @accept json
// @param passenger body passengerCreateReq true "Passenger request entity"
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

// @tags Passengers
// @accept json
// @param passenger body passengerCreateReq true "Passenger request entity"
// @param id path string true "Passenger id (uuid)"
// @response 200
// @response 422
// @response 500
// @router /passengers/{id} [PUT]
func (g *passengerGroup) replace(c *gin.Context) {
	params := id{}

	if err := c.ShouldBindUri(&params); err != nil {
		setBindError(c, err)
		return
	}

	req := passengerCreateReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		setBindError(c, err)
		return
	}

	err := g.passengerG.ReplacePassenger(
		c.Request.Context(),
		entities.Passenger{
			Id:         params.Value,
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			MiddleName: req.MiddleName,
		},
	)
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// @tags Passengers
// @param id path string true "Passenger id (uuid)"
// @response 200
// @response 204
// @response 500
// @router /passengers/{id} [DELETE]
func (g *passengerGroup) delete(c *gin.Context) {
	params := id{}

	if err := c.ShouldBindUri(&params); err != nil {
		setBindError(c, err)
		return
	}

	err := g.passengerG.DeletePassenger(c.Request.Context(), params.Value)
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// @tags Passengers
// @response 200
// @response 204
// @response 500
// @router /passengers/ [GET]
func (g *passengerGroup) all(c *gin.Context) {
	passengers, err := g.passengerG.GetAllPassengers(c.Request.Context())
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.JSON(http.StatusOK, passengers)
}
