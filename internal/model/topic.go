package model

import (
	"time"

	"github.com/google/uuid"
)

type Topic struct {
	ID        uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name      string      `gorm:"type:varchar(120);uniqueIndex;not null"`
	Slug      string      `gorm:"type:varchar(140);uniqueIndex;not null"`
	Status    string      `gorm:"type:varchar(20);index;default:'approved';not null"` // pending, approved, rejected
	CreatedBy *uuid.UUID  `gorm:"type:uuid;index"`
	CreatedAt time.Time   `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time   `gorm:"not null;default:CURRENT_TIMESTAMP"`
	Users     []UserTopic `gorm:"foreignKey:TopicID;constraint:OnDelete:CASCADE"`
}

type UserTopic struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_user_topic,priority:1;index;not null"`
	TopicID   uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_user_topic,priority:2;index;not null"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Topic     Topic     `gorm:"foreignKey:TopicID;constraint:OnDelete:CASCADE"`
}
