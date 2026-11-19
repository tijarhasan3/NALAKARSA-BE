package assethandler

import (
	"net/http"
	"strings"

	"nalakarsa/internal/config"
	"nalakarsa/internal/dto"
	userrepository "nalakarsa/internal/repository/user"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AssetHandler struct {
	cfg      *config.Config
	userRepo userrepository.UserRepository
}

func NewAssetHandler(cfg *config.Config, userRepo userrepository.UserRepository) *AssetHandler {
	return &AssetHandler{cfg: cfg, userRepo: userRepo}
}

func (h *AssetHandler) Upload(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "File is required")
		return
	}
	category := strings.TrimSpace(c.DefaultPostForm("category", "general"))
	path, signedURL, err := utils.UploadAssetToSupabase(userID, category, file, h.cfg)
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "not configured") {
			status = http.StatusServiceUnavailable
		}
		utils.ErrorJSONResponseWithMessage(c, status, err.Error())
		return
	}
	if strings.EqualFold(category, "avatar") || strings.EqualFold(category, "avatars") {
		publicURL := utils.NormalizeSupabasePublicURL(path, h.cfg.SupabaseAssetBucket, h.cfg)
		if err := h.userRepo.UpdateAvatar(userID, publicURL); err != nil {
			utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, "Asset uploaded but profile avatar could not be updated: "+err.Error())
			return
		}
	}
	utils.JSONResponse(c, http.StatusCreated, "Asset uploaded successfully", dto.AssetUploadResponse{
		Path: path, URL: signedURL, Name: file.Filename, MimeType: file.Header.Get("Content-Type"), Size: file.Size, Category: category,
	}, nil)
}
