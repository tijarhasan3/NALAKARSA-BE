package discussionroutes

import (
	"nalakarsa/internal/config"
	"nalakarsa/internal/handler/discussion"
	"nalakarsa/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public *gin.RouterGroup, protected *gin.RouterGroup, h *discussionhandler.DiscussionHandler, cfg *config.Config) {
	discussions := public.Group("/discussions")
	discussions.GET("", middleware.OptionalAuthMiddleware(cfg), h.List)
	discussions.GET("/:id", h.GetByID)
	discussions.GET("/:id/replies", h.GetReplies)
	discussions.GET("/:id/collaboration", h.GetByID)

	protected.POST("/discussions", h.Create)
	protected.PATCH("/discussions/:id", h.Update)
	protected.DELETE("/discussions/:id", h.Delete)
	protected.POST("/discussions/:id/collaboration", h.MarkCollaboration)
	protected.POST("/discussions/:id/replies", h.AddReply)
	protected.DELETE("/discussions/:id/replies/:reply_id", h.DeleteReply)
	protected.POST("/discussions/:id/votes", h.Vote)
	protected.DELETE("/discussions/:id/votes", h.Unvote)
	protected.POST("/discussions/:id/upvote", h.Vote)
	protected.DELETE("/discussions/:id/upvote", h.Unvote)
}
