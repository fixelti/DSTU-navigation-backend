package pathBuilding

import (
	"navigation/internal/logging"
	"navigation/internal/transport/rest/handlers"

	"github.com/gin-gonic/gin"
)

const pathBuildingURL = "/path-building"

type handler struct {
	logger *logging.Logger
}

func NewHandler(logger *logging.Logger) handlers.Handler {
	return &handler{logger: logger}
}

func (h *handler) Register(router *gin.RouterGroup) {
	pathBuilding := router.Group(pathBuildingURL)
	pathBuilding.GET("", h.pathBuilding)
}

func (h *handler) pathBuilding(c *gin.Context) {
	return nil
}
