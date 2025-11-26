// Package repository provides database access layer for the application.
// It contains repository implementations for data persistence using PostgreSQL.
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

func (r *OrganizationRepository) CreateWithUser(ctx context.Context, co domain.CreateOrganization) (domain.DBResponse, error) {
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

		const insertAddressQuery = `
			INSERT INTO addresses (organization_id, zip_code, complement, public_place, city, state)
			VALUES (@orgID, @zipCode, @complement, @publicPlace, @city, @state)
		`
		args = pgx.StrictNamedArgs{
			"orgID":       orgID,
			"zipCode":     co.ZipCode,
			"complement":  co.Complement,
			"publicPlace": co.PublicPlace,
			"city":        co.City,
			"state":       co.State,
		}

		_, err = tx.Exec(ctx, insertAddressQuery, args)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return domain.DBResponse{Message: "erro ao criar organização"}, err
	}

	return domain.DBResponse{Success: true, Data: orgID}, nil
}

func (r *OrganizationRepository) GetOrganizationByUserID(ctx context.Context, userID uuid.UUID) (domain.DBResponse, error) {
	const query = `
		SELECT
			o.id,
			o.name,
			o.created_by,
			o.created_at,
			a.id,
			a.organization_id,
			a.zip_code,
			a.complement,
			a.public_place,
			a.city,
			a.state
		FROM
			organizations o
		JOIN
			organization_users ou ON o.id = ou.organization_id
		LEFT JOIN
			addresses a ON o.id = a.organization_id
		WHERE
			ou.user_id = @userID
		LIMIT 1
	`
	args := pgx.StrictNamedArgs{
		"userID": userID,
	}

	var org domain.Organization
	var addr domain.Address
	var addrID, addrOrgID *uuid.UUID
	var zipCode, complement, publicPlace, city, state *string

	err := r.DB.QueryRow(ctx, query, args).Scan(
		&org.ID, &org.Name, &org.CreatedBy, &org.CreatedAt,
		&addrID, &addrOrgID, &zipCode, &complement, &publicPlace, &city, &state,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.DBResponse{Success: false, Message: "usuário não pertence a nenhuma organização"}, nil
		}
		return domain.DBResponse{Success: false, Message: "erro ao buscar organização"}, err
	}

	if addrID != nil {
		addr.ID = *addrID
		addr.OrganizationID = *addrOrgID
		if zipCode != nil {
			addr.ZipCode = *zipCode
		}
		if complement != nil {
			addr.Complement = *complement
		}
		if publicPlace != nil {
			addr.PublicPlace = *publicPlace
		}
		if city != nil {
			addr.City = *city
		}
		if state != nil {
			addr.State = *state
		}
		org.Address = &addr
	}

	return domain.DBResponse{Success: true, Data: org}, nil
}

func (r *OrganizationRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.DBResponse, error) {
	const query = `
		SELECT
			o.id,
			o.name,
			o.created_by,
			o.created_at,
			a.id,
			a.organization_id,
			a.zip_code,
			a.complement,
			a.public_place,
			a.city,
			a.state
		FROM
			organizations o
		LEFT JOIN
			addresses a ON o.id = a.organization_id
		WHERE
			o.id = @id
	`
	args := pgx.StrictNamedArgs{
		"id": id,
	}

	var org domain.Organization
	var addr domain.Address
	var addrID, addrOrgID *uuid.UUID
	var zipCode, complement, publicPlace, city, state *string

	err := r.DB.QueryRow(ctx, query, args).Scan(
		&org.ID, &org.Name, &org.CreatedBy, &org.CreatedAt,
		&addrID, &addrOrgID, &zipCode, &complement, &publicPlace, &city, &state,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.DBResponse{Success: false, Message: "organização não encontrada"}, nil
		}
		return domain.DBResponse{Success: false, Message: "erro ao buscar organização"}, err
	}

	if addrID != nil {
		addr.ID = *addrID
		addr.OrganizationID = *addrOrgID
		if zipCode != nil {
			addr.ZipCode = *zipCode
		}
		if complement != nil {
			addr.Complement = *complement
		}
		if publicPlace != nil {
			addr.PublicPlace = *publicPlace
		}
		if city != nil {
			addr.City = *city
		}
		if state != nil {
			addr.State = *state
		}
		org.Address = &addr
	}

	return domain.DBResponse{Success: true, Data: org}, nil
}

