package topichandler

import (
	"net/http"

	"nalakarsa/internal/dto"
	topicservice "nalakarsa/internal/service/topic"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TopicHandler struct{ service topicservice.TopicService }

func NewTopicHandler(service topicservice.TopicService) *TopicHandler {
	return &TopicHandler{service: service}
}

func (h *TopicHandler) List(c *gin.Context) {
	data, err := h.service.List(c.Query("search"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(c, http.StatusOK, "Topics retrieved successfully", data, nil)
}

func (h *TopicHandler) Suggest(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	var req dto.SuggestTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}
	data, err := h.service.Suggest(userID, req)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusConflict, err.Error())
		return
	}
	utils.JSONResponse(c, http.StatusCreated, "Topic suggestion submitted successfully", data, nil)
}

func (h *TopicHandler) ListMine(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	data, err := h.service.ListUserTopics(userID)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(c, http.StatusOK, "User topics retrieved successfully", data, nil)
}

func (h *TopicHandler) UpdateMine(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	var req dto.UpdateInterestsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}
	data, err := h.service.UpdateUserTopics(userID, req)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSONResponse(c, http.StatusOK, "User topics updated successfully", data, nil)
}
