package topicrepository

import (
	"errors"
	"strings"

	"nalakarsa/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TopicRepository interface {
	List(search string) ([]model.Topic, error)
	Create(topic *model.Topic) error
	GetApprovedByIDs(ids []uuid.UUID) ([]model.Topic, error)
	ReplaceUserTopics(userID uuid.UUID, topicIDs []uuid.UUID) error
	ListUserTopics(userID uuid.UUID) ([]model.Topic, error)
}

type pgTopicRepository struct{ db *gorm.DB }

func NewTopicRepository(db *gorm.DB) TopicRepository { return &pgTopicRepository{db: db} }

func (r *pgTopicRepository) List(search string) ([]model.Topic, error) {
	var topics []model.Topic
	q := r.db.Where("status = ?", "approved").Order("name asc")
	if strings.TrimSpace(search) != "" {
		q = q.Where("name ILIKE ?", "%"+strings.TrimSpace(search)+"%")
	}
	return topics, q.Find(&topics).Error
}

func (r *pgTopicRepository) Create(topic *model.Topic) error {
	return r.db.Create(topic).Error
}

func (r *pgTopicRepository) GetApprovedByIDs(ids []uuid.UUID) ([]model.Topic, error) {
	var topics []model.Topic
	if len(ids) == 0 {
		return topics, nil
	}
	err := r.db.Where("id IN ? AND status = ?", ids, "approved").Find(&topics).Error
	return topics, err
}

func (r *pgTopicRepository) ReplaceUserTopics(userID uuid.UUID, topicIDs []uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var topics []model.Topic
		if err := tx.Where("id IN ? AND status = ?", topicIDs, "approved").Find(&topics).Error; err != nil {
			return err
		}
		if len(topics) != len(topicIDs) {
			return errors.New("one or more topics are invalid or not approved")
		}
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserTopic{}).Error; err != nil {
			return err
		}
		rows := make([]model.UserTopic, 0, len(topicIDs))
		for _, id := range topicIDs {
			rows = append(rows, model.UserTopic{UserID: userID, TopicID: id})
		}
		return tx.Create(&rows).Error
	})
}

func (r *pgTopicRepository) ListUserTopics(userID uuid.UUID) ([]model.Topic, error) {
	var topics []model.Topic
	err := r.db.Model(&model.Topic{}).
		Joins("JOIN user_topics ON user_topics.topic_id = topics.id").
		Where("user_topics.user_id = ?", userID).
		Order("topics.name asc").Find(&topics).Error
	return topics, err
}
