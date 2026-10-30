package dto

import (
	"time"

	"github.com/google/uuid"
)

type TopicResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type SuggestTopicRequest struct {
	Name string `json:"name" binding:"required,min=2,max=120"`
}

type UpdateInterestsRequest struct {
	TopicIDs []uuid.UUID `json:"topic_ids" binding:"required,min=3"`
}
