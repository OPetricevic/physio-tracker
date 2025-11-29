package doctors

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type service struct {
	repo   Repository
	hasher PasswordHasher
}

func NewService(repo Repository, hasher PasswordHasher) Service {
	return &service{repo: repo, hasher: hasher}
}

func (s *service) Register(ctx context.Context, input RegisterInput) (Doctor, error) {
	if strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.Password) == "" {
		return Doctor{}, errors.New("email i lozinka su obavezni")
	}
	if len(input.Password) < 8 {
		return Doctor{}, errors.New("lozinka mora imati barem 8 znakova")
	}
	if s.hasher == nil {
		return Doctor{}, errors.New("password hasher nije konfiguriran")
	}
	hashed, err := s.hasher.Hash(input.Password)
	if err != nil {
		return Doctor{}, err
	}
	now := time.Now().UTC().Unix()
	doc := Doctor{
		UUID:         uuid.NewString(),
		Email:        strings.ToLower(strings.TrimSpace(input.Email)),
		FirstName:    strings.TrimSpace(input.FirstName),
		LastName:     strings.TrimSpace(input.LastName),
		PasswordHash: hashed,
		CreatedAt:    now,
	}
	if err := s.repo.Insert(ctx, doc); err != nil {
		return Doctor{}, err
	}
	return doc, nil
}

func (s *service) Authenticate(ctx context.Context, email, password string) (Doctor, error) {
	if strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		return Doctor{}, errors.New("email i lozinka su obavezni")
	}
	doc, err := s.repo.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return Doctor{}, err
	}
	if s.hasher == nil {
		return Doctor{}, errors.New("password hasher nije konfiguriran")
	}
	if err := s.hasher.Compare(password, doc.PasswordHash); err != nil {
		return Doctor{}, errors.New("neispravni podaci za prijavu")
	}
	return doc, nil
}

func (s *service) GetByUUID(ctx context.Context, uuid string) (Doctor, error) {
	if strings.TrimSpace(uuid) == "" {
		return Doctor{}, errors.New("uuid je obavezan")
	}
	return s.repo.GetByUUID(ctx, uuid)
}
