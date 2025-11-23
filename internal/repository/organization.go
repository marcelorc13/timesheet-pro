package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcelorc13/timesheet-pro/internal/domain"
)

type OrganizationRepository struct {
	DB *pgxpool.Pool
}

func NewOrganizationRepository(db *pgxpool.Pool) *OrganizationRepository {
	return &OrganizationRepository{db}
}

func (r *OrganizationRepository) CreateWithUser(ctx context.Context, co domain.CreateOrganization) (*uuid.UUID, error) {
	var orgID uuid.UUID
	err := pgx.BeginFunc(ctx, r.DB, func(tx pgx.Tx) error {
		const insertOrgQuery = `
			INSERT INTO organizations (name, created_by)
			VALUES (@name, @userID)
			RETURNING id
			`
		args := pgx.StrictNamedArgs{
			"name":   co.Name,
			"userID": co.UserID,
		}

		err := tx.QueryRow(ctx, insertOrgQuery, args).Scan(&orgID)
		if err != nil {
			return err
		}

		const insertUserOrgQuery = `
			INSERT INTO organization_users (user_id, organization_id, organization_role_id)
			SELECT @userID, @orgID, id 
			FROM organization_roles
			WHERE name = @role
			`
		args = pgx.StrictNamedArgs{
			"orgID":  orgID,
			"userID": co.UserID,
			"role":   domain.Admin.String(),
		}

		_, err = tx.Exec(ctx, insertUserOrgQuery, args)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &orgID, nil
}

func (r *OrganizationRepository) GetOrganizationsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Organization, error) {
	const query = `
		SELECT
			o.id,
			o.name,
			o.created_by,
			o.created_at
		FROM
			organizations o
		JOIN
			organization_users ou ON o.id = ou.organization_id
		WHERE
			ou.user_id = @userID
	`
	args := pgx.StrictNamedArgs{
		"userID": userID,
	}

	rows, err := r.DB.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	orgs, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (org domain.Organization, err error) {
		err = rows.Scan(&org.ID, &org.Name, &org.CreatedBy, &org.CreatedAt)

		return
	})
	if err != nil {
		return nil, err
	}

	if len(orgs) == 0 {
		orgs = make([]domain.Organization, 0)
	}

	return orgs, nil
}

func (r *OrganizationRepository) IsUserInOrganization(ctx context.Context, userID, organizationID uuid.UUID) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM organization_users
			WHERE user_id = @userID AND organization_id = @organizationID
		)
	`
	args := pgx.StrictNamedArgs{
		"userID":         userID,
		"organizationID": organizationID,
	}

	var exists bool
	err := r.DB.QueryRow(ctx, query, args).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
