package handler

import "github.com/gin-gonic/gin"

func Route(c *gin.Engine, h *Handler) {
	c.POST("/api/add", h.GetNewCars)
	c.POST("api/all", h.GetData)
	c.DELETE("api/delete/:id", h.DeleteData)
	c.PUT("api/update/:id", h.UpdateData)
}
