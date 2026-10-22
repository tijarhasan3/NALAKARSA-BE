package publicationservice

import (
	"errors"
	"strings"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	publicationrepository "nalakarsa/internal/repository/publication"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrInvalidScope = errors.New("scope must be either for_you or following")

func IsNotFound(err error) bool { return errors.Is(err, gorm.ErrRecordNotFound) }

type PublicationService interface {
	Feed(userID uuid.UUID, scope, category, search, cursor string, limit int) (dto.PublicationFeedResponse, error)
	GetByID(userID, publicationID uuid.UUID) (dto.PublicationResponse, error)
	Save(userID, publicationID uuid.UUID) (dto.PublicationSaveResponse, error)
	Unsave(userID, publicationID uuid.UUID) (dto.PublicationSaveResponse, error)
}

type publicationService struct {
	repo publicationrepository.PublicationRepository
}

func NewPublicationService(repo publicationrepository.PublicationRepository) PublicationService {
	return &publicationService{repo: repo}
}

func (s *publicationService) Feed(userID uuid.UUID, scope, category, search, cursor string, limit int) (dto.PublicationFeedResponse, error) {
	scope = strings.TrimSpace(scope)
	if scope == "" {
		scope = "for_you"
	}
	if scope != "for_you" && scope != "following" {
		return dto.PublicationFeedResponse{}, ErrInvalidScope
	}
	if limit < 1 {
		limit = 12
	}
	if limit > 50 {
		limit = 50
	}
	items, hasMore, nextCursor, err := s.repo.ListPublished(userID, scope, strings.TrimSpace(category), strings.TrimSpace(search), cursor, limit)
	if err != nil {
		return dto.PublicationFeedResponse{}, err
	}
	result := dto.PublicationFeedResponse{Meta: dto.PublicationFeedMeta{HasMore: hasMore, Limit: limit}}
	if nextCursor != "" {
		result.Meta.NextCursor = &nextCursor
	}
	for _, item := range items {
		mapped, err := s.mapPublication(item, userID)
		if err != nil {
			return dto.PublicationFeedResponse{}, err
		}
		result.Data = append(result.Data, mapped)
	}
	if result.Data == nil {
		result.Data = []dto.PublicationResponse{}
	}
	return result, nil
}

func (s *publicationService) GetByID(userID, publicationID uuid.UUID) (dto.PublicationResponse, error) {
	item, err := s.repo.GetPublishedByID(publicationID)
	if err != nil {
		return dto.PublicationResponse{}, err
	}
	return s.mapPublication(*item, userID)
}

func (s *publicationService) Save(userID, publicationID uuid.UUID) (dto.PublicationSaveResponse, error) {
	count, err := s.repo.Save(publicationID, userID)
	if err != nil {
		return dto.PublicationSaveResponse{}, err
	}
	return dto.PublicationSaveResponse{PublicationID: publicationID, IsSaved: true, SaveCount: count}, nil
}

func (s *publicationService) Unsave(userID, publicationID uuid.UUID) (dto.PublicationSaveResponse, error) {
	count, err := s.repo.Unsave(publicationID, userID)
	if err != nil {
		return dto.PublicationSaveResponse{}, err
	}
	return dto.PublicationSaveResponse{PublicationID: publicationID, IsSaved: false, SaveCount: count}, nil
}

func (s *publicationService) mapPublication(item model.Publication, userID uuid.UUID) (dto.PublicationResponse, error) {
	isSaved, err := s.repo.IsSaved(item.ID, userID)
	if err != nil {
		return dto.PublicationResponse{}, err
	}
	response := dto.PublicationResponse{
		ID: item.ID, PublicationType: item.PublicationType, Title: item.Title,
		Summary: item.Summary, Content: item.Content, Category: item.Category,
		Status: item.Status, CoverURL: item.CoverURL, ViewCount: item.ViewCount,
		SaveCount: item.SaveCount, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt,
		IsSaved: isSaved, Collaborators: []dto.PublicationCollaboratorResponse{},
	}
	for _, collaborator := range item.Collaborators {
		avatar := collaborator.User.AvatarURL
		var avatarURL *string
		if avatar != "" {
			avatarURL = &avatar
		}
		response.Collaborators = append(response.Collaborators, dto.PublicationCollaboratorResponse{
			ID: collaborator.User.ID, Name: collaborator.User.FullName,
			Role: collaborator.ContributionRole, Affiliation: collaborator.User.Affiliation,
			AvatarURL: avatarURL,
		})
	}
	return response, nil
}
