package dto

import (
	"time"

	"github.com/google/uuid"
)

type PublicationCollaboratorResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Role        string    `json:"role"`
	Affiliation string    `json:"affiliation"`
	AvatarURL   *string   `json:"avatar_url"`
}

type PublicationResponse struct {
	ID              uuid.UUID                         `json:"id"`
	PublicationType string                            `json:"publication_type"`
	Title           string                            `json:"title"`
	Summary         string                            `json:"summary"`
	Content         string                            `json:"content"`
	Category        string                            `json:"category"`
	Status          string                            `json:"status"`
	CoverURL        string                            `json:"cover_url"`
	ViewCount       int                               `json:"view_count"`
	SaveCount       int                               `json:"save_count"`
	CreatedAt       time.Time                         `json:"created_at"`
	UpdatedAt       time.Time                         `json:"updated_at"`
	IsSaved         bool                              `json:"is_saved"`
	Collaborators   []PublicationCollaboratorResponse `json:"collaborators"`
}

type PublicationFeedMeta struct {
	NextCursor *string `json:"next_cursor"`
	HasMore    bool    `json:"has_more"`
	Limit      int     `json:"limit"`
}

type PublicationFeedResponse struct {
	Data []PublicationResponse `json:"data"`
	Meta PublicationFeedMeta   `json:"meta"`
}

type PublicationSaveResponse struct {
	PublicationID uuid.UUID `json:"publication_id"`
	IsSaved       bool      `json:"is_saved"`
	SaveCount     int       `json:"save_count"`
}
