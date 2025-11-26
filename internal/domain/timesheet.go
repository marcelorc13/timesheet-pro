package domain

import (
	"time"

	"github.com/google/uuid"
)

// TimesheetStatus represents the status of a daily timesheet
type TimesheetStatus int

const (
	StatusOpen TimesheetStatus = iota + 1
	StatusClosed
	StatusAbsent
	StatusApproved
)

// EntryType represents the type of timesheet entry (clock in/out)
type EntryType int

const (
	EntryTypeIn EntryType = iota + 1
	EntryTypeOut
)

// DailyTimesheet represents a user's timesheet for a specific day
type DailyTimesheet struct {
	ID             uuid.UUID       `json:"id"`
	UserID         uuid.UUID       `json:"user_id"`
	OrganizationID uuid.UUID       `json:"organization_id"`
	Date           time.Time       `json:"date"`
	StatusID       TimesheetStatus `json:"status_id"`
	TotalMinutes   int64           `json:"total_minutes"`
	CreatedAt      time.Time       `json:"created_at"`
	Entries        []TimesheetEntry `json:"entries,omitempty"`
}

// TimesheetEntry represents a single clock in/out entry
type TimesheetEntry struct {
	ID             uuid.UUID `json:"id"`
	TimesheetID    uuid.UUID `json:"timesheet_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	TypeID         EntryType `json:"type_id"`
	Timestamp      time.Time `json:"timestamp"`
}

// UserTimesheet combines user information with their timesheet
type UserTimesheet struct {
	DailyTimesheet
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}
