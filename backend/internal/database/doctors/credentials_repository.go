package doctors

import (
	"context"
	"fmt"
	"time"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// CredentialsRepository stores password hashes.
type CredentialsRepository interface {
	Create(ctx context.Context, c *pb.DoctorCredentials) (*pb.DoctorCredentials, error)
	GetByDoctor(ctx context.Context, doctorUUID string) (*pb.DoctorCredentials, error)
	Update(ctx context.Context, c *pb.DoctorCredentials) (*pb.DoctorCredentials, error)
}

type credentialsRepo struct {
	db *gorm.DB
}

func NewCredentialsRepository(db *gorm.DB) CredentialsRepository {
	return &credentialsRepo{db: db}
}

func (r *credentialsRepo) Create(ctx context.Context, c *pb.DoctorCredentials) (*pb.DoctorCredentials, error) {
	if err := r.db.WithContext(ctx).Table("doctor_credentials").Create(c).Error; err != nil {
		return nil, fmt.Errorf("insert doctor credentials: %w", err)
	}
	return c, nil
}

func (r *credentialsRepo) GetByDoctor(ctx context.Context, doctorUUID string) (*pb.DoctorCredentials, error) {
	var c pb.DoctorCredentials
	if err := r.db.WithContext(ctx).Table("doctor_credentials").Where("doctor_uuid = ?", doctorUUID).First(&c).Error; err != nil {
		return nil, fmt.Errorf("get doctor credentials: %w", err)
	}
	return &c, nil
}

func (r *credentialsRepo) Update(ctx context.Context, c *pb.DoctorCredentials) (*pb.DoctorCredentials, error) {
	res := r.db.WithContext(ctx).Table("doctor_credentials").Where("uuid = ?", c.GetUuid()).Updates(c)
	if res.Error != nil {
		return nil, fmt.Errorf("update doctor credentials: %w", res.Error)
	}
	return c, nil
}
