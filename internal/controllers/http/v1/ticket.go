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
	}
}

type tickerCreateReq struct {
	Provider string    `json:"provider" example:"Emirates" binding:"required,max=255"`
	FlyFrom  string    `json:"flyFrom" example:"Moscow" binding:"required,max=255"`
	FlyTo    string    `json:"flyTo" example:"Hanoi" binding:"required,max=255"`
	FlyAt    time.Time `json:"flyAt" example:"2022-01-02T15:04:05+03:00" binding:"required"`
	ArriveAt time.Time `json:"arriveAt" example:"2022-01-03T15:04:05+07:00" binding:"required,gtfield=FlyAt"`
}

// @tags Tickets
// @accept json
// @param Ticket body tickerCreateReq true "Ticket request entity"
// @response 201
// @header 201 {string} Location "Return /v1/tickets/:id resource"
// @response 422
// @response 500
// @router /tickets/ [post]
func (g *ticketGroup) create(c *gin.Context) {
	req := tickerCreateReq{}

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

	// TODO: parse header route
	_ = id

	c.Status(http.StatusCreated)
}
