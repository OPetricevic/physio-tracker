package doctorprofiles

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	pt "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	out "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/outbound/doctorprofiles"
	re "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/repoerrors"
	se "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/serviceerrors"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

type Service interface {
	Get(ctx context.Context, doctorUUID string) (*pt.DoctorProfile, error)
	Upsert(ctx context.Context, doctorUUID string, req *pt.UpsertDoctorProfileRequest) (*pt.DoctorProfile, error)
}

type service struct {
	repo out.Repository
}

func NewService(repo out.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Get(ctx context.Context, doctorUUID string) (*pt.DoctorProfile, error) {
	if strings.TrimSpace(doctorUUID) == "" {
		return nil, fmt.Errorf("get doctor profile: %w", se.ErrInvalidRequest)
	}
	profile, err := s.repo.GetByDoctor(ctx, doctorUUID)
	if err != nil {
		if errors.Is(err, re.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get doctor profile: %w", se.ErrNotFound)
		}
		return nil, fmt.Errorf("get doctor profile: %w", err)
	}
	return profile, nil
}

func (s *service) Upsert(ctx context.Context, doctorUUID string, req *pt.UpsertDoctorProfileRequest) (*pt.DoctorProfile, error) {
	if strings.TrimSpace(doctorUUID) == "" || req == nil || req.Profile == nil {
		return nil, fmt.Errorf("upsert doctor profile: %w", se.ErrInvalidRequest)
	}
	p := req.GetProfile()
	trim := func(v string) string { return strings.TrimSpace(v) }

	required := []struct {
		name string
		val  string
	}{
		{"practice_name", p.GetPracticeName()},
		{"address", p.GetAddress()},
		{"phone", p.GetPhone()},
	}
	for _, r := range required {
		if trim(r.val) == "" {
			return nil, fmt.Errorf("upsert doctor profile: missing %s: %w", r.name, se.ErrInvalidRequest)
		}
	}

	now := timestamppb.New(time.Now().UTC())
	if trim(p.GetUuid()) == "" {
		p.Uuid = uuid.NewString()
		p.CreatedAt = now
	}
	p.DoctorUuid = doctorUUID
	p.PracticeName = trim(p.GetPracticeName())
	p.Department = trim(p.GetDepartment())
	p.RoleTitle = trim(p.GetRoleTitle())
	p.Address = trim(p.GetAddress())
	p.Phone = trim(p.GetPhone())
	p.Email = trim(p.GetEmail())
	p.Website = trim(p.GetWebsite())
	p.LogoPath = trim(p.GetLogoPath())
	p.ProtocolPrefix = trim(p.GetProtocolPrefix())
	p.HeaderNote = trim(p.GetHeaderNote())
	p.FooterNote = trim(p.GetFooterNote())
	p.UpdatedAt = now

	var existing *pt.DoctorProfile
	if ex, err := s.repo.GetByDoctor(ctx, doctorUUID); err == nil {
		existing = ex
	}

	// Preserve existing logo if caller sends empty; delete old file if replaced.
	if existing != nil {
		if p.GetLogoPath() == "" {
			p.LogoPath = existing.GetLogoPath()
		} else if existing.GetLogoPath() != "" && existing.GetLogoPath() != p.GetLogoPath() {
			_ = deleteLocalStatic(existing.GetLogoPath())
		}
	}

	saved, err := s.repo.Upsert(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("upsert doctor profile: %w", err)
	}
	return saved, nil
}

var _ Service = (*service)(nil)

func deleteLocalStatic(path string) error {
	if path == "" {
		return nil
	}
	const prefix = "/static/"
	if !strings.HasPrefix(path, prefix) {
		return nil
	}
	rel := strings.TrimPrefix(path, prefix)
	local := filepath.Join("uploads", rel)
	if _, err := os.Stat(local); err == nil {
		return os.Remove(local)
	}
	return nil
}
