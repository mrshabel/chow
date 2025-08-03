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
	ErrJointAlreadyExist       = errors.New("joint already exists")
	ErrJointNotFound           = errors.New("joint not found")
	ErrMaxSearchRadiusExceeded = errors.New("maximum search radius exceeded")
)

type JointService struct {
	cfg       *config.Config
	jointRepo *repository.JointRepository
	voteRepo  *repository.VoteRepository
}

func NewJointService(cfg *config.Config, jointRepo *repository.JointRepository, voteRepo *repository.VoteRepository) *JointService {
	return &JointService{
		cfg:       cfg,
		jointRepo: jointRepo,
		voteRepo:  voteRepo,
	}
}

func (s *JointService) CreateJoint(ctx context.Context, data *model.Joint) (*model.Joint, error) {
	joint, err := s.jointRepo.Create(ctx, data)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExist) {
			return nil, ErrJointAlreadyExist
		}
		return nil, err
	}

	// TODO: notify admins that joint has been created
	return joint, err
}

// VoteForJoint updates the vote metrics for a given joint.
// If an attempt is made to upvote an upvoted joint by the same user or vice versa, the request returns immediately with no error or joint information
func (s *JointService) VoteForJoint(ctx context.Context, id uuid.UUID, data *model.Vote) (*model.Joint, error) {
	// get previous vote record if present
	vote, err := s.voteRepo.GetUserJointVote(ctx, data.UserID, data.JointID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	// skip processing if directions are the same
	voteExists := vote != nil
	if voteExists && data.Direction == vote.Direction {
		return nil, nil
	}

	// start outer transaction
	tx, err := s.jointRepo.GetTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// upsert vote record
	_, err = s.voteRepo.Upsert(ctx, tx, data)
	if err != nil {
		return nil, err
	}

	// apply vote. if it's an existing vote, the direction vote count is incremented while its opposite direction vote count is decremented
	joint, err := s.jointRepo.UpdateVotes(ctx, tx, id, data.Direction, voteExists)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrJointNotFound
		}
		return nil, err
	}

	// finally commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return joint, err
}

// GetNearbyJoints returns the closest joints from the provided coordinates with the specified radius
func (s *JointService) GetNearbyJoints(ctx context.Context, coord model.Coordinate, radius float64, offset, limit int) ([]*model.Joint, error) {
	// validate max radius
	if radius > s.cfg.MaxNearbyRadius {
		return nil, ErrMaxSearchRadiusExceeded
	}

	// perform search
	return s.jointRepo.GetNearby(ctx, coord, radius, offset, limit)
}

func (s *JointService) GetJointByID(ctx context.Context, id uuid.UUID) (*model.Joint, error) {
	joint, err := s.jointRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrJointNotFound
		}
		return nil, err
	}

	return joint, err
}

func (s *JointService) SearchForJointByNameOrDescription(ctx context.Context, query string, offset, limit int) ([]*model.Joint, error) {
	return s.jointRepo.Search(ctx, query, offset, limit)
}

func (s *JointService) GetAllJoints(ctx context.Context, offset, limit int) ([]*model.Joint, error) {
	return s.jointRepo.GetAll(ctx, offset, limit)
}

func (s *JointService) UpdateJointByID(ctx context.Context, id uuid.UUID, data *model.Joint) (*model.Joint, error) {
	joint, err := s.jointRepo.UpdateByID(ctx, id, data)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrJointNotFound
		}
		return nil, err
	}

	return joint, err
}

func (s *JointService) DeleteJointByID(ctx context.Context, id uuid.UUID) error {
	if err := s.jointRepo.DeleteByID(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrJointNotFound
		}
		return err
	}

	return nil
}
