// Package domain defines the core business entities and data structures.
// It contains domain models, DTOs, and business rules independent of infrastructure.
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
	Address   *Address  `json:"address,omitempty"`
}

type Address struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	ZipCode        string    `json:"zip_code"`
	Complement     string    `json:"complement"`
	PublicPlace    string    `json:"public_place"`
	City           string    `json:"city"`
	State          string    `json:"state"`
}

type CreateOrganization struct {
	UserID      string `json:"user_id" form:"user_id" validate:"required"`
	Name        string `json:"name" form:"name" validate:"required,min=3,max=100"`
	ZipCode     string `json:"zip_code" form:"zip_code" validate:"required"`
	Complement  string `json:"complement" form:"complement" validate:"required"`
	PublicPlace string `json:"public_place" form:"public_place" validate:"required"`
	City        string `json:"city" form:"city" validate:"required"`
	State       string `json:"state" form:"state" validate:"required"`
}

type UpdateOrganization struct {
	Name        string `json:"name" form:"name" validate:"required,min=3,max=100"`
	ZipCode     string `json:"zip_code" form:"zip_code" validate:"required"`
	Complement  string `json:"complement" form:"complement" validate:"required"`
	PublicPlace string `json:"public_place" form:"public_place" validate:"required"`
	City        string `json:"city" form:"city" validate:"required"`
	State       string `json:"state" form:"state" validate:"required"`
}

type AddUserToOrganization struct {
	Email string `json:"email" form:"email" validate:"required,email"`
	Role  string `json:"role" form:"role" validate:"required,oneof=member admin"`
}

type OrganizationUser struct {
	UserID   uuid.UUID `json:"user_id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
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
