package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcelorc13/timesheet-pro/internal/domain"
)

type TimesheetRepository struct {
	DB *pgxpool.Pool
}

func NewTimesheetRepository(db *pgxpool.Pool) *TimesheetRepository {
	return &TimesheetRepository{db}
}

func (r TimesheetRepository) ClockIn(ctx context.Context, orgID, userID uuid.UUID) (domain.DBResponse, error) {
	err := pgx.BeginFunc(ctx, r.DB, func(tx pgx.Tx) error {
		now := time.Now()
		today := now.Truncate(24 * time.Hour) 

		var timesheetID uuid.UUID

		const findQuery = `
			SELECT id FROM daily_timesheets 
			WHERE user_id = @userID AND date = @date
			FOR UPDATE
		`
		err := tx.QueryRow(ctx, findQuery, pgx.NamedArgs{"userID": userID, "date": today}).Scan(&timesheetID)

		if err == pgx.ErrNoRows {
			const insertSheetQuery = `
				INSERT INTO daily_timesheets (user_id, organization_id, date, status_id)
				VALUES (@userID, @orgID, @date, 1)
				RETURNING id
			`
			err = tx.QueryRow(ctx, insertSheetQuery, pgx.NamedArgs{
				"userID": userID,
				"orgID":  orgID,
				"date":   today,
			}).Scan(&timesheetID)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		var count int
		const countQuery = `SELECT count(*) FROM timesheet_entries WHERE timesheet_id = @id`
		err = tx.QueryRow(ctx, countQuery, pgx.NamedArgs{"id": timesheetID}).Scan(&count)
		if err != nil {
			return err
		}

		nextType := 1 // "in"
		if count%2 != 0 {
			nextType = 2 // "out"
		}

		// 5. Inserir a Batida (Entry)
		const insertEntryQuery = `
			INSERT INTO timesheet_entries (timesheet_id, organization_id, timestamp, type_id)
			VALUES (@sheetID, @orgID, @timestamp, @type)
		`
		args := pgx.StrictNamedArgs{
			"sheetID":   timesheetID,
			"orgID":   orgID,
			"timestamp": now,
			"type":      nextType,
		}
		_, err = tx.Exec(ctx, insertEntryQuery, args)
		return err
	})
	if err != nil {
		return domain.DBResponse{Message: err.Error()}, err
	}

	return domain.DBResponse{Success: true}, nil
}

func (r *TimesheetRepository) GetUserTimesheet(ctx context.Context, userID, orgID uuid.UUID, date time.Time) (domain.DBResponse, error) {
	dateOnly := date.Truncate(24 * time.Hour)

	const timesheetQuery = `
		SELECT 
			dt.id,
			dt.user_id,
			dt.organization_id,
			dt.date,
			dt.status_id,
			dt.total_minutes,
			dt.created_at,
			u.name as user_name,
			u.email as user_email
		FROM daily_timesheets dt
		JOIN users u ON dt.user_id = u.id
		WHERE dt.user_id = @userID 
			AND dt.organization_id = @orgID 
			AND dt.date = @date
	`
	args := pgx.NamedArgs{
		"userID": userID,
		"orgID":  orgID,
		"date":   dateOnly,
	}

	var timesheet domain.UserTimesheet
	err := r.DB.QueryRow(ctx, timesheetQuery, args).Scan(
		&timesheet.ID,
		&timesheet.UserID,
		&timesheet.OrganizationID,
		&timesheet.Date,
		&timesheet.StatusID,
		&timesheet.TotalMinutes,
		&timesheet.CreatedAt,
		&timesheet.UserName,
		&timesheet.UserEmail,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.DBResponse{Success: false, Message: "timesheet não encontrado para esta data"}, nil
		}
		return domain.DBResponse{Success: false, Message: "erro ao buscar timesheet"}, err
	}

	// Get all entries for this timesheet
	const entriesQuery = `
		SELECT 
			id,
			timesheet_id,
			organization_id,
			type_id,
			timestamp
		FROM timesheet_entries
		WHERE timesheet_id = @timesheetID
		ORDER BY timestamp ASC
	`
	entriesArgs := pgx.NamedArgs{
		"timesheetID": timesheet.ID,
	}

	rows, err := r.DB.Query(ctx, entriesQuery, entriesArgs)
	if err != nil {
		return domain.DBResponse{Success: false, Message: "erro ao buscar entradas do timesheet"}, err
	}
	defer rows.Close()

	var entries []domain.TimesheetEntry
	for rows.Next() {
		var entry domain.TimesheetEntry
		err := rows.Scan(
			&entry.ID,
			&entry.TimesheetID,
			&entry.OrganizationID,
			&entry.TypeID,
			&entry.Timestamp,
		)
		if err != nil {
			return domain.DBResponse{Success: false, Message: "erro ao ler entrada do timesheet"}, err
		}
		entries = append(entries, entry)
	}

	if rows.Err() != nil {
		return domain.DBResponse{Success: false, Message: "erro ao processar entradas"}, rows.Err()
	}

	timesheet.Entries = entries

	return domain.DBResponse{Success: true, Data: timesheet}, nil
}

