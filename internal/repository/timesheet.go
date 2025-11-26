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
			INSERT INTO timesheet_entries (timesheet_id, timestamp, type_id)
			VALUES (@sheetID, @timestamp, @type)
		`
		args := pgx.StrictNamedArgs{
			"sheetID":   timesheetID,
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
