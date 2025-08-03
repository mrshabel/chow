package repository

import (
	"context"
	"database/sql"

	"chow/internal/model"

	"github.com/google/uuid"
)

// VoteRepository handles database operations for votes
type VoteRepository struct {
	db *sql.DB
}

// NewVoteRepository creates a new vote repository
func NewVoteRepository(db *sql.DB) *VoteRepository {
	return &VoteRepository{db: db}
}

// GetTx returns a transaction that can be passed down to other repository functions. The transaction should be rolled back on error or committed on success by the caller
func (r *VoteRepository) GetTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
}

func (r *VoteRepository) Create(ctx context.Context, tx *sql.Tx, data *model.Vote) (*model.Vote, error) {
	var vote model.Vote
	query := `
        INSERT INTO votes(user_id, joint_id, direction)
        VALUES ($1, $2, $3)
		RETURNING id, user_id, joint_id, direction, created_at, updated_at
    `
	if err := r.db.QueryRowContext(ctx, query, data.UserID, data.JointID, data.Direction).Scan(
		&vote.ID,
		&vote.UserID,
		&vote.JointID,
		&vote.Direction,
		&vote.CreatedAt,
		&vote.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &vote, nil
}

// Upsert adds a new vote record or update the direction if it already exists
func (r *VoteRepository) Upsert(ctx context.Context, tx *sql.Tx, data *model.Vote) (*model.Vote, error) {
	var vote model.Vote
	query := `
        INSERT INTO votes(user_id, joint_id, direction)
        VALUES ($1, $2, $3)
		ON CONFLICT(user_id, joint_id) 
		DO UPDATE SET
			direction = EXCLUDED.direction, updated_at = NOW()
		RETURNING id, user_id, joint_id, direction, created_at, updated_at
    `
	if err := tx.QueryRowContext(ctx, query, data.UserID, data.JointID, data.Direction).Scan(
		&vote.ID,
		&vote.UserID,
		&vote.JointID,
		&vote.Direction,
		&vote.CreatedAt,
		&vote.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &vote, nil
}

func (r *VoteRepository) GetUserJointVote(ctx context.Context, userID, jointID uuid.UUID) (*model.Vote, error) {
	query := `
		SELECT id, user_id, joint_id, direction, created_at, updated_at
		FROM votes  
		WHERE user_id = $1 AND joint_id = $2
		`
	var vote model.Vote
	err := r.db.QueryRowContext(ctx, query, userID, jointID).Scan(
		&vote.ID,
		&vote.UserID,
		&vote.JointID,
		&vote.Direction,
		&vote.CreatedAt,
		&vote.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &vote, nil
}

func (r *VoteRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Vote, error) {
	query := `
		SELECT id, user_id, joint_id, direction, created_at, updated_at
		FROM votes  
		WHERE id = $1
		`
	var vote model.Vote
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&vote.ID,
		&vote.UserID,
		&vote.JointID,
		&vote.Direction,
		&vote.CreatedAt,
		&vote.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &vote, nil
}

func (r *VoteRepository) GetVotersByJointID(ctx context.Context, jointID uuid.UUID, offset, limit int) ([]*model.JointVoter, error) {
	query := `
		SELECT v.user_id, v.joint_id, v.direction, v.created_at, v.updated_at, u.username 
		FROM votes v
		JOIN users u
		ON v.user_id = u.id
		WHERE v.joint_id = $1
		ORDER BY v.created_at DESC
		LIMIT $2 OFFSET $3
		`
	rows, err := r.db.QueryContext(ctx, query, jointID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var voters []*model.JointVoter
	for rows.Next() {
		var voter model.JointVoter
		if err := rows.Scan(
			&voter.UserID,
			&voter.Username,
			&voter.JointID,
			&voter.Direction,
			&voter.CreatedAt,
			&voter.UpdatedAt,
		); err != nil {
			return nil, err
		}

		voters = append(voters, &voter)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return voters, nil
}

func (r *VoteRepository) UpdateVoteDirectionByID(ctx context.Context, id uuid.UUID, direction model.VoteDirection) (*model.Vote, error) {
	query := `
        UPDATE votes 
        SET direction = $1, updated_at = NOW()
        WHERE id = $2
        RETURNING id, user_id, joint_id, direction, created_at, updated_at
    `
	var vote model.Vote
	err := r.db.QueryRowContext(ctx, query, id, direction).Scan(
		&vote.ID,
		&vote.UserID,
		&vote.JointID,
		&vote.Direction,
		&vote.CreatedAt,
		&vote.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &vote, nil
}