// GetTimesheetByID retrieves a single timesheet by its ID with all entries
func (r *TimesheetRepository) GetTimesheetByID(ctx context.Context, timesheetID uuid.UUID) (domain.DBResponse, error) {
	const timesheetQuery = `
		SELECT 
			dt.id,
			dt.user_id,
			dt.organization_id,
			dt.date,
			dt.status_id,
			dt.total_minutes,
			dt.created_at,
			u.name as user_name,
			u.email as user_email
		FROM daily_timesheets dt
		JOIN users u ON dt.user_id = u.id
		WHERE dt.id = @timesheetID
	`
	args := pgx.NamedArgs{
		"timesheetID": timesheetID,
	}

	var timesheet domain.UserTimesheet
	err := r.DB.QueryRow(ctx, timesheetQuery, args).Scan(
		&timesheet.ID,
		&timesheet.UserID,
		&timesheet.OrganizationID,
		&timesheet.Date,
		&timesheet.StatusID,
		&timesheet.TotalMinutes,
		&timesheet.CreatedAt,
		&timesheet.UserName,
		&timesheet.UserEmail,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.DBResponse{Success: false, Message: "timesheet não encontrado"}, nil
		}
		return domain.DBResponse{Success: false, Message: "erro ao buscar timesheet"}, err
	}

	// Get all entries for this timesheet
	const entriesQuery = `
		SELECT 
			id,
			timesheet_id,
			organization_id,
			type_id,
			timestamp
		FROM timesheet_entries
		WHERE timesheet_id = @timesheetID
		ORDER BY timestamp ASC
	`
	entriesArgs := pgx.NamedArgs{
		"timesheetID": timesheet.ID,
	}

	rows, err := r.DB.Query(ctx, entriesQuery, entriesArgs)
	if err != nil {
		return domain.DBResponse{Success: false, Message: "erro ao buscar entradas do timesheet"}, err
	}
	defer rows.Close()

	var entries []domain.TimesheetEntry
	for rows.Next() {
		var entry domain.TimesheetEntry
		err := rows.Scan(
			&entry.ID,
			&entry.TimesheetID,
			&entry.OrganizationID,
			&entry.TypeID,
			&entry.Timestamp,
		)
		if err != nil {
			return domain.DBResponse{Success: false, Message: "erro ao ler entrada do timesheet"}, err
		}
		entries = append(entries, entry)
	}

	if rows.Err() != nil {
		return domain.DBResponse{Success: false, Message: "erro ao processar entradas"}, rows.Err()
	}

	timesheet.Entries = entries

	return domain.DBResponse{Success: true, Data: timesheet}, nil
}

