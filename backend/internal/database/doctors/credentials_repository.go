package doctors

import (
	"context"
	"errors"
	"fmt"

	pt "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	re "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/repoerrors"
	"github.com/jackc/pgconn"
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
	orm, err := c.ToORM(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating credentials: convert to ORM: %w", err)
	}
	if err := r.db.WithContext(ctx).Create(&orm).Error; err != nil {
		if isCredUniqueViolation(err) {
			return nil, fmt.Errorf("creating credentials: %w", re.ErrConflict)
		}
		return nil, fmt.Errorf("creating credentials: insert: %w", err)
	}
	pbObj, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating credentials: convert to PB: %w", err)
	}
	return &pbObj, nil
}

func (r *credentialsRepo) GetByDoctor(ctx context.Context, doctorUUID string) (*pt.DoctorCredentials, error) {
	var orm pt.DoctorCredentialsORM
	if err := r.db.WithContext(ctx).Where("doctor_uuid = ?", doctorUUID).First(&orm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("getting credentials: %w", re.ErrNotFound)
		}
		return nil, fmt.Errorf("getting credentials: %w", err)
	}
	pbObj, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting credentials: convert to PB: %w", err)
	}
	return &pbObj, nil
}

func (r *credentialsRepo) Update(ctx context.Context, c *pt.DoctorCredentials) (*pt.DoctorCredentials, error) {
	orm, err := c.ToORM(ctx)
	if err != nil {
		return nil, fmt.Errorf("updating credentials: convert to ORM: %w", err)
	}
	res := r.db.WithContext(ctx).Model(&orm).Where("uuid = ?", c.GetUuid()).Updates(&orm)
	if res.Error != nil {
		if isCredUniqueViolation(res.Error) {
			return nil, fmt.Errorf("updating credentials: %w", re.ErrConflict)
		}
		return nil, fmt.Errorf("updating credentials: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("updating credentials: %w", re.ErrNotFound)
	}
	pbObj, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("updating credentials: convert to PB: %w", err)
	}
	return &pbObj, nil
}

func isCredUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
