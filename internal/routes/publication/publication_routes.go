package publicationroutes

import (
	publicationhandler "nalakarsa/internal/handler/publication"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(protected *gin.RouterGroup, handler *publicationhandler.PublicationHandler) {
	publications := protected.Group("/publications")
	publications.GET("/feed", handler.Feed)
	publications.GET("/:publicationId", handler.GetByID)
	publications.POST("/:publicationId/save", handler.Save)
	publications.DELETE("/:publicationId/save", handler.Unsave)
}
