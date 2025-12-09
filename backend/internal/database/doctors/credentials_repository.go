package doctors

import (
	"context"
	"fmt"

	pt "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	"gorm.io/gorm"
)

// CredentialsRepository stores password hashes.
type CredentialsRepository interface {
	Create(ctx context.Context, c *pt.DoctorCredentials) (*pt.DoctorCredentials, error)
	GetByDoctor(ctx context.Context, doctorUUID string) (*pt.DoctorCredentials, error)
	Update(ctx context.Context, c *pt.DoctorCredentials) (*pt.DoctorCredentials, error)
}

type credentialsRepo struct {
	db *gorm.DB
}

func NewCredentialsRepository(db *gorm.DB) CredentialsRepository {
	return &credentialsRepo{db: db}
}

func (r *credentialsRepo) Create(ctx context.Context, c *pt.DoctorCredentials) (*pt.DoctorCredentials, error) {
	if err := r.db.WithContext(ctx).Table("doctor_credentials").Create(c).Error; err != nil {
		return nil, fmt.Errorf("insert doctor credentials: %w", err)
	}
	return c, nil
}

func (r *credentialsRepo) GetByDoctor(ctx context.Context, doctorUUID string) (*pt.DoctorCredentials, error) {
	var c pt.DoctorCredentials
	if err := r.db.WithContext(ctx).Table("doctor_credentials").Where("doctor_uuid = ?", doctorUUID).First(&c).Error; err != nil {
		return nil, fmt.Errorf("get doctor credentials: %w", err)
	}
	return &c, nil
}

func (r *credentialsRepo) Update(ctx context.Context, c *pt.DoctorCredentials) (*pt.DoctorCredentials, error) {
	res := r.db.WithContext(ctx).Table("doctor_credentials").Where("uuid = ?", c.GetUuid()).Updates(c)
	if res.Error != nil {
		return nil, fmt.Errorf("update doctor credentials: %w", res.Error)
	}
	return c, nil
}
