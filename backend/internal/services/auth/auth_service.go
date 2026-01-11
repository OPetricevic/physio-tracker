package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	"github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/outbound/auth"
	doctorsout "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/outbound/doctors"
	re "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/repoerrors"
	se "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/serviceerrors"
	credrepo "github.com/OPetricevic/physio-tracker/backend/internal/database/doctors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

var (
	ErrInvalidRequest = se.ErrInvalidRequest
	ErrNotFound       = se.ErrNotFound
	ErrConflict       = se.ErrConflict
	ErrUnauthorized   = se.ErrUnauthorized
)

type Service interface {
	Register(ctx context.Context, req *RegisterRequest) (*pb.AuthToken, error)
	Login(ctx context.Context, usernameOrEmail, password string) (*pb.AuthToken, error)
	Logout(ctx context.Context, token string) error
	ChangePassword(ctx context.Context, doctorUUID, currentPassword, newPassword string) error
}

type service struct {
	doctorRepo doctorsout.Repository
	credRepo   credrepo.CredentialsRepository
	tokenRepo  auth.Repository
}

type RegisterRequest struct {
	Email     string
	Username  string
	FirstName string
	LastName  string
	Password  string
}

func NewService(doctorRepo doctorsout.Repository, credRepo credrepo.CredentialsRepository, tokenRepo auth.Repository) Service {
	return &service{
		doctorRepo: doctorRepo,
		credRepo:   credRepo,
		tokenRepo:  tokenRepo,
	}
}

func (s *service) Register(ctx context.Context, req *RegisterRequest) (*pb.AuthToken, error) {
	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		return nil, fmt.Errorf("register: %w", ErrInvalidRequest)
	}
	if strings.TrimSpace(req.FirstName) == "" || strings.TrimSpace(req.LastName) == "" {
		return nil, fmt.Errorf("register: %w", ErrInvalidRequest)
	}
	now := time.Now().UTC()
	doc := &pb.Doctor{
		Uuid:      uuid.NewString(),
		Email:     strings.TrimSpace(req.Email),
		Username:  strings.TrimSpace(req.Username),
		FirstName: strings.TrimSpace(req.FirstName),
		LastName:  strings.TrimSpace(req.LastName),
		CreatedAt: timestamppb.New(now),
		UpdatedAt: nil,
	}
	doc, err := s.doctorRepo.Create(ctx, doc)
	if err != nil {
		if errors.Is(err, re.ErrConflict) {
			return nil, fmt.Errorf("create doctor: %w", ErrConflict)
		}
		return nil, fmt.Errorf("create doctor: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	cred := &pb.DoctorCredentials{
		Uuid:              uuid.NewString(),
		DoctorUuid:        doc.GetUuid(),
		PasswordHash:      string(hash),
		PasswordUpdatedAt: timestamppb.New(now),
	}
	if _, err := s.credRepo.Create(ctx, cred); err != nil {
		if errors.Is(err, re.ErrConflict) {
			return nil, fmt.Errorf("create credentials: %w", se.ErrConflict)
		}
		if errors.Is(err, re.ErrNotFound) {
			// doctor missing / fk violation: treat as conflict for user-facing purposes
			return nil, fmt.Errorf("create credentials: %w", se.ErrConflict)
		}
		return nil, fmt.Errorf("create credentials: %w", err)
	}

	return s.issueToken(ctx, doc.GetUuid())
}

func (s *service) Login(ctx context.Context, usernameOrEmail, password string) (*pb.AuthToken, error) {
	if strings.TrimSpace(usernameOrEmail) == "" || strings.TrimSpace(password) == "" {
		return nil, fmt.Errorf("login: %w", ErrInvalidRequest)
	}
	// find doctor by username or email
	doc, err := s.findDoctor(ctx, usernameOrEmail)
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}
	cred, err := s.credRepo.GetByDoctor(ctx, doc.GetUuid())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("login: %w", ErrUnauthorized)
		}
		return nil, fmt.Errorf("login: get credentials: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(cred.GetPasswordHash()), []byte(password)); err != nil {
		// Legacy plain-text fallback (should be removed once data is clean)
		if cred.GetPasswordHash() != password {
			return nil, fmt.Errorf("login: %w", ErrUnauthorized)
		}
	}
	return s.issueToken(ctx, doc.GetUuid())
}

func (s *service) Logout(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return ErrInvalidRequest
	}
	if err := s.tokenRepo.Delete(ctx, token); err != nil {
		return fmt.Errorf("logout: %w", err)
	}
	return nil
}

func (s *service) ChangePassword(ctx context.Context, doctorUUID, currentPassword, newPassword string) error {
	if strings.TrimSpace(doctorUUID) == "" || strings.TrimSpace(currentPassword) == "" || strings.TrimSpace(newPassword) == "" {
		return fmt.Errorf("change password: %w", ErrInvalidRequest)
	}
	cred, err := s.credRepo.GetByDoctor(ctx, doctorUUID)
	if err != nil {
		if errors.Is(err, re.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("change password: %w", ErrNotFound)
		}
		return fmt.Errorf("change password: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(cred.GetPasswordHash()), []byte(currentPassword)); err != nil {
		return fmt.Errorf("change password: %w", ErrUnauthorized)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("change password: %w", err)
	}
	cred.PasswordHash = string(hash)
	cred.PasswordUpdatedAt = timestamppb.New(time.Now().UTC())
	if _, err := s.credRepo.Update(ctx, cred); err != nil {
		return fmt.Errorf("change password: %w", err)
	}
	return nil
}

func (s *service) issueToken(ctx context.Context, doctorUUID string) (*pb.AuthToken, error) {
	now := time.Now().UTC()
	tok := &pb.AuthToken{
		Uuid:       uuid.NewString(),
		DoctorUuid: doctorUUID,
		Token:      uuid.NewString(),
		ExpiresAt:  timestamppb.New(now.Add(24 * time.Hour)), // adjust lifetime as needed
		CreatedAt:  timestamppb.New(now),
	}
	created, err := s.tokenRepo.Create(ctx, tok)
	if err != nil {
		return nil, fmt.Errorf("issue token: %w", err)
	}
	return created, nil
}

func (s *service) findDoctor(ctx context.Context, identifier string) (*pb.Doctor, error) {
	doc, err := s.doctorRepo.GetByIdentifier(ctx, identifier)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, re.ErrNotFound) || errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("find doctor: %w", ErrUnauthorized)
		}
		return nil, fmt.Errorf("find doctor: %w", err)
	}
	return doc, nil
}
