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
		passengerG.POST("/bound-to-ticket/", group.boundToTicket)
		passengerG.POST("/unbound-from-ticket/", group.unboundToTicket)
		passengerG.GET("/by-ticket-id/:id", group.allByTicketId)

		if gin.Mode() == gin.DebugMode {
			passengerG.GET("/", group.all)
		}
	}
}

type passengerCreateReq struct {
	FirstName  string `json:"firstName" example:"Riley" binding:"required,max=255,names"`
	LastName   string `json:"lastName" example:"Scott" binding:"required,max=255,names"`
	MiddleName string `json:"middleName" example:"Reed" binding:"required,max=255,names"`
}

// @tags Passengers
// @accept json
// @param passenger body passengerCreateReq true "Passenger request entity"
// @response 201 {object} entities.Id
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

	id, err := g.passengerG.CreatePassenger(
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

	c.JSON(http.StatusCreated, id)
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
// @response 422
// @response 500
// @router /passengers/{id} [DELETE]
func (g *passengerGroup) delete(c *gin.Context) {
	params := id{}

	if err := c.ShouldBindUri(&params); err != nil {
		setBindError(c, err)
		return
	}

	err := g.passengerG.DeletePassenger(
		c.Request.Context(),
		entities.Id{params.Value},
	)
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

type passengerBoundingTicketReq struct {
	Id       string `json:"id" example:"uuid" binding:"required,uuid"`
	TicketId string `json:"ticketId" example:"uuid" binding:"required,uuid"`
}

// @tags Passengers
// @accept json
// @param ids body passengerBoundingTicketReq true "Bounding request entity"
// @response 201
// @response 409
// @response 422
// @response 500
// @router /passengers/bound-to-ticket/ [POST]
func (g *passengerGroup) boundToTicket(c *gin.Context) {
	req := passengerBoundingTicketReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		setBindError(c, err)
		return
	}

	err := g.passengerG.BoundToTicket(
		c.Request.Context(),
		entities.Id{req.Id},
		entities.Id{req.TicketId},
	)
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

// @tags Passengers
// @accept json
// @param ids body passengerBoundingTicketReq true "Unbounding request entity"
// @response 200
// @response 204
// @response 422
// @response 500
// @router /passengers/unbound-from-ticket/ [POST]
func (g *passengerGroup) unboundToTicket(c *gin.Context) {
	req := passengerBoundingTicketReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		setBindError(c, err)
		return
	}

	err := g.passengerG.UnboundToTicket(
		c.Request.Context(),
		entities.Id{req.Id},
		entities.Id{req.TicketId},
	)
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// @tags Passengers
// @param id path string true "Ticket id (uuid)"
// @response 200 {array} entities.Passenger
// @response 204
// @response 422
// @response 500
// @router /passengers/by-ticket-id/{id} [GET]
func (g *passengerGroup) allByTicketId(c *gin.Context) {
	params := id{}

	if err := c.ShouldBindUri(&params); err != nil {
		setBindError(c, err)
		return
	}

	passengers, err := g.passengerG.GetPassengersByTicketId(
		c.Request.Context(),
		entities.Id{params.Value},
	)
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.JSON(http.StatusOK, passengers)
}

// @tags Passengers
// @description Support endpoint (not by terms). Avaible only within gin debug
// @response 200
// @response 204
// @response 500
// @router /passengers/ [GET]
func (g *passengerGroup) all(c *gin.Context) {
	passengers, err := g.passengerG.GetPassengers(c.Request.Context())
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.JSON(http.StatusOK, passengers)
}