// GetUserTimesheets retrieves all timesheets for a specific user within a date range
func (r *TimesheetRepository) GetUserTimesheets(ctx context.Context, userID, orgID uuid.UUID, startDate, endDate time.Time) (domain.DBResponse, error) {
	start := startDate.Truncate(24 * time.Hour)
	end := endDate.Truncate(24 * time.Hour)

	const query = `
		SELECT 
			dt.id,
			dt.user_id,
			dt.organization_id,
			dt.date,
			dt.status_id,
			dt.total_minutes,
			dt.created_at,
			u.name as user_name,
			u.email as user_email
		FROM daily_timesheets dt
		JOIN users u ON dt.user_id = u.id
		WHERE dt.user_id = @userID 
			AND dt.organization_id = @orgID 
			AND dt.date >= @startDate
			AND dt.date <= @endDate
		ORDER BY dt.date DESC
	`
	args := pgx.NamedArgs{
		"userID":    userID,
		"orgID":     orgID,
		"startDate": start,
		"endDate":   end,
	}

	rows, err := r.DB.Query(ctx, query, args)
	if err != nil {
		return domain.DBResponse{Success: false, Message: "erro ao buscar timesheets"}, err
	}
	defer rows.Close()

	var timesheets []domain.UserTimesheet
	for rows.Next() {
		var ts domain.UserTimesheet
		err := rows.Scan(
			&ts.ID,
			&ts.UserID,
			&ts.OrganizationID,
			&ts.Date,
			&ts.StatusID,
			&ts.TotalMinutes,
			&ts.CreatedAt,
			&ts.UserName,
			&ts.UserEmail,
		)
		if err != nil {
			return domain.DBResponse{Success: false, Message: "erro ao ler timesheet"}, err
		}
		
		// Get entries for this timesheet
		const entriesQuery = `
			SELECT 
				id,
				timesheet_id,
				organization_id,
				type_id,
				timestamp
			FROM timesheet_entries
			WHERE timesheet_id = @timesheetID
			ORDER BY timestamp ASC
		`
		entriesRows, err := r.DB.Query(ctx, entriesQuery, pgx.NamedArgs{"timesheetID": ts.ID})
		if err != nil {
			return domain.DBResponse{Success: false, Message: "erro ao buscar entradas"}, err
		}

		var entries []domain.TimesheetEntry
		for entriesRows.Next() {
			var entry domain.TimesheetEntry
			err := entriesRows.Scan(
				&entry.ID,
				&entry.TimesheetID,
				&entry.OrganizationID,
				&entry.TypeID,
				&entry.Timestamp,
			)
			if err != nil {
				entriesRows.Close()
				return domain.DBResponse{Success: false, Message: "erro ao ler entrada"}, err
			}
			entries = append(entries, entry)
		}
		entriesRows.Close()

		ts.Entries = entries
		timesheets = append(timesheets, ts)
	}

	if rows.Err() != nil {
		return domain.DBResponse{Success: false, Message: "erro ao processar timesheets"}, rows.Err()
	}

	// Return empty slice if no timesheets found
	if timesheets == nil {
		timesheets = []domain.UserTimesheet{}
	}

	return domain.DBResponse{Success: true, Data: timesheets}, nil
}

// GetOrganizationTimesheets retrieves all timesheets for an organization on a specific date
func (r *TimesheetRepository) GetOrganizationTimesheets(ctx context.Context, orgID uuid.UUID, date time.Time) (domain.DBResponse, error) {
	dateOnly := date.Truncate(24 * time.Hour)

	const query = `
		SELECT 
			dt.id,
			dt.user_id,
			dt.organization_id,
			dt.date,
			dt.status_id,
			dt.total_minutes,
			dt.created_at,
			u.name as user_name,
			u.email as user_email
		FROM daily_timesheets dt
		JOIN users u ON dt.user_id = u.id
		WHERE dt.organization_id = @orgID 
			AND dt.date = @date
		ORDER BY u.name ASC
	`
	args := pgx.NamedArgs{
		"orgID": orgID,
		"date":  dateOnly,
	}

	rows, err := r.DB.Query(ctx, query, args)
	if err != nil {
		return domain.DBResponse{Success: false, Message: "erro ao buscar timesheets da organização"}, err
	}
	defer rows.Close()

	var timesheets []domain.UserTimesheet
	for rows.Next() {
		var ts domain.UserTimesheet
		err := rows.Scan(
			&ts.ID,
			&ts.UserID,
			&ts.OrganizationID,
			&ts.Date,
			&ts.StatusID,
			&ts.TotalMinutes,
			&ts.CreatedAt,
			&ts.UserName,
			&ts.UserEmail,
		)
		if err != nil {
			return domain.DBResponse{Success: false, Message: "erro ao ler timesheet"}, err
		}

		// Get entries for this timesheet
		const entriesQuery = `
			SELECT 
				id,
				timesheet_id,
				organization_id,
				type_id,
				timestamp
			FROM timesheet_entries
			WHERE timesheet_id = @timesheetID
			ORDER BY timestamp ASC
		`
		entriesRows, err := r.DB.Query(ctx, entriesQuery, pgx.NamedArgs{"timesheetID": ts.ID})
		if err != nil {
			return domain.DBResponse{Success: false, Message: "erro ao buscar entradas"}, err
		}

		var entries []domain.TimesheetEntry
		for entriesRows.Next() {
			var entry domain.TimesheetEntry
			err := entriesRows.Scan(
				&entry.ID,
				&entry.TimesheetID,
				&entry.OrganizationID,
				&entry.TypeID,
				&entry.Timestamp,
			)
			if err != nil {
				entriesRows.Close()
				return domain.DBResponse{Success: false, Message: "erro ao ler entrada"}, err
			}
			entries = append(entries, entry)
		}
		entriesRows.Close()

		ts.Entries = entries
		timesheets = append(timesheets, ts)
	}

	if rows.Err() != nil {
		return domain.DBResponse{Success: false, Message: "erro ao processar timesheets"}, rows.Err()
	}

	// Return empty slice if no timesheets found
	if timesheets == nil {
		timesheets = []domain.UserTimesheet{}
	}

	return domain.DBResponse{Success: true, Data: timesheets}, nil
}
