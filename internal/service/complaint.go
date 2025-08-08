package service

import (
	"chow/internal/config"
	"chow/internal/model"
	"chow/internal/repository"
	"context"
	"errors"

	"github.com/google/uuid"
)

// errors
var (
	ErrComplaintAlreadyExist = errors.New("complaint already exists")
	ErrComplaintNotFound     = errors.New("complaint not found")
)

type ComplaintService struct {
	cfg           *config.Config
	complaintRepo *repository.ComplaintRepository
}

func NewComplaintService(cfg *config.Config, complaintRepo *repository.ComplaintRepository) *ComplaintService {
	return &ComplaintService{
		cfg:           cfg,
		complaintRepo: complaintRepo,
	}
}

func (s *ComplaintService) CreateComplaint(ctx context.Context, data *model.Complaint) (*model.Complaint, error) {
	// set complaint status to open
	data.Status = model.OpenComplaint
	complaint, err := s.complaintRepo.Create(ctx, data)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExist) {
			return nil, ErrComplaintAlreadyExist
		}
		return nil, err
	}

	// TODO: notify admins that new complaint has been submitted
	return complaint, err
}

func (s *ComplaintService) GetAllComplaints(ctx context.Context, offset, limit int) ([]*model.Complaint, error) {
	return s.complaintRepo.GetAll(ctx, offset, limit)
}

// GetUserJointComplaints retrieves all complaints made against a joint by a specific user
func (s *ComplaintService) GetUserJointComplaints(ctx context.Context, userID, jointID uuid.UUID, offset, limit int) ([]*model.Complaint, error) {
	return s.complaintRepo.GetUserJointComplaints(ctx, userID, jointID, offset, limit)
}

// GetJointComplaints retrieves all complaints made against a joint
func (s *ComplaintService) GetJointComplaints(ctx context.Context, jointID uuid.UUID, offset, limit int) ([]*model.Complaint, error) {
	return s.complaintRepo.GetJointComplaints(ctx, jointID, offset, limit)
}

// GetUserComplaints retrieves all complaints made by a user
func (s *ComplaintService) GetUserComplaints(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*model.Complaint, error) {
	return s.complaintRepo.GetUserComplaints(ctx, userID, offset, limit)
}

func (s *ComplaintService) GetComplaintByID(ctx context.Context, id uuid.UUID) (*model.Complaint, error) {
	complaint, err := s.complaintRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrComplaintNotFound
		}
		return nil, err
	}

	return complaint, err
}

func (s *ComplaintService) UpdateComplaintStatusByID(ctx context.Context, id uuid.UUID, status model.ComplaintStatus) (*model.Complaint, error) {
	complaint, err := s.complaintRepo.UpdateComplaintStatus(ctx, id, status)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrComplaintNotFound
		}
		return nil, err
	}

	return complaint, err
}