func (r *OrganizationRepository) Update(ctx context.Context, orgID uuid.UUID, uo domain.UpdateOrganization) (domain.DBResponse, error) {
	err := pgx.BeginFunc(ctx, r.DB, func(tx pgx.Tx) error {
		const updateOrgQuery = `
			UPDATE organizations
			SET name = @name
			WHERE id = @id
		`
		args := pgx.StrictNamedArgs{
			"id":   orgID,
			"name": uo.Name,
		}

		_, err := tx.Exec(ctx, updateOrgQuery, args)
		if err != nil {
			return err
		}

		// Upsert address
		const upsertAddressQuery = `
			INSERT INTO addresses (organization_id, zip_code, complement, public_place, city, state)
			VALUES (@orgID, @zipCode, @complement, @publicPlace, @city, @state)
			ON CONFLICT (organization_id) DO UPDATE
			SET
				zip_code = EXCLUDED.zip_code,
				complement = EXCLUDED.complement,
				public_place = EXCLUDED.public_place,
				city = EXCLUDED.city,
				state = EXCLUDED.state
		`
		addrArgs := pgx.StrictNamedArgs{
			"orgID":       orgID,
			"zipCode":     uo.ZipCode,
			"complement":  uo.Complement,
			"publicPlace": uo.PublicPlace,
			"city":        uo.City,
			"state":       uo.State,
		}

		_, err = tx.Exec(ctx, upsertAddressQuery, addrArgs)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return domain.DBResponse{Success: false, Message: "erro ao atualizar organização"}, err
	}

	return domain.DBResponse{Success: true, Message: "organização atualizada com sucesso"}, nil
}

func (r *OrganizationRepository) Delete(ctx context.Context, id uuid.UUID) (domain.DBResponse, error) {
	const query = `
		DELETE FROM organizations
		WHERE id = @id
	`
	args := pgx.StrictNamedArgs{
		"id": id,
	}

	res, err := r.DB.Exec(ctx, query, args)
	if err != nil {
		return domain.DBResponse{Message: "erro ao deletar organização"}, err
	}

	rows := res.RowsAffected()
	if rows != 1 {
		return domain.DBResponse{Message: "organização não encontrada"}, nil
	}

	return domain.DBResponse{Success: true}, nil
}

func (r *OrganizationRepository) IsUserInOrganization(ctx context.Context, userID, organizationID uuid.UUID) (domain.DBResponse, error) {
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
		return domain.DBResponse{Message: "erro ao verificar usuário na organização"}, err
	}

	return domain.DBResponse{Success: true, Data: exists}, nil
}

func (r *OrganizationRepository) IsUserAdmin(ctx context.Context, userID, organizationID uuid.UUID) (domain.DBResponse, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM organization_users
			WHERE user_id = @userID AND organization_id = @organizationID AND organization_role_id = 2
		)
	`
	args := pgx.StrictNamedArgs{
		"userID":         userID,
		"organizationID": organizationID,
	}

	var exists bool
	err := r.DB.QueryRow(ctx, query, args).Scan(&exists)
	if err != nil {
		return domain.DBResponse{Message: "erro ao verificar permissão de admin"}, err
	}

	return domain.DBResponse{Success: true, Data: exists}, nil
}

func (r *OrganizationRepository) AddUserToOrganization(ctx context.Context, userID, organizationID uuid.UUID, role string) (domain.DBResponse, error) {
	const query = `
		INSERT INTO organization_users (user_id, organization_id, organization_role_id)
		SELECT @userID, @organizationID, id 
		FROM organization_roles
		WHERE name = @role
	`
	args := pgx.StrictNamedArgs{
		"userID":         userID,
		"organizationID": organizationID,
		"role":           role,
	}

	res, err := r.DB.Exec(ctx, query, args)
	if err != nil {
		return domain.DBResponse{Message: "erro ao adicionar usuário à organização"}, err
	}

	rows := res.RowsAffected()
	if rows != 1 {
		return domain.DBResponse{Message: "erro ao adicionar usuário à organização"}, nil
	}

	return domain.DBResponse{Success: true}, nil
}

// GetOrganizationMembers retrieves all members of an organization with their details
func (r *OrganizationRepository) GetOrganizationMembers(ctx context.Context, organizationID uuid.UUID) (domain.DBResponse, error) {
	const query = `
		SELECT 
			u.id, 
			u.name, 
			u.email,
			r.name as role,
			ou.joined_at as joined_at
		FROM users u
		JOIN organization_users ou ON u.id = ou.user_id
		JOIN organization_roles r ON ou.organization_role_id = r.id
		WHERE ou.organization_id = @organizationID
		ORDER BY ou.joined_at
	`
	args := pgx.StrictNamedArgs{
		"organizationID": organizationID,
	}

	rows, err := r.DB.Query(ctx, query, args)
	if err != nil {
		return domain.DBResponse{Success: false, Message: "erro ao buscar membros"}, err
	}
	defer rows.Close()

	var members []domain.OrganizationUser
	for rows.Next() {
		var member domain.OrganizationUser
		var roleStr string

		err := rows.Scan(
			&member.UserID,
			&member.Name,
			&member.Email,
			&roleStr,
			&member.JoinedAt,
		)
		if err != nil {
			return domain.DBResponse{Success: false, Message: "erro ao ler dados do membro"}, err
		}

		role, err := domain.ParseRole(roleStr)
		if err != nil {
			continue // Skip invalid roles
		}
		member.Role = role

		members = append(members, member)
	}

	if rows.Err() != nil {
		return domain.DBResponse{Success: false, Message: "erro ao processar membros"}, rows.Err()
	}

	// Return empty slice if no members, not nil
	if members == nil {
		members = []domain.OrganizationUser{}
	}

	return domain.DBResponse{Success: true, Data: members}, nil
}

func (r *OrganizationRepository) RemoveUserFromOrganization(ctx context.Context, organizationID, userID uuid.UUID) (domain.DBResponse, error) {
	const query = `
		DELETE FROM organization_users
		WHERE organization_id = @organizationID AND user_id = @userID
	`
	args := pgx.StrictNamedArgs{
		"organizationID": organizationID,
		"userID":         userID,
	}

	res, err := r.DB.Exec(ctx, query, args)
	if err != nil {
		return domain.DBResponse{Message: "erro ao remover usuário da organização"}, err
	}

	rows := res.RowsAffected()
	if rows != 1 {
		return domain.DBResponse{Message: "usuário não encontrado na organização"}, nil
	}

	return domain.DBResponse{Success: true}, nil
}
