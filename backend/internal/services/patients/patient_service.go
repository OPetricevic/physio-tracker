package patients

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	pt "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	out "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/outbound/patients"
	re "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/repoerrors"
	se "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/serviceerrors"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, req *pt.CreatePatientRequest) (*pt.Patient, error)
	Update(ctx context.Context, req *pt.UpdatePatientRequest) (*pt.Patient, error)
	List(ctx context.Context, req *pt.ListPatientsRequest, doctorUUID string, pageSize, currentPage int) ([]*pt.Patient, error)
	Delete(ctx context.Context, uuid string) error
}

type service struct {
	repo out.Repository
}

func NewService(repo out.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req *pt.CreatePatientRequest) (*pt.Patient, error) {
	if strings.TrimSpace(req.GetDoctorUuid()) == "" {
		return nil, fmt.Errorf("create patient: %w", se.ErrInvalidRequest)
	}
	if strings.TrimSpace(req.GetFirstName()) == "" || strings.TrimSpace(req.GetLastName()) == "" {
		return nil, fmt.Errorf("create patient: %w", se.ErrInvalidRequest)
	}
	now := time.Now().UTC()
	p := &pt.Patient{
		Uuid:        uuid.NewString(),
		DoctorUuid:  strings.TrimSpace(req.GetDoctorUuid()),
		FirstName:   strings.TrimSpace(req.GetFirstName()),
		LastName:    strings.TrimSpace(req.GetLastName()),
		Phone:       normalizeWrapper(req.Phone),
		Address:     normalizeWrapper(req.Address),
		DateOfBirth: normalizeWrapper(req.DateOfBirth),
		Sex:         normalizeWrapper(req.Sex),
		CreatedAt:   timestamppb.New(now),
		UpdatedAt:   nil,
	}
	created, err := s.repo.Create(ctx, p)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("create patient: %w", se.ErrConflict)
		}
		return nil, fmt.Errorf("create patient: %w", err)
	}
	return created, nil
}

func (s *service) Update(ctx context.Context, req *pt.UpdatePatientRequest) (*pt.Patient, error) {
	if strings.TrimSpace(req.GetUuid()) == "" {
		return nil, fmt.Errorf("update patient: %w", se.ErrInvalidRequest)
	}
	if strings.TrimSpace(req.GetDoctorUuid()) == "" {
		return nil, fmt.Errorf("update patient: %w", se.ErrInvalidRequest)
	}
	if strings.TrimSpace(req.GetFirstName()) == "" || strings.TrimSpace(req.GetLastName()) == "" {
		return nil, fmt.Errorf("update patient: %w", se.ErrInvalidRequest)
	}
	existing, err := s.repo.Get(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, re.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("update patient: %w", se.ErrNotFound)
		}
		return nil, fmt.Errorf("load patient for update: %w", err)
	}
	existing.FirstName = strings.TrimSpace(req.GetFirstName())
	existing.LastName = strings.TrimSpace(req.GetLastName())
	existing.Phone = normalizeWrapper(req.Phone)
	existing.Address = normalizeWrapper(req.Address)
	existing.DateOfBirth = normalizeWrapper(req.DateOfBirth)
	existing.Sex = normalizeWrapper(req.Sex)
	now := time.Now().UTC()
	existing.UpdatedAt = timestamppb.New(now)
	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		if errors.Is(err, re.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("update patient: %w", se.ErrNotFound)
		}
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("update patient: %w", se.ErrConflict)
		}
		return nil, fmt.Errorf("update patient: %w", err)
	}
	return updated, nil
}

func (s *service) List(ctx context.Context, req *pt.ListPatientsRequest, doctorUUID string, pageSize, currentPage int) ([]*pt.Patient, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if currentPage <= 0 {
		currentPage = 1
	}
	offset := (currentPage - 1) * pageSize
	list, err := s.repo.List(ctx, req, doctorUUID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("list patients: %w", err)
	}
	return list, nil
}

func (s *service) Delete(ctx context.Context, uuid string) error {
	if strings.TrimSpace(uuid) == "" {
		return fmt.Errorf("delete patient: %w", se.ErrInvalidRequest)
	}
	if err := s.repo.Delete(ctx, uuid); err != nil {
		if errors.Is(err, re.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("delete patient: %w", se.ErrNotFound)
		}
		return fmt.Errorf("delete patient: %w", err)
	}
	return nil
}

func normalizeWrapper(w *wrapperspb.StringValue) *wrapperspb.StringValue {
	if w == nil {
		return nil
	}
	val := strings.TrimSpace(w.GetValue())
	if val == "" {
		return nil
	}
	return &wrapperspb.StringValue{Value: val}
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
