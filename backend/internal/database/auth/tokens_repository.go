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
	rec := protoToTokenRecord(t)
	if err := r.db.WithContext(ctx).Table("auth_tokens").Create(rec).Error; err != nil {
		return nil, fmt.Errorf("insert auth token: %w", err)
	}
	return tokenRecordToProto(*rec), nil
}

func (r *TokensRepository) Get(ctx context.Context, token string) (*pb.AuthToken, error) {
	var rec tokenRecord
	if err := r.db.WithContext(ctx).Table("auth_tokens").Where("token = ?", token).First(&rec).Error; err != nil {
		return nil, fmt.Errorf("get auth token: %w", err)
	}
	return tokenRecordToProto(rec), nil
}

func (r *TokensRepository) Delete(ctx context.Context, token string) error {
	if err := r.db.WithContext(ctx).Table("auth_tokens").Where("token = ?", token).Delete(&tokenRecord{}).Error; err != nil {
		return fmt.Errorf("delete auth token: %w", err)
	}
	return nil
}

func (r *TokensRepository) DeleteByDoctor(ctx context.Context, doctorUUID string) error {
	if err := r.db.WithContext(ctx).Table("auth_tokens").Where("doctor_uuid = ?", doctorUUID).Delete(&tokenRecord{}).Error; err != nil {
		return fmt.Errorf("delete auth tokens by doctor: %w", err)
	}
	return nil
}

type tokenRecord struct {
	UUID       string    `gorm:"column:uuid"`
	DoctorUUID string    `gorm:"column:doctor_uuid"`
	Token      string    `gorm:"column:token"`
	ExpiresAt  time.Time `gorm:"column:expires_at"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func protoToTokenRecord(t *pb.AuthToken) *tokenRecord {
	return &tokenRecord{
		UUID:       t.GetUuid(),
		DoctorUUID: t.GetDoctorUuid(),
		Token:      t.GetToken(),
		ExpiresAt:  tsToTime(t.ExpiresAt),
		CreatedAt:  tsToTime(t.CreatedAt),
	}
}

func tokenRecordToProto(rec tokenRecord) *pb.AuthToken {
	return &pb.AuthToken{
		Uuid:       rec.UUID,
		DoctorUuid: rec.DoctorUUID,
		Token:      rec.Token,
		ExpiresAt:  timestamppb.New(rec.ExpiresAt),
		CreatedAt:  timestamppb.New(rec.CreatedAt),
	}
}

func tsToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}
