package handler

import (
	"Proxy/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler struct {
	usecase usecase.UsecaseI
}

func NewHandler(usecase usecase.UsecaseI) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (handler *Handler) SetupRoutes() *gin.Engine {
	router := gin.New()
	api := router.Group("/")
	{
		api.GET("/requests", handler.AllRequests)
		api.GET("/requests/:id", handler.GetRequest)
		api.GET("/repeat/:id", handler.RepeatRequest)
	}

	return router
}

// AllRequests /requests
func (h *Handler) AllRequests(c *gin.Context) {
	requests, err := h.usecase.AllRequests()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, requests)
}

// GetRequest /requests/{id}
func (h *Handler) GetRequest(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, "no such request")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	request, err := h.usecase.GetRequest(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, request)
}

// RepeatRequest /repeat/{id}
func (h *Handler) RepeatRequest(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, "no such request")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	request, err := h.usecase.RepeatRequest(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, request)
}
