package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthController who handle health check
type HealthController struct {
}

// HealthCheck is the response object of the API
type healthCheck struct {
	// API return code
	Code string `json:"code" example:"SUCCESS"`
	// Health status of the service
	Health string `json:"health" example:"RUNNING"`
}

// HealthCheck return the Status of the current app
// @Summary Health check endpoint
// @Description health check endpoint to know if the service is up
// @Tags healthcheck
// @Produce  json
// @Success 200 {object} api.healthCheck
// @Router /health [get]
func (h HealthController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, healthCheck{
		Code:   Success,
		Health: "RUNNING",
	})
}
