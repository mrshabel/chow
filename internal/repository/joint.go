package repository

import (
	"context"
	"database/sql"
	"log"

	"chow/internal/model"

	"github.com/google/uuid"
)

// JointRepository handles database operations for joints
type JointRepository struct {
	db *sql.DB
}

// NewJointRepository creates a new joint repository
func NewJointRepository(db *sql.DB) *JointRepository {
	return &JointRepository{db: db}
}

// GetTx returns a transaction that can be passed down to other repository functions. The transaction should be rolled back on error or committed on success by the caller
func (r *JointRepository) GetTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
}

func (r *JointRepository) Create(ctx context.Context, data *model.Joint) (*model.Joint, error) {
	var joint model.Joint
	log.Println(*data)
	query := `
        INSERT INTO joints(name, latitude, longitude, location, description, is_approved, creator_id, photo_url)
        VALUES ($1, $2, $3, ST_Point($4, $5), $6, $7, $8, $9)
		RETURNING id, name, latitude, longitude, description, is_approved, creator_id, photo_url, upvotes, downvotes, created_at, updated_at
    `
	if err := r.db.QueryRowContext(ctx, query, data.Name, data.Latitude, data.Longitude, data.Longitude, data.Latitude, data.Description, data.IsApproved, data.CreatorID, data.PhotoURL).Scan(
		&joint.ID,
		&joint.Name,
		&joint.Latitude,
		&joint.Longitude,
		&joint.Description,
		&joint.IsApproved,
		&joint.CreatorID,
		&joint.PhotoURL,
		&joint.UpVotes,
		&joint.DownVotes,
		&joint.CreatedAt,
		&joint.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &joint, nil
}

func (r *JointRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Joint, error) {
	query := `
		SELECT id, name, latitude, longitude, description, is_approved, creator_id, photo_url, upvotes, downvotes, created_at, updated_at
		FROM joints  
		WHERE id = $1
		`
	var joint model.Joint
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&joint.ID,
		&joint.Name,
		&joint.Latitude,
		&joint.Longitude,
		&joint.Description,
		&joint.IsApproved,
		&joint.CreatorID,
		&joint.PhotoURL,
		&joint.UpVotes,
		&joint.DownVotes,
		&joint.CreatedAt,
		&joint.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &joint, nil
}

func (r *JointRepository) GetAll(ctx context.Context, offset, limit int) ([]*model.Joint, error) {
	query := `
		SELECT id, name, latitude, longitude, description, is_approved, creator_id, photo_url, upvotes, downvotes, created_at, updated_at
		FROM joints
		WHERE is_approved = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
		`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	joints := make([]*model.Joint, 0)
	for rows.Next() {
		var joint model.Joint
		if err := rows.Scan(
			&joint.ID,
			&joint.Name,
			&joint.Latitude,
			&joint.Longitude,
			&joint.Description,
			&joint.IsApproved,
			&joint.CreatorID,
			&joint.PhotoURL,
			&joint.UpVotes,
			&joint.DownVotes,
			&joint.CreatedAt,
			&joint.UpdatedAt,
		); err != nil {
			return nil, err
		}
		log.Println(joint)
		joints = append(joints, &joint)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return joints, nil
}

func (r *JointRepository) UpdateByID(ctx context.Context, id uuid.UUID, data *model.Joint) (*model.Joint, error) {
	query := `
        UPDATE joints 
        SET name = $1, latitude = $2, longitude = $3, location = ST_Point($4, $5), description = $6, is_approved = $7, photo_url = $8, updated_at = NOW()
        WHERE id = $9
        RETURNING id, name, latitude, longitude, description, is_approved, creator_id, photo_url, upvotes, downvotes, created_at, updated_at
    `

	var joint model.Joint
	err := r.db.QueryRowContext(ctx, query, data.Name, data.Latitude, data.Longitude, data.Longitude, data.Latitude, data.Description, data.IsApproved, data.PhotoURL, id).Scan(
		&joint.ID,
		&joint.Name,
		&joint.Latitude,
		&joint.Longitude,
		&joint.Description,
		&joint.IsApproved,
		&joint.CreatorID,
		&joint.PhotoURL,
		&joint.UpVotes,
		&joint.DownVotes,
		&joint.CreatedAt,
		&joint.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &joint, nil
}

func (r *JointRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM joints WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// UpdateVote updates the upvotes/downvotes count for a given joint. A parent transaction should be passed as an argument to ensure that inserting a vote record and upvoting/downvoting is atomic
func (r *JointRepository) UpdateVotes(ctx context.Context, tx *sql.Tx, id uuid.UUID, direction model.VoteDirection, voteExists bool) (*model.Joint, error) {
	query := "UPDATE JOINT SET "
	// increment one direction while decrementing the other for existing votes
	if voteExists {
		if direction == model.UpVote {
			query += "upvotes = upvotes + 1, downvotes = downvotes - 1"
		} else {
			query += "downvotes = downvotes + 1, upvotes = upvotes - 1"
		}
	} else {
		if direction == model.UpVote {
			query += "upvotes = upvotes + 1"
		} else {
			query += "downvotes = downvotes + 1"
		}
	}

	query += `
		WHERE id = $1
		RETURNING id, name, latitude, longitude, description, is_approved, creator_id, photo_url, upvotes, downvotes, created_at, updated_at
	`

	var joint model.Joint
	err := tx.QueryRowContext(ctx, query, id).Scan(
		&joint.ID,
		&joint.Name,
		&joint.Latitude,
		&joint.Longitude,
		&joint.Description,
		&joint.IsApproved,
		&joint.CreatorID,
		&joint.PhotoURL,
		&joint.UpVotes,
		&joint.DownVotes,
		&joint.CreatedAt,
		&joint.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &joint, nil
}

// Search searches for a given joint by name or description
func (r *JointRepository) Search(ctx context.Context, q string, offset, limit int) ([]*model.Joint, error) {
	query := `
		SELECT id, name, latitude, longitude, description, is_approved, creator_id, photo_url, upvotes, downvotes, created_at, updated_at
		FROM joints
		WHERE is_approved = true AND name ILIKE $1 OR description ILIKE $1
		ORDER BY name
		LIMIT $2 OFFSET $3
		`
	rows, err := r.db.QueryContext(ctx, query, "%"+q+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	joints := make([]*model.Joint, 0)
	for rows.Next() {
		var joint model.Joint
		if err := rows.Scan(
			&joint.ID,
			&joint.Name,
			&joint.Latitude,
			&joint.Longitude,
			&joint.Description,
			&joint.IsApproved,
			&joint.CreatorID,
			&joint.PhotoURL,
			&joint.UpVotes,
			&joint.DownVotes,
			&joint.CreatedAt,
			&joint.UpdatedAt,
		); err != nil {
			return nil, err
		}

		joints = append(joints, &joint)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return joints, nil
}

// GetNearby finds nearby joints within the from the given coordinates within the provided radius in meters
func (r *JointRepository) GetNearby(ctx context.Context, cord model.Coordinate, radius float64, offset, limit int) ([]*model.Joint, error) {
	query := `
		SELECT id, name, latitude, longitude, description, is_approved, creator_id, photo_url, upvotes, downvotes, created_at, updated_at, ST_Distance(location, ST_Point($1, $2)::GEOGRAPHY) AS distance
		FROM joints
		WHERE is_approved = true AND 
		-- nearby distance relative to the location. (lon, lat)
		ST_DWithin(location, ST_Point($1, $2)::GEOGRAPHY, $3)
		ORDER BY distance
		LIMIT $4 OFFSET $5
		`
	rows, err := r.db.QueryContext(ctx, query, cord.Longitude, cord.Latitude, radius, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	joints := make([]*model.Joint, 0)
	for rows.Next() {
		var joint model.Joint
		var distance float64

		if err := rows.Scan(
			&joint.ID,
			&joint.Name,
			&joint.Latitude,
			&joint.Longitude,
			&joint.Description,
			&joint.IsApproved,
			&joint.CreatorID,
			&joint.PhotoURL,
			&joint.UpVotes,
			&joint.DownVotes,
			&joint.CreatedAt,
			&joint.UpdatedAt,
			&distance,
		); err != nil {
			return nil, err
		}

		// update distance
		joint.Distance = &distance
		joints = append(joints, &joint)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return joints, nil
}
