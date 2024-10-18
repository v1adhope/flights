package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/v1adhope/flights/internal/entities"
)

type ticketGroup struct {
	rg      *gin.RouterGroup
	ticketU TicketUsecaser
}

func registerTicketGroup(group *ticketGroup) {
	ticketG := group.rg.Group("/tickets")
	{
		ticketG.POST("/", group.create)
		ticketG.PUT("/:id", group.replace)
	}
}

type ticketCreateReq struct {
	Provider string    `json:"provider" example:"Emirates" binding:"required,max=255"`
	FlyFrom  string    `json:"flyFrom" example:"Moscow" binding:"required,max=255"`
	FlyTo    string    `json:"flyTo" example:"Hanoi" binding:"required,max=255"`
	FlyAt    time.Time `json:"flyAt" example:"2022-01-02T15:04:05+03:00" binding:"required"`
	ArriveAt time.Time `json:"arriveAt" example:"2022-01-03T15:04:05+07:00" binding:"required,gtfield=FlyAt"`
}

// @tags Tickets
// @accept json
// @param ticket body ticketCreateReq true "Ticket request entity"
// @response 201
// @header 201 {string} Location "Return /v1/tickets/:id resource"
// @response 422
// @response 500
// @router /tickets/ [post]
func (g *ticketGroup) create(c *gin.Context) {
	req := ticketCreateReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		setBindError(c, err)
		return
	}

	id, err := g.ticketU.Create(c.Request.Context(), entities.Ticket{
		FlyFrom:  req.FlyFrom,
		FlyTo:    req.FlyTo,
		Provider: req.Provider,
		FlyAt:    req.FlyAt,
		ArriveAt: req.ArriveAt,
	})
	if err != nil {
		setAnyError(c, err)
		return
	}

	setLocationHeader(c, "/v1/tickets/", id)

	c.Status(http.StatusCreated)
}

// @tags Tickets
// @accept json
// @param ticket body ticketCreateReq true "Ticket request entity"
// @param id path string true "Ticket Id (uuid)"
// @response 200
// @response 422
// @response 500
// @router /tickets/{id} [PUT]
func (g *ticketGroup) replace(c *gin.Context) {
	params := id{}

	if err := c.ShouldBindUri(&params); err != nil {
		setBindError(c, err)
		return
	}

	req := ticketCreateReq{}

	if err := c.ShouldBind(&req); err != nil {
		setBindError(c, err)
		return
	}

	err := g.ticketU.Replace(c.Request.Context(), entities.Ticket{
		Id:       params.Id,
		FlyFrom:  req.FlyFrom,
		FlyTo:    req.FlyTo,
		Provider: req.Provider,
		FlyAt:    req.FlyAt,
		ArriveAt: req.ArriveAt,
	})
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
