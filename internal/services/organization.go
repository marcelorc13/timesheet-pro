package service

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/marcelorc13/timesheet-pro/internal/domain"
	"github.com/marcelorc13/timesheet-pro/internal/repository"
)
type OrganizationService struct {
	repository repository.OrganizationRepository
}

func NewOrganizationService(userRepository repository.OrganizationRepository) *OrganizationService {
	return &OrganizationService{repository: userRepository}
}

func (s *OrganizationService) CreateWithUser(ctx context.Context, co domain.CreateOrganization) (*uuid.UUID, error) {
	validate := validator.New()
	err := validate.Struct(co)
	if err != nil {
		return nil, err.(validator.ValidationErrors)
	}

	orgID, err := s.repository.CreateWithUser(ctx, co)
	if err != nil {
		return nil, err
	}

	return orgID, nil
}
