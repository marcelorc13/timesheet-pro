// Package service contains business logic layer for the application.
// It implements validation, authorization, and coordinates between repositories and handlers.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/marcelorc13/timesheet-pro/internal/domain"
	"github.com/marcelorc13/timesheet-pro/internal/repository"
)

type TimesheetService struct {
	timesheetRepo *repository.TimesheetRepository
	orgRepo       *repository.OrganizationRepository
}

func NewTimesheetService(timesheetRepo *repository.TimesheetRepository, orgRepo *repository.OrganizationRepository) *TimesheetService {
	return &TimesheetService{
		timesheetRepo: timesheetRepo,
		orgRepo:       orgRepo,
	}
}

// ClockIn handles clock in/out for a user in an organization
func (s *TimesheetService) ClockIn(ctx context.Context, userID, orgID uuid.UUID) error {
	// Verify user is member of the organization
	memberRes, err := s.orgRepo.IsUserInOrganization(ctx, userID, orgID)
	if err != nil {
		return err
	}

	if !memberRes.Success {
		return fmt.Errorf("%s", memberRes.Message)
	}

	isMember, ok := memberRes.Data.(bool)
	if !ok || !isMember {
		return fmt.Errorf("usuário não é membro desta organização")
	}

	// Call repository to clock in/out
	res, err := s.timesheetRepo.ClockIn(ctx, orgID, userID)
	if err != nil {
		return err
	}

	if !res.Success {
		return fmt.Errorf("%s", res.Message)
	}

	return nil
}

// GetUserTimesheet retrieves a user's timesheet for a specific date
func (s *TimesheetService) GetUserTimesheet(ctx context.Context, userID, orgID uuid.UUID, date time.Time) (*domain.UserTimesheet, error) {
	// Verify user is member of the organization
	memberRes, err := s.orgRepo.IsUserInOrganization(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}

	if !memberRes.Success {
		return nil, fmt.Errorf("%s", memberRes.Message)
	}

	isMember, ok := memberRes.Data.(bool)
	if !ok || !isMember {
		return nil, fmt.Errorf("usuário não é membro desta organização")
	}

	// Get timesheet from repository
	res, err := s.timesheetRepo.GetUserTimesheet(ctx, userID, orgID, date)
	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, fmt.Errorf("%s", res.Message)
	}

	timesheet, ok := res.Data.(domain.UserTimesheet)
	if !ok {
		return nil, fmt.Errorf("erro ao converter dados do timesheet")
	}

	return &timesheet, nil
}

// GetCurrentStatus returns the current clock in/out status for a user
func (s *TimesheetService) GetCurrentStatus(ctx context.Context, userID, orgID uuid.UUID) (string, *time.Time, error) {
	// Verify user is member of the organization
	memberRes, err := s.orgRepo.IsUserInOrganization(ctx, userID, orgID)
	if err != nil {
		return "", nil, err
	}

	if !memberRes.Success {
		return "", nil, fmt.Errorf("%s", memberRes.Message)
	}

	isMember, ok := memberRes.Data.(bool)
	if !ok || !isMember {
		return "", nil, fmt.Errorf("usuário não é membro desta organização")
	}

	// Get today's timesheet
	today := time.Now().Truncate(24 * time.Hour)
	res, err := s.timesheetRepo.GetUserTimesheet(ctx, userID, orgID, today)
	if err != nil {
		return "", nil, err
	}

	// If no timesheet exists for today, user is clocked out
	if !res.Success {
		return "out", nil, nil
	}

	timesheet, ok := res.Data.(domain.UserTimesheet)
	if !ok {
		return "", nil, fmt.Errorf("erro ao converter dados do timesheet")
	}

	// If no entries, user is clocked out
	if len(timesheet.Entries) == 0 {
		return "out", nil, nil
	}

	// Get last entry
	lastEntry := timesheet.Entries[len(timesheet.Entries)-1]
	
	if lastEntry.TypeID == domain.EntryTypeIn {
		return "in", &lastEntry.Timestamp, nil
	}

	return "out", &lastEntry.Timestamp, nil
}

// GetOrganizationTimesheets retrieves all timesheets for an organization on a specific date
// Requesting user must be admin
func (s *TimesheetService) GetOrganizationTimesheets(ctx context.Context, adminUserID, orgID uuid.UUID, date time.Time) ([]domain.UserTimesheet, error) {
	// Verify user is member of the organization
	memberRes, err := s.orgRepo.IsUserInOrganization(ctx, adminUserID, orgID)
	if err != nil {
		return nil, err
	}

	if !memberRes.Success {
		return nil, fmt.Errorf("%s", memberRes.Message)
	}

	isMember, ok := memberRes.Data.(bool)
	if !ok || !isMember {
		return nil, fmt.Errorf("usuário não é membro desta organização")
	}

	// Check if user is admin
	adminRes, err := s.orgRepo.IsUserAdmin(ctx, adminUserID, orgID)
	if err != nil {
		return nil, err
	}

	if !adminRes.Success {
		return nil, fmt.Errorf("%s", adminRes.Message)
	}

	isAdmin, ok := adminRes.Data.(bool)
	if !ok || !isAdmin {
		return nil, fmt.Errorf("apenas administradores podem visualizar timesheets da organização")
	}

	// TODO: Implement GetOrganizationTimesheets repository method
	// For now, return empty slice
	return []domain.UserTimesheet{}, fmt.Errorf("método não implementado ainda")
}
