package assetroutes

import (
	assethandler "nalakarsa/internal/handler/asset"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(protected *gin.RouterGroup, h *assethandler.AssetHandler) {
	protected.POST("/assets", h.Upload)
}
