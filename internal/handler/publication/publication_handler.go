package publicationhandler

import (
	"net/http"
	"strconv"

	publicationservice "nalakarsa/internal/service/publication"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PublicationHandler struct {
	service publicationservice.PublicationService
}

func NewPublicationHandler(service publicationservice.PublicationService) *PublicationHandler {
	return &PublicationHandler{service: service}
}

func (h *PublicationHandler) Feed(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	limit := 12
	if raw := c.Query("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 {
			utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "limit must be a positive integer")
			return
		}
		limit = parsed
	}
	result, err := h.service.Feed(userID, c.DefaultQuery("scope", "for_you"), c.Query("category"), c.Query("search"), c.Query("cursor"), limit)
	if err != nil {
		if err == publicationservice.ErrInvalidScope || err.Error() == "invalid cursor" {
			utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *PublicationHandler) GetByID(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	publicationID, err := uuid.Parse(c.Param("publicationId"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "invalid publication ID format")
		return
	}
	result, err := h.service.GetByID(userID, publicationID)
	if err != nil {
		status := http.StatusInternalServerError
		if publicationservice.IsNotFound(err) {
			status = http.StatusNotFound
		}
		utils.ErrorJSONResponseWithMessage(c, status, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *PublicationHandler) Save(c *gin.Context) {
	h.save(c, true)
}

func (h *PublicationHandler) Unsave(c *gin.Context) {
	h.save(c, false)
}

func (h *PublicationHandler) save(c *gin.Context, save bool) {
	userID := c.MustGet("user_id").(uuid.UUID)
	publicationID, err := uuid.Parse(c.Param("publicationId"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "invalid publication ID format")
		return
	}
	var result interface{}
	if save {
		result, err = h.service.Save(userID, publicationID)
	} else {
		result, err = h.service.Unsave(userID, publicationID)
	}
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}
