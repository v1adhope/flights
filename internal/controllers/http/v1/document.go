package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/v1adhope/flights/internal/entities"
)

type documentGroup struct {
	rg        *gin.RouterGroup
	documentU DocumentUsecaser
}

func registerDocumentGroup(group *documentGroup) {
	documentG := group.rg.Group("/documents")
	{
		documentG.POST("/", group.create)
		documentG.PUT("/:id", group.replace)
		documentG.DELETE("/:id", group.delete)
		documentG.GET("/by-passenger/:id", group.allByPassengerId)
	}
}

type documentCreateReq struct {
	Type        string `json:"type" example:"Passport" binding:"required,oneof='Passport' 'Id card' 'International passport'"`
	Number      string `json:"number" example:"5555444444" binding:"required,max=255"`
	PassengerId string `json:"passengerId" example:"uuid" binding:"required,uuid"`
}

// @tags Documents
// @description One of Passport, Id card, International passport
// @accept json
// @param document body documentCreateReq true "Document request entity"
// @response 201 {object} entities.Id
// @response 409
// @response 422
// @response 500
// @router /documents/ [POST]
func (g *documentGroup) create(c *gin.Context) {
	req := documentCreateReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		setBindError(c, err)
		return
	}

	id, err := g.documentU.CreateDocument(
		c.Request.Context(),
		entities.Document{
			Type:        req.Type,
			Number:      req.Number,
			PassengerId: req.PassengerId,
		},
	)
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.JSON(http.StatusCreated, entities.Id{id})
}

// @tags Documents
// @accept json
// @param document body documentCreateReq true "Document request entity"
// @param id path string true "Document id (uuid)"
// @response 200
// @response 409
// @response 422
// @response 500
// @router /documents/{id} [PUT]
func (g *documentGroup) replace(c *gin.Context) {
	params := id{}

	if err := c.ShouldBindUri(&params); err != nil {
		setBindError(c, err)
		return
	}

	req := documentCreateReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		setBindError(c, err)
		return
	}

	err := g.documentU.ReplaceDocument(
		c.Request.Context(),
		entities.Document{
			Id:          params.Value,
			Type:        req.Type,
			Number:      req.Number,
			PassengerId: req.PassengerId,
		},
	)
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// @tags Documents
// @param id path string true "Document id (uuid)"
// @response 200
// @response 204
// @response 422
// @response 500
// @router /documents/{id} [DELETE]
func (g *documentGroup) delete(c *gin.Context) {
	params := id{}

	if err := c.ShouldBindUri(&params); err != nil {
		setBindError(c, err)
		return
	}

	err := g.documentU.DeleteDocument(
		c.Request.Context(),
		entities.Id{params.Value},
	)
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// @tags Documents
// @param id path string true "Passenger id (uuid)"
// @response 200 {array} entities.Document
// @response 204
// @response 422
// @response 500
// @router /documents/by-passenger/{id} [GET]
func (g *documentGroup) allByPassengerId(c *gin.Context) {
	params := id{}

	if err := c.ShouldBindUri(&params); err != nil {
		setBindError(c, err)
		return
	}

	documents, err := g.documentU.GetDocumentsByPassengerId(
		c.Request.Context(),
		entities.Id{params.Value},
	)
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.JSON(http.StatusOK, documents)
}
