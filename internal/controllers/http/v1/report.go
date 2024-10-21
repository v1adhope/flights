package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/v1adhope/flights/internal/entities"
)

type reportGroup struct {
	rg      *gin.RouterGroup
	reportU ReportUsecaser
}

func registerReportGroup(group *reportGroup) {
	reportG := group.rg.Group("/reports")
	{
		reportG.GET("/by-passenger-id-for-period/:id", group.byPassengerIdForPeriod)
	}
}

type reportByPassengerIdForPeriodQuery struct {
	From string `form:"from" binding:"required"`
	To   string `form:"to" binding:"required"`
}

// @tags Reports
// @param id path string true "Passenger id (uuid)"
// @param from query string true "Perion start value" format(rfc3339Time)
// @param to query string true "Perion end value" format(rfc3339Time)
// @response 200 {array} entities.ReportRowByPassengerForPeriod
// @response 204
// @response 422
// @response 500
// @router /reports/by-passenger-id-for-period/{id} [GET]
func (g *reportGroup) byPassengerIdForPeriod(c *gin.Context) {
	params := id{}

	if err := c.ShouldBindUri(&params); err != nil {
		setBindError(c, err)
		return
	}

	query := reportByPassengerIdForPeriodQuery{}

	if err := c.ShouldBind(&query); err != nil {
		setBindError(c, err)
		return
	}

	reportRows, err := g.reportU.GetRowsByPassengerIdForPeriod(
		c.Request.Context(),
		entities.Id{params.Value},
		entities.PeriodFilter{
			From: query.From,
			To:   query.To,
		},
	)
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.JSON(http.StatusOK, reportRows)
}
