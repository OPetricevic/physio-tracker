package auth

import (
	"context"
	"fmt"
	"time"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	"gorm.io/gorm"
)

type TokensRepository struct {
	db *gorm.DB
}

func NewTokensRepository(db *gorm.DB) *TokensRepository {
	return &TokensRepository{db: db}
}

func (r *TokensRepository) Create(ctx context.Context, t *pb.AuthToken) (*pb.AuthToken, error) {
	if err := r.db.WithContext(ctx).Table("auth_tokens").Create(t).Error; err != nil {
		return nil, fmt.Errorf("insert auth token: %w", err)
	}
	return t, nil
}

func (r *TokensRepository) Get(ctx context.Context, token string) (*pb.AuthToken, error) {
	var t pb.AuthToken
	if err := r.db.WithContext(ctx).Table("auth_tokens").Where("token = ?", token).First(&t).Error; err != nil {
		return nil, fmt.Errorf("get auth token: %w", err)
	}
	return &t, nil
}

func (r *TokensRepository) Delete(ctx context.Context, token string) error {
	if err := r.db.WithContext(ctx).Table("auth_tokens").Where("token = ?", token).Delete(&tokenModel{}).Error; err != nil {
		return fmt.Errorf("delete auth token: %w", err)
	}
	return nil
}

func (r *TokensRepository) DeleteByDoctor(ctx context.Context, doctorUUID string) error {
	if err := r.db.WithContext(ctx).Table("auth_tokens").Where("doctor_uuid = ?", doctorUUID).Delete(&tokenModel{}).Error; err != nil {
		return fmt.Errorf("delete auth tokens by doctor: %w", err)
	}
	return nil
}

type tokenModel struct {
	DoctorUUID string    `gorm:"column:doctor_uuid"`
	Token      string    `gorm:"column:token"`
	ExpiresAt  time.Time `gorm:"column:expires_at"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}
