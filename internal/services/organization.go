// Package service contains business logic layer for the application.
// It implements validation, authorization, and coordinates between repositories and handlers.
package service

import (
	"context"
	"fmt"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/marcelorc13/timesheet-pro/internal/domain"
	"github.com/marcelorc13/timesheet-pro/internal/repository"
)

type OrganizationService struct {
	repository     repository.OrganizationRepository
	userRepository repository.UserRepository
}

func NewOrganizationService(organizationRepository repository.OrganizationRepository, userRepository repository.UserRepository) *OrganizationService {
	return &OrganizationService{
		repository:     organizationRepository,
		userRepository: userRepository,
	}
}

func (s *OrganizationService) CreateWithUser(ctx context.Context, co domain.CreateOrganization) (*uuid.UUID, error) {
	validate := validator.New()
	err := validate.Struct(co)
	if err != nil {
		return nil, err.(validator.ValidationErrors)
	}

	res, err := s.repository.CreateWithUser(ctx, co)
	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, fmt.Errorf("%s", res.Message)
	}

	orgID, ok := res.Data.(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("erro ao converter dados")
	}

	return &orgID, nil
}

func (s *OrganizationService) GetOrganizationByUserID(ctx context.Context, userID uuid.UUID) (*domain.Organization, error) {
	res, err := s.repository.GetOrganizationByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, fmt.Errorf("%s", res.Message)
	}

	org, ok := res.Data.(domain.Organization)
	if !ok {
		return nil, fmt.Errorf("erro ao converter dados")
	}

	return &org, nil
}

func (s *OrganizationService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	res, err := s.repository.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, nil
	}

	org, ok := res.Data.(domain.Organization)
	if !ok {
		return nil, fmt.Errorf("erro ao converter dados")
	}

	return &org, nil
}

func (s *OrganizationService) Update(ctx context.Context, userID, orgID uuid.UUID, uo domain.UpdateOrganization) error {
	validate := validator.New()
	err := validate.Struct(uo)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	// Check if user is admin
	adminRes, err := s.repository.IsUserAdmin(ctx, userID, orgID)
	if err != nil {
		return err
	}

	if !adminRes.Success {
		return fmt.Errorf("%s", adminRes.Message)
	}

	isAdmin, ok := adminRes.Data.(bool)
	if !ok {
		return fmt.Errorf("erro ao verificar permissões")
	}

	if !isAdmin {
		return fmt.Errorf("usuário não tem permissão para atualizar esta organização")
	}

	res, err := s.repository.Update(ctx, orgID, uo)
	if err != nil {
		return err
	}

	if !res.Success {
		return fmt.Errorf("%s", res.Message)
	}

	return nil
}

func (s *OrganizationService) Delete(ctx context.Context, userID, orgID uuid.UUID) error {
	// Check if user is admin
	adminRes, err := s.repository.IsUserAdmin(ctx, userID, orgID)
	if err != nil {
		return err
	}

	if !adminRes.Success {
		return fmt.Errorf("%s", adminRes.Message)
	}

	isAdmin, ok := adminRes.Data.(bool)
	if !ok {
		return fmt.Errorf("erro ao verificar permissões")
	}

	if !isAdmin {
		return fmt.Errorf("usuário não tem permissão para deletar esta organização")
	}

	res, err := s.repository.Delete(ctx, orgID)
	if err != nil {
		return err
	}

	if !res.Success {
		return fmt.Errorf("%s", res.Message)
	}

	return nil
}

func (s *OrganizationService) IsUserInOrganization(ctx context.Context, userID, organizationID uuid.UUID) (bool, error) {
	res, err := s.repository.IsUserInOrganization(ctx, userID, organizationID)

	if err != nil {
		return false, err
	}

	if !res.Success {
		return false, fmt.Errorf("%s", res.Message)
	}

	exists, ok := res.Data.(bool)
	if !ok {
		return false, fmt.Errorf("erro ao converter dados")
	}

	return exists, nil
}

func (s *OrganizationService) IsUserAdmin(ctx context.Context, userID, organizationID uuid.UUID) (bool, error) {
	res, err := s.repository.IsUserAdmin(ctx, userID, organizationID)

	if err != nil {
		return false, err
	}

	if !res.Success {
		return false, fmt.Errorf("%s", res.Message)
	}

	isAdmin, ok := res.Data.(bool)
	if !ok {
		return false, fmt.Errorf("erro ao converter dados")
	}

	return isAdmin, nil
}

