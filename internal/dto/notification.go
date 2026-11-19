package dto

import (
	"time"

	"github.com/google/uuid"
)

type NotificationResponse struct {
	ID             uuid.UUID          `json:"id"`
	Type           string             `json:"type"`
	ActorID        *uuid.UUID         `json:"actor_id,omitempty"`
	ActorName      string             `json:"actor_name,omitempty"`
	ResourceType   string             `json:"resource_type,omitempty"`
	ResourceID     *uuid.UUID         `json:"resource_id,omitempty"`
	Message        string             `json:"message"`
	Title          string             `json:"title,omitempty"`
	Actor          *NotificationActor `json:"actor,omitempty"`
	DiscussionID   *uuid.UUID         `json:"discussionId,omitempty"`
	ProjectID      *uuid.UUID         `json:"projectId,omitempty"`
	RequestID      *uuid.UUID         `json:"requestId,omitempty"`
	Status         string             `json:"status,omitempty"`
	ReadAt         *time.Time         `json:"read_at"`
	IsRead         bool               `json:"is_read"`
	IsReadCamel    bool               `json:"isRead"`
	CreatedAt      time.Time          `json:"created_at"`
	CreatedAtCamel time.Time          `json:"createdAt"`
}

type NotificationActor struct {
	ID          uuid.UUID `json:"id"`
	FullName    string    `json:"fullName"`
	Role        string    `json:"role"`
	Affiliation string    `json:"affiliation"`
	AvatarURL   string    `json:"avatarUrl"`
}
