package doctors

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	out "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/outbound/doctors"
	re "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/repoerrors"
	se "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/serviceerrors"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, req *pb.CreateDoctorRequest) (*pb.Doctor, error)
	Update(ctx context.Context, req *pb.UpdateDoctorRequest) (*pb.Doctor, error)
	Delete(ctx context.Context, uuid string) error
	List(ctx context.Context, query string, pageSize, currentPage int) ([]*pb.Doctor, error)
}

type service struct {
	repo out.Repository
}

func NewService(repo out.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req *pb.CreateDoctorRequest) (*pb.Doctor, error) {
	if strings.TrimSpace(req.GetEmail()) == "" || strings.TrimSpace(req.GetUsername()) == "" {
		return nil, fmt.Errorf("create doctor: %w", se.ErrInvalidRequest)
	}
	if strings.TrimSpace(req.GetFirstName()) == "" || strings.TrimSpace(req.GetLastName()) == "" {
		return nil, fmt.Errorf("create doctor: %w", se.ErrInvalidRequest)
	}
	now := time.Now().UTC()
	doc := &pb.Doctor{
		Uuid:      uuid.NewString(),
		Email:     strings.TrimSpace(req.GetEmail()),
		Username:  strings.TrimSpace(req.GetUsername()),
		FirstName: strings.TrimSpace(req.GetFirstName()),
		LastName:  strings.TrimSpace(req.GetLastName()),
		CreatedAt: timestamppb.New(now),
		UpdatedAt: nil,
	}
	created, err := s.repo.Create(ctx, doc)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("create doctor: %w", se.ErrConflict)
		}
		return nil, fmt.Errorf("create doctor: %w", err)
	}
	return created, nil
}

func (s *service) Update(ctx context.Context, req *pb.UpdateDoctorRequest) (*pb.Doctor, error) {
	if strings.TrimSpace(req.GetUuid()) == "" {
		return nil, fmt.Errorf("update doctor: %w", se.ErrInvalidRequest)
	}
	if strings.TrimSpace(req.GetEmail()) == "" || strings.TrimSpace(req.GetUsername()) == "" {
		return nil, fmt.Errorf("update doctor: %w", se.ErrInvalidRequest)
	}
	if strings.TrimSpace(req.GetFirstName()) == "" || strings.TrimSpace(req.GetLastName()) == "" {
		return nil, fmt.Errorf("update doctor: %w", se.ErrInvalidRequest)
	}
	existing, err := s.repo.Get(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, se.ErrNotFound) || errors.Is(err, re.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("update doctor: %w", se.ErrNotFound)
		}
		return nil, fmt.Errorf("update doctor: load doctor: %w", err)
	}
	existing.Email = strings.TrimSpace(req.GetEmail())
	existing.Username = strings.TrimSpace(req.GetUsername())
	existing.FirstName = strings.TrimSpace(req.GetFirstName())
	existing.LastName = strings.TrimSpace(req.GetLastName())
	now := time.Now().UTC()
	existing.UpdatedAt = timestamppb.New(now)

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		if errors.Is(err, se.ErrNotFound) || errors.Is(err, re.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("update doctor: %w", se.ErrNotFound)
		}
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("update doctor: %w", se.ErrConflict)
		}
		return nil, fmt.Errorf("update doctor: %w", err)
	}
	return updated, nil
}

func (s *service) Delete(ctx context.Context, uuid string) error {
	if strings.TrimSpace(uuid) == "" {
		return fmt.Errorf("delete doctor: %w", se.ErrInvalidRequest)
	}
	if err := s.repo.Delete(ctx, uuid); err != nil {
		if errors.Is(err, se.ErrNotFound) || errors.Is(err, re.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("delete doctor: %w", se.ErrNotFound)
		}
		return fmt.Errorf("delete doctor: %w", err)
	}
	return nil
}

func (s *service) List(ctx context.Context, query string, pageSize, currentPage int) ([]*pb.Doctor, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if currentPage <= 0 {
		currentPage = 1
	}
	offset := (currentPage - 1) * pageSize
	list, err := s.repo.List(ctx, query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("list doctors: %w", err)
	}
	return list, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
