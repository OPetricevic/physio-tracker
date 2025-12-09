package auth

import (
	"context"
	"fmt"
	"time"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type TokensRepository struct {
	db *gorm.DB
}

func NewTokensRepository(db *gorm.DB) *TokensRepository {
	return &TokensRepository{db: db}
}

func (r *TokensRepository) Create(ctx context.Context, t *pb.AuthToken) (*pb.AuthToken, error) {
	orm, err := t.ToORM(ctx)
	if err != nil {
		return nil, fmt.Errorf("to ORM: %w", err)
	}
	if err := r.db.WithContext(ctx).Create(&orm).Error; err != nil {
		return nil, fmt.Errorf("insert auth token: %w", err)
	}
	pbTok, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("to PB: %w", err)
	}
	return &pbTok, nil
}

func (r *TokensRepository) Get(ctx context.Context, token string) (*pb.AuthToken, error) {
	var orm pb.AuthTokenORM
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&orm).Error; err != nil {
		return nil, fmt.Errorf("get auth token: %w", err)
	}
	pbTok, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("to PB: %w", err)
	}
	return &pbTok, nil
}

func (r *TokensRepository) Delete(ctx context.Context, token string) error {
	if err := r.db.WithContext(ctx).Where("token = ?", token).Delete(&pb.AuthTokenORM{}).Error; err != nil {
		return fmt.Errorf("delete auth token: %w", err)
	}
	return nil
}

func (r *TokensRepository) DeleteByDoctor(ctx context.Context, doctorUUID string) error {
	if err := r.db.WithContext(ctx).Where("doctor_uuid = ?", doctorUUID).Delete(&pb.AuthTokenORM{}).Error; err != nil {
		return fmt.Errorf("delete auth tokens by doctor: %w", err)
	}
	return nil
}
