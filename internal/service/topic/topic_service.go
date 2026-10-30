package topicservice

import (
	"errors"
	"strings"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	topicrepository "nalakarsa/internal/repository/topic"

	"github.com/google/uuid"
)

type TopicService interface {
	List(search string) ([]dto.TopicResponse, error)
	Suggest(userID uuid.UUID, req dto.SuggestTopicRequest) (dto.TopicResponse, error)
	ListUserTopics(userID uuid.UUID) ([]dto.TopicResponse, error)
	UpdateUserTopics(userID uuid.UUID, req dto.UpdateInterestsRequest) ([]dto.TopicResponse, error)
}

type topicService struct {
	repo topicrepository.TopicRepository
}

func NewTopicService(repo topicrepository.TopicRepository) TopicService {
	return &topicService{repo: repo}
}

func (s *topicService) List(search string) ([]dto.TopicResponse, error) {
	topics, err := s.repo.List(search)
	return mapTopics(topics), err
}

func (s *topicService) Suggest(userID uuid.UUID, req dto.SuggestTopicRequest) (dto.TopicResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return dto.TopicResponse{}, errors.New("topic name is required")
	}
	topic := &model.Topic{Name: name, Slug: slugify(name), Status: "pending", CreatedBy: &userID}
	if err := s.repo.Create(topic); err != nil {
		return dto.TopicResponse{}, err
	}
	return mapTopic(*topic), nil
}

func (s *topicService) ListUserTopics(userID uuid.UUID) ([]dto.TopicResponse, error) {
	topics, err := s.repo.ListUserTopics(userID)
	return mapTopics(topics), err
}

func (s *topicService) UpdateUserTopics(userID uuid.UUID, req dto.UpdateInterestsRequest) ([]dto.TopicResponse, error) {
	if len(req.TopicIDs) < 3 {
		return nil, errors.New("at least 3 topics are required")
	}
	seen := map[uuid.UUID]bool{}
	for _, id := range req.TopicIDs {
		if id == uuid.Nil || seen[id] {
			return nil, errors.New("topic_ids must contain unique valid IDs")
		}
		seen[id] = true
	}
	if err := s.repo.ReplaceUserTopics(userID, req.TopicIDs); err != nil {
		return nil, err
	}
	return s.ListUserTopics(userID)
}

func mapTopic(topic model.Topic) dto.TopicResponse {
	return dto.TopicResponse{ID: topic.ID, Name: topic.Name, Slug: topic.Slug, Status: topic.Status, CreatedAt: topic.CreatedAt}
}

func mapTopics(topics []model.Topic) []dto.TopicResponse {
	result := make([]dto.TopicResponse, 0, len(topics))
	for _, topic := range topics {
		result = append(result, mapTopic(topic))
	}
	return result
}

func slugify(value string) string {
	return strings.Join(strings.Fields(strings.ToLower(value)), "-")
}
