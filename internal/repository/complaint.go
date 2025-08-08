package repository

import (
	"context"
	"database/sql"

	"chow/internal/model"

	"github.com/google/uuid"
)

// ComplaintRepository handles database operations for complaints
type ComplaintRepository struct {
	db *sql.DB
}

// NewComplaintRepository creates a new complaint repository
func NewComplaintRepository(db *sql.DB) *ComplaintRepository {
	return &ComplaintRepository{db: db}
}

func (r *ComplaintRepository) Create(ctx context.Context, data *model.Complaint) (*model.Complaint, error) {
	var complaint model.Complaint
	query := `
        INSERT INTO complaints(joint_id, user_id, reason, status)
        VALUES ($1, $2, $3, $4)
		RETURNING id, joint_id, user_id, reason, status, created_at, updated_at
    `
	if err := r.db.QueryRowContext(ctx, query, data.JointID, data.UserID, data.Reason, data.Status).Scan(
		&complaint.ID,
		&complaint.JointID,
		&complaint.UserID,
		&complaint.Reason,
		&complaint.Status,
		&complaint.CreatedAt,
		&complaint.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &complaint, nil
}

func (r *ComplaintRepository) GetAll(ctx context.Context, offset, limit int) ([]*model.Complaint, error) {
	query := `
		SELECT id, joint_id, user_id, reason, status, created_at, updated_at
		FROM complaints  
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
		`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	complaints := make([]*model.Complaint, 0, limit)

	for rows.Next() {
		var complaint model.Complaint
		if err := rows.Scan(
			&complaint.ID,
			&complaint.JointID,
			&complaint.UserID,
			&complaint.Reason,
			&complaint.Status,
			&complaint.CreatedAt,
			&complaint.UpdatedAt,
		); err != nil {
			return nil, err
		}
		complaints = append(complaints, &complaint)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return complaints, nil
}

func (r *ComplaintRepository) GetUserComplaints(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*model.Complaint, error) {
	query := `
		SELECT id, joint_id, user_id, reason, status, created_at, updated_at
		FROM complaints  
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
		`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	complaints := make([]*model.Complaint, 0, limit)

	for rows.Next() {
		var complaint model.Complaint
		if err := rows.Scan(
			&complaint.ID,
			&complaint.JointID,
			&complaint.UserID,
			&complaint.Reason,
			&complaint.Status,
			&complaint.CreatedAt,
			&complaint.UpdatedAt,
		); err != nil {
			return nil, err
		}
		complaints = append(complaints, &complaint)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return complaints, nil
}

func (r *ComplaintRepository) GetJointComplaints(ctx context.Context, jointID uuid.UUID, offset, limit int) ([]*model.Complaint, error) {
	query := `
		SELECT id, joint_id, user_id, reason, status, created_at, updated_at
		FROM complaints  
		WHERE joint_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
		`
	rows, err := r.db.QueryContext(ctx, query, jointID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	complaints := make([]*model.Complaint, 0, limit)
	for rows.Next() {
		var complaint model.Complaint
		if err := rows.Scan(
			&complaint.ID,
			&complaint.JointID,
			&complaint.UserID,
			&complaint.Reason,
			&complaint.Status,
			&complaint.CreatedAt,
			&complaint.UpdatedAt,
		); err != nil {
			return nil, err
		}
		complaints = append(complaints, &complaint)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return complaints, nil
}

// GetUserJointComplaints returns all complaints made by a user on a particular joint
func (r *ComplaintRepository) GetUserJointComplaints(ctx context.Context, userID, jointID uuid.UUID, offset, limit int) ([]*model.Complaint, error) {
	query := `
		SELECT id, joint_id, user_id, reason, status, created_at, updated_at
		FROM complaints  
		WHERE user_id = $1 AND joint_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
		`
	rows, err := r.db.QueryContext(ctx, query, userID, jointID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	complaints := make([]*model.Complaint, 0, limit)

	for rows.Next() {
		var complaint model.Complaint
		if err := rows.Scan(
			&complaint.ID,
			&complaint.JointID,
			&complaint.UserID,
			&complaint.Reason,
			&complaint.Status,
			&complaint.CreatedAt,
			&complaint.UpdatedAt,
		); err != nil {
			return nil, err
		}
		complaints = append(complaints, &complaint)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return complaints, nil
}

func (r *ComplaintRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Complaint, error) {
	query := `
		SELECT id, joint_id, user_id, reason, status, created_at, updated_at
		FROM complaints  
		WHERE id = $1
		`
	var complaint model.Complaint
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&complaint.ID,
		&complaint.JointID,
		&complaint.UserID,
		&complaint.Reason,
		&complaint.Status,
		&complaint.CreatedAt,
		&complaint.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &complaint, nil
}

func (r *ComplaintRepository) UpdateComplaintStatus(ctx context.Context, id uuid.UUID, status model.ComplaintStatus) (*model.Complaint, error) {
	query := `
        UPDATE complaints 
        SET status = $1, updated_at = NOW()
        WHERE id = $2
        RETURNING id, joint_id, user_id, reason, status, created_at, updated_at
    `
	var complaint model.Complaint
	err := r.db.QueryRowContext(ctx, query, status, id).Scan(
		&complaint.ID,
		&complaint.JointID,
		&complaint.UserID,
		&complaint.Reason,
		&complaint.Status,
		&complaint.CreatedAt,
		&complaint.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &complaint, nil
}
