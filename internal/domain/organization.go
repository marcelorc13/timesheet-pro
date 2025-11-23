package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedBy uuid.UUID `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateOrganization struct {
	UserID string `json:"user_id" form:"user_id" validate:"required"`
	Name   string    `json:"name" form:"name" validate:"required,min=3,max=100"`
}

type OrganizationUser struct {
	UserID   uuid.UUID `json:"user_id"`
	Name     string    `json:"name"`
	Role     Role      `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type Role string

const (
	Member Role = "member"
	Admin  Role = "admin"
)

func (r Role) String() string {
	return string(r)
}

func ParseRole(s string) (Role, error) {
	r := Role(s)
	switch r {
	case Member:
		return r, nil
	case Admin:
		return r, nil
	default:
		return "", fmt.Errorf("invalid role: %q", s)
	}
}
