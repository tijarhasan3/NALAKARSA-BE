package topicroutes

import (
	topichandler "nalakarsa/internal/handler/topic"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public *gin.RouterGroup, protected *gin.RouterGroup, h *topichandler.TopicHandler) {
	public.GET("/topics", h.List)
	protected.GET("/profile/interests", h.ListMine)
	protected.PATCH("/profile/interests", h.UpdateMine)
	protected.POST("/topics/suggestions", h.Suggest)
}
