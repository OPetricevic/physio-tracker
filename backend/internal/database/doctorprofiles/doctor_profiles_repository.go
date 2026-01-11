package doctorprofiles

import (
	"context"
	"errors"
	"fmt"

	pt "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	out "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/outbound/doctorprofiles"
	re "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/repoerrors"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByDoctor(ctx context.Context, doctorUUID string) (*pt.DoctorProfile, error) {
	var orm pt.DoctorProfileORM
	if err := r.db.WithContext(ctx).Where("doctor_uuid = ?", doctorUUID).First(&orm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get doctor profile: %w", re.ErrNotFound)
		}
		return nil, fmt.Errorf("get doctor profile: %w", err)
	}
	pbObj, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("get doctor profile: convert to PB: %w", err)
	}
	return &pbObj, nil
}

func (r *Repository) Upsert(ctx context.Context, profile *pt.DoctorProfile) (*pt.DoctorProfile, error) {
	orm, err := profile.ToORM(ctx)
	if err != nil {
		return nil, fmt.Errorf("upsert doctor profile: convert to ORM: %w", err)
	}
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing pt.DoctorProfileORM
		if err := tx.Where("doctor_uuid = ?", profile.GetDoctorUuid()).First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return tx.Create(&orm).Error
			}
			return err
		}
		orm.Uuid = existing.Uuid
		// Preserve created_at while allowing empty strings to overwrite optional fields.
		orm.CreatedAt = existing.CreatedAt
		return tx.Model(&existing).
			Select(
				"practice_name",
				"department",
				"role_title",
				"address",
				"phone",
				"email",
				"website",
				"logo_path",
				"oib_owner",
				"updated_at",
			).
			Where("doctor_uuid = ?", profile.GetDoctorUuid()).
			Updates(&orm).Error
	})
	if err != nil {
		return nil, fmt.Errorf("upsert doctor profile: %w", err)
	}
	pbObj, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("upsert doctor profile: convert to PB: %w", err)
	}
	return &pbObj, nil
}

var _ out.Repository = (*Repository)(nil)
