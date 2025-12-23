package anamneses

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	out "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/outbound/anamneses"
	re "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/repoerrors"
	se "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/serviceerrors"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, doctorUUID string, req *pb.CreateAnamnesisRequest) (*pb.Anamnesis, error)
	Update(ctx context.Context, doctorUUID string, req *pb.UpdateAnamnesisRequest) (*pb.Anamnesis, error)
	List(ctx context.Context, doctorUUID, patientUUID, query string, pageSize, currentPage int) ([]*pb.Anamnesis, error)
	Delete(ctx context.Context, doctorUUID, uuid string) error
	Get(ctx context.Context, doctorUUID, uuid string) (*pb.Anamnesis, error)
}

type service struct {
	repo out.Repository
}

func NewService(repo out.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, doctorUUID string, req *pb.CreateAnamnesisRequest) (*pb.Anamnesis, error) {
	if strings.TrimSpace(req.GetPatientUuid()) == "" {
		return nil, fmt.Errorf("create anamnesis: %w", se.ErrInvalidRequest)
	}
	if strings.TrimSpace(req.GetAnamnesis()) == "" || strings.TrimSpace(req.GetDiagnosis()) == "" || strings.TrimSpace(req.GetTherapy()) == "" {
		return nil, fmt.Errorf("create anamnesis: %w", se.ErrInvalidRequest)
	}
	now := time.Now().UTC()
	a := &pb.Anamnesis{
		Uuid:             uuid.NewString(),
		PatientUuid:      strings.TrimSpace(req.GetPatientUuid()),
		Anamnesis:        strings.TrimSpace(req.GetAnamnesis()),
		Diagnosis:        strings.TrimSpace(req.GetDiagnosis()),
		Therapy:          strings.TrimSpace(req.GetTherapy()),
		OtherInfo:        strings.TrimSpace(req.GetOtherInfo()),
		IncludeVisitUuids: req.IncludeVisitUuids,
		CreatedAt:        timestamppb.New(now),
		UpdatedAt:        nil,
	}
	created, err := s.repo.Create(ctx, a)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("create anamnesis: %w", se.ErrConflict)
		}
		return nil, fmt.Errorf("create anamnesis: %w", err)
	}
	return created, nil
}

func (s *service) Update(ctx context.Context, doctorUUID string, req *pb.UpdateAnamnesisRequest) (*pb.Anamnesis, error) {
	if strings.TrimSpace(req.GetUuid()) == "" || strings.TrimSpace(req.GetPatientUuid()) == "" {
		return nil, fmt.Errorf("update anamnesis: %w", se.ErrInvalidRequest)
	}
	existing, err := s.repo.Get(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, re.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("update anamnesis: %w", se.ErrNotFound)
		}
		return nil, fmt.Errorf("load anamnesis for update: %w", err)
	}
	if strings.TrimSpace(existing.GetPatientUuid()) != strings.TrimSpace(req.GetPatientUuid()) {
		return nil, fmt.Errorf("update anamnesis: %w", se.ErrInvalidRequest)
	}
	if strings.TrimSpace(req.GetAnamnesis()) != "" {
		existing.Anamnesis = strings.TrimSpace(req.GetAnamnesis())
	}
	if strings.TrimSpace(req.GetDiagnosis()) != "" {
		existing.Diagnosis = strings.TrimSpace(req.GetDiagnosis())
	}
	if strings.TrimSpace(req.GetTherapy()) != "" {
		existing.Therapy = strings.TrimSpace(req.GetTherapy())
	}
	if strings.TrimSpace(req.GetOtherInfo()) != "" {
		existing.OtherInfo = strings.TrimSpace(req.GetOtherInfo())
	}
	if req.IncludeVisitUuids != nil {
		existing.IncludeVisitUuids = req.IncludeVisitUuids
	}
	existing.UpdatedAt = timestamppb.New(time.Now().UTC())

	updated, err := s.repo.Update(ctx, existing)
		if err != nil {
		if errors.Is(err, re.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("update anamnesis: %w", se.ErrNotFound)
		}
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("update anamnesis: %w", se.ErrConflict)
		}
		return nil, fmt.Errorf("update anamnesis: %w", err)
	}
	return updated, nil
}

func (s *service) List(ctx context.Context, doctorUUID, patientUUID, query string, pageSize, currentPage int) ([]*pb.Anamnesis, error) {
	if strings.TrimSpace(patientUUID) == "" {
		return nil, fmt.Errorf("list anamneses: %w", se.ErrInvalidRequest)
	}
	if pageSize <= 0 {
		pageSize = 5
	}
	if currentPage <= 0 {
		currentPage = 1
	}
	offset := (currentPage - 1) * pageSize
	list, err := s.repo.List(ctx, patientUUID, doctorUUID, query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("list anamneses: %w", err)
	}
	return list, nil
}

func (s *service) Delete(ctx context.Context, doctorUUID, uuid string) error {
	if strings.TrimSpace(uuid) == "" {
		return fmt.Errorf("delete anamnesis: %w", se.ErrInvalidRequest)
	}
	if err := s.repo.Delete(ctx, uuid); err != nil {
		if errors.Is(err, re.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("delete anamnesis: %w", se.ErrNotFound)
		}
		return fmt.Errorf("delete anamnesis: %w", err)
	}
	return nil
}

func (s *service) Get(ctx context.Context, doctorUUID, uuid string) (*pb.Anamnesis, error) {
	if strings.TrimSpace(uuid) == "" {
		return nil, fmt.Errorf("get anamnesis: %w", se.ErrInvalidRequest)
	}
	anm, err := s.repo.Get(ctx, uuid)
	if err != nil {
		if errors.Is(err, re.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get anamnesis: %w", se.ErrNotFound)
		}
		return nil, fmt.Errorf("get anamnesis: %w", err)
	}
	return anm, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