func (s *OrganizationService) AddUserByEmail(ctx context.Context, requestingUserID, orgID uuid.UUID, addUser domain.AddUserToOrganization) error {
	validate := validator.New()
	err := validate.Struct(addUser)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	// Check if requesting user is admin
	adminRes, err := s.repository.IsUserAdmin(ctx, requestingUserID, orgID)
	if err != nil {
		return err
	}

	if !adminRes.Success {
		return fmt.Errorf("%s", adminRes.Message)
	}

	isAdmin, ok := adminRes.Data.(bool)
	if !ok {
		return fmt.Errorf("erro ao verificar permissões")
	}

	if !isAdmin {
		return fmt.Errorf("usuário não tem permissão para adicionar usuários a esta organização")
	}

	// Get user by email
	userRes, err := s.userRepository.GetByEmail(ctx, addUser.Email)
	if err != nil {
		return err
	}

	if !userRes.Success {
		return fmt.Errorf("usuário com email %s não encontrado", addUser.Email)
	}

	user, ok := userRes.Data.(domain.User)
	if !ok {
		return fmt.Errorf("erro ao processar dados do usuário")
	}

	// Check if user is already in the organization
	inOrgRes, err := s.repository.IsUserInOrganization(ctx, user.ID, orgID)
	if err != nil {
		return err
	}

	if !inOrgRes.Success {
		return fmt.Errorf("%s", inOrgRes.Message)
	}

	alreadyMember, ok := inOrgRes.Data.(bool)
	if !ok {
		return fmt.Errorf("erro ao verificar membros")
	}

	if alreadyMember {
		return fmt.Errorf("usuário já é membro desta organização")
	}

	// Add user to organization
	res, err := s.repository.AddUserToOrganization(ctx, user.ID, orgID, addUser.Role)
	if err != nil {
		return err
	}

	if !res.Success {
		return fmt.Errorf("%s", res.Message)
	}

	return nil
}

// GetMembers retrieves all members of an organization
func (s *OrganizationService) GetMembers(ctx context.Context, organizationID uuid.UUID) (*[]domain.OrganizationUser, error) {
	res, err := s.repository.GetOrganizationMembers(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, fmt.Errorf("%s", res.Message)
	}

	members, ok := res.Data.([]domain.OrganizationUser)
	if !ok {
		return nil, fmt.Errorf("erro ao converter dados dos membros")
	}

	return &members, nil
}

func (s *OrganizationService) RemoveUserFromOrganization(ctx context.Context, requestorID, organizationID, targetUserID uuid.UUID) error {
	// Check if requestor is admin
	adminRes, err := s.repository.IsUserAdmin(ctx, requestorID, organizationID)
	if err != nil {
		return err
	}

	if !adminRes.Success {
		return fmt.Errorf("%s", adminRes.Message)
	}

	isAdmin, ok := adminRes.Data.(bool)
	if !ok {
		return fmt.Errorf("erro ao verificar permissões")
	}

	if !isAdmin {
		return fmt.Errorf("usuário não tem permissão para remover membros desta organização")
	}

	// Prevent removing yourself (optional but good practice)
	if requestorID == targetUserID {
		return fmt.Errorf("você não pode remover a si mesmo da organização")
	}

	// Remove user
	res, err := s.repository.RemoveUserFromOrganization(ctx, organizationID, targetUserID)
	if err != nil {
		return err
	}

	if !res.Success {
		return fmt.Errorf("%s", res.Message)
	}

	return nil
}

func (s *OrganizationService) LeaveOrganization(ctx context.Context, userID, organizationID uuid.UUID) error {
	// Check if user is in the organization
	inOrgRes, err := s.repository.IsUserInOrganization(ctx, userID, organizationID)
	if err != nil {
		return err
	}

	if !inOrgRes.Success {
		return fmt.Errorf("%s", inOrgRes.Message)
	}

	isMember, ok := inOrgRes.Data.(bool)
	if !ok {
		return fmt.Errorf("erro ao verificar membros")
	}

	if !isMember {
		return fmt.Errorf("você não é membro desta organização")
	}

	// Remove user
	res, err := s.repository.RemoveUserFromOrganization(ctx, organizationID, userID)
	if err != nil {
		return err
	}

	if !res.Success {
		return fmt.Errorf("%s", res.Message)
	}

	return nil
}
