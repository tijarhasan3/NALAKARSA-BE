package model

import (
	"time"

	"github.com/google/uuid"
)

// Publication is a published knowledge item displayed by the publication feed.
type Publication struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OwnerID         uuid.UUID  `gorm:"type:uuid;index;not null"`
	PublicationType string     `gorm:"type:varchar(40);index;not null"`
	Title           string     `gorm:"type:varchar(255);not null"`
	Summary         string     `gorm:"type:text;not null"`
	Content         string     `gorm:"type:text;not null"`
	Category        string     `gorm:"type:varchar(100);index;not null"`
	Status          string     `gorm:"type:varchar(20);index;not null;default:'DRAFT'"`
	CoverURL        string     `gorm:"type:varchar(500);default:'';not null"`
	ViewCount       int        `gorm:"not null;default:0"`
	SaveCount       int        `gorm:"not null;default:0"`
	PublishedAt     *time.Time `gorm:"index"`
	CreatedAt       time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`

	Owner         User                      `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`
	Collaborators []PublicationCollaborator `gorm:"foreignKey:PublicationID;constraint:OnDelete:CASCADE"`
}

type PublicationCollaborator struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	PublicationID    uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_publication_collaborator,priority:1;index;not null"`
	UserID           uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_publication_collaborator,priority:2;index;not null"`
	ContributionRole string    `gorm:"type:varchar(100);default:'';not null"`
	IsOwner          bool      `gorm:"not null;default:false"`
	CreatedAt        time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`

	Publication Publication `gorm:"foreignKey:PublicationID;constraint:OnDelete:CASCADE"`
	User        User        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type PublicationSave struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	PublicationID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_publication_save,priority:1;index;not null"`
	UserID        uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_publication_save,priority:2;index;not null"`
	CreatedAt     time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
