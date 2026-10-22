package publicationrepository

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"nalakarsa/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Cursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
}

type PublicationRepository interface {
	ListPublished(userID uuid.UUID, scope, category, search, cursor string, limit int) ([]model.Publication, bool, string, error)
	GetPublishedByID(id uuid.UUID) (*model.Publication, error)
	IsSaved(publicationID, userID uuid.UUID) (bool, error)
	Save(publicationID, userID uuid.UUID) (int, error)
	Unsave(publicationID, userID uuid.UUID) (int, error)
}

type publicationRepository struct{ db *gorm.DB }

func NewPublicationRepository(db *gorm.DB) PublicationRepository {
	return &publicationRepository{db: db}
}

func EncodeCursor(createdAt time.Time, id uuid.UUID) string {
	b, _ := json.Marshal(Cursor{CreatedAt: createdAt, ID: id})
	return base64.RawURLEncoding.EncodeToString(b)
}

func DecodeCursor(value string) (Cursor, error) {
	if strings.TrimSpace(value) == "" {
		return Cursor{}, nil
	}
	b, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return Cursor{}, errors.New("invalid cursor")
	}
	var cursor Cursor
	if err := json.Unmarshal(b, &cursor); err != nil || cursor.ID == uuid.Nil || cursor.CreatedAt.IsZero() {
		return Cursor{}, errors.New("invalid cursor")
	}
	return cursor, nil
}

func (r *publicationRepository) ListPublished(userID uuid.UUID, scope, category, search, cursorValue string, limit int) ([]model.Publication, bool, string, error) {
	cursor, err := DecodeCursor(cursorValue)
	if err != nil {
		return nil, false, "", err
	}
	query := r.db.Model(&model.Publication{}).
		Where("publications.status = ?", "PUBLISHED").
		Preload("Owner").
		Preload("Collaborators.User")

	if scope == "following" {
		query = query.Where("publications.owner_id IN (SELECT following_id FROM user_followers WHERE follower_id = ?)", userID)
	}
	if category != "" {
		query = query.Where("publications.category = ?", category)
	}
	if search != "" {
		term := "%" + strings.ToLower(search) + "%"
		query = query.Where("(LOWER(publications.title) LIKE ? OR LOWER(publications.summary) LIKE ? OR LOWER(publications.content) LIKE ?)", term, term, term)
	}
	if cursorValue != "" {
		query = query.Where("(publications.created_at < ?) OR (publications.created_at = ? AND publications.id < ?)", cursor.CreatedAt, cursor.CreatedAt, cursor.ID)
	}

	if limit < 1 {
		limit = 12
	}
	if limit > 50 {
		limit = 50
	}
	var items []model.Publication
	err = query.Order("publications.created_at DESC").Order("publications.id DESC").Limit(limit + 1).Find(&items).Error
	if err != nil {
		return nil, false, "", err
	}
	hasMore := len(items) > limit
	if hasMore {
		items = items[:limit]
	}
	nextCursor := ""
	if hasMore && len(items) > 0 {
		last := items[len(items)-1]
		nextCursor = EncodeCursor(last.CreatedAt, last.ID)
	}
	return items, hasMore, nextCursor, nil
}

func (r *publicationRepository) GetPublishedByID(id uuid.UUID) (*model.Publication, error) {
	var publication model.Publication
	err := r.db.Where("id = ? AND status = ?", id, "PUBLISHED").
		Preload("Owner").Preload("Collaborators.User").First(&publication).Error
	if err != nil {
		return nil, err
	}
	return &publication, nil
}

func (r *publicationRepository) IsSaved(publicationID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&model.PublicationSave{}).Where("publication_id = ? AND user_id = ?", publicationID, userID).Count(&count).Error
	return count > 0, err
}

func (r *publicationRepository) Save(publicationID, userID uuid.UUID) (int, error) {
	var count int64
	err := r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.PublicationSave{PublicationID: publicationID, UserID: userID})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 1 {
			if err := tx.Model(&model.Publication{}).Where("id = ? AND status = ?", publicationID, "PUBLISHED").UpdateColumn("save_count", gorm.Expr("save_count + 1")).Error; err != nil {
				return err
			}
		}
		return tx.Model(&model.Publication{}).Where("id = ?", publicationID).Select("save_count").Row().Scan(&count)
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *publicationRepository) Unsave(publicationID, userID uuid.UUID) (int, error) {
	var count int64
	err := r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("publication_id = ? AND user_id = ?", publicationID, userID).Delete(&model.PublicationSave{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 1 {
			if err := tx.Model(&model.Publication{}).Where("id = ?", publicationID).UpdateColumn("save_count", gorm.Expr("GREATEST(save_count - 1, 0)")).Error; err != nil {
				return err
			}
		}
		return tx.Model(&model.Publication{}).Where("id = ?", publicationID).Select("save_count").Row().Scan(&count)
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func IsNotFound(err error) bool { return errors.Is(err, gorm.ErrRecordNotFound) }
