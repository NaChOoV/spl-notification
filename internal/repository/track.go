package repository

import (
	"database/sql"
	"spl-notification/internal/dto/request"
	"spl-notification/internal/errors"
	"spl-notification/internal/model"
	"time"
)

type trackRepositoryImpl struct {
	db *sql.DB
}

func NewTrackRepositoryImpl(db *sql.DB) TrackRepository {
	return &trackRepositoryImpl{db: db}
}

func (r *trackRepositoryImpl) GetAll() ([]*model.Track, *errors.AppError) {
	query := `
		SELECT 
			id, 
			chat_id, 
			external_id, 
			run, 
			full_name, 
			alias, 
			last_entry, 
			last_exit
		FROM track
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, r.error(err)
	}
	defer rows.Close()

	var tracks []*model.Track
	for rows.Next() {
		track := &model.Track{}

		var lastEntryStr sql.NullString
		var lastExitStr sql.NullString
		var alias sql.NullString
		err := rows.Scan(
			&track.ID,
			&track.ChatID,
			&track.ExternalID,
			&track.Run,
			&track.FullName,
			&alias,
			&lastEntryStr,
			&lastExitStr,
		)
		if err != nil {
			return nil, r.error(err)
		}

		if alias.Valid {
			track.Alias = &alias.String
		}

		if lastEntryStr.Valid {
			lastEntry, err := time.Parse(time.RFC3339, lastEntryStr.String)
			if err != nil {
				return nil, r.error(err)
			}
			track.LastEntry = &lastEntry
		}
		if lastExitStr.Valid {
			lastExit, err := time.Parse(time.RFC3339, lastExitStr.String)
			if err != nil {
				return nil, r.error(err)
			}
			track.LastExit = &lastExit
		}

		tracks = append(tracks, track)
	}

	if err = rows.Err(); err != nil {
		return nil, r.error(err)
	}

	return tracks, nil
}

func (r *trackRepositoryImpl) GetTracksByChatId(chatId string) ([]*model.Track, *errors.AppError) {
	query := `
		SELECT 
			id, 
			chat_id, 
			external_id, 
			run, 
			full_name, 
			alias, 
			last_entry, 
			last_exit
		FROM track 
		WHERE chat_id = ?
	`

	rows, err := r.db.Query(query, chatId)
	if err != nil {
		return nil, r.error(err)
	}
	defer rows.Close()

	var tracks []*model.Track
	for rows.Next() {
		track := &model.Track{}

		var lastEntryStr sql.NullString
		var lastExitStr sql.NullString
		var alias sql.NullString
		err := rows.Scan(
			&track.ID,
			&track.ChatID,
			&track.ExternalID,
			&track.Run,
			&track.FullName,
			&alias,
			&lastEntryStr,
			&lastExitStr,
		)
		if err != nil {
			return nil, r.error(err)
		}

		if alias.Valid {
			track.Alias = &alias.String
		}

		if lastEntryStr.Valid {
			lastEntry, err := time.Parse(time.RFC3339, lastEntryStr.String)
			if err != nil {
				return nil, r.error(err)
			}
			track.LastEntry = &lastEntry
		}
		if lastExitStr.Valid {
			lastExit, err := time.Parse(time.RFC3339, lastExitStr.String)
			if err != nil {
				return nil, r.error(err)
			}
			track.LastExit = &lastExit
		}

		tracks = append(tracks, track)
	}

	if err = rows.Err(); err != nil {
		return nil, r.error(err)
	}

	return tracks, nil
}

func (r *trackRepositoryImpl) UpdateEntryAt(accessArray []*model.Access) *errors.AppError {
	if len(accessArray) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return r.error(err)
	}
	defer tx.Rollback()

	query := `
        UPDATE track
        SET last_entry = ?, updated_at = CURRENT_TIMESTAMP
        WHERE external_id = ?
    `

	stmt, err := tx.Prepare(query)
	if err != nil {
		return r.error(err)
	}
	defer stmt.Close()

	for _, access := range accessArray {
		_, err := stmt.Exec(access.EntryAt.Format(time.RFC3339), access.ExternalID)
		if err != nil {
			return r.error(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return r.error(err)
	}

	return nil
}

func (r *trackRepositoryImpl) UpdateExitAt(accessArray []*model.Access) *errors.AppError {
	if len(accessArray) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return r.error(err)
	}
	defer tx.Rollback()

	query := `
        UPDATE track
        SET last_exit = ?, updated_at = CURRENT_TIMESTAMP
        WHERE external_id = ?
    `

	stmt, err := tx.Prepare(query)
	if err != nil {
		return r.error(err)
	}
	defer stmt.Close()

	for _, access := range accessArray {
		_, err := stmt.Exec(access.ExitAt.Format(time.RFC3339), access.ExternalID)
		if err != nil {
			return r.error(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return r.error(err)
	}

	return nil
}

func (r *trackRepositoryImpl) Create(trackDTO *request.CreateTrackDTO) *errors.AppError {
	query := `
		INSERT INTO track (
			chat_id, external_id, run, full_name, alias, last_entry, last_exit
		) VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(chat_id, run) DO NOTHING
	`

	var lastEntry, lastExit interface{}
	if trackDTO.LastEntry != nil {
		lastEntry = trackDTO.LastEntry.Format(time.RFC3339)
	}
	if trackDTO.LastExit != nil {
		lastExit = trackDTO.LastExit.Format(time.RFC3339)
	}

	_, err := r.db.Exec(
		query,
		trackDTO.ChatID,
		trackDTO.ExternalID,
		trackDTO.Run,
		trackDTO.FullName,
		trackDTO.Alias,
		lastEntry,
		lastExit,
	)

	if err != nil {
		return r.error(err)
	}

	return nil
}

func (r *trackRepositoryImpl) Delete(trackDTO *request.DeleteTrackDTO) *errors.AppError {
	query := `
		DELETE FROM track
		WHERE chat_id = ? AND run = ?
	`

	_, err := r.db.Exec(query, trackDTO.ChatID, trackDTO.Run)
	if err != nil {
		return r.error(err)
	}

	return nil
}

func (r *trackRepositoryImpl) error(err error) *errors.AppError {
	return errors.NewAppError("TrackRepository", err)
}
