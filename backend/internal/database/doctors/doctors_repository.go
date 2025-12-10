package doctors

import (
	"context"
	"errors"
	"fmt"
	"strings"

	pt "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	re "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/repoerrors"
	"gorm.io/gorm"
)

type DoctorsRepository struct {
	db *gorm.DB
}

func NewDoctorsRepository(db *gorm.DB) *DoctorsRepository {
	return &DoctorsRepository{db: db}
}

func (r *DoctorsRepository) Create(ctx context.Context, d *pt.Doctor) (*pt.Doctor, error) {
	orm, err := d.ToORM(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating doctor: convert to ORM: %w", err)
	}
	if err := r.db.WithContext(ctx).Create(&orm).Error; err != nil {
		return nil, fmt.Errorf("creating doctor: insert: %w", err)
	}
	pbDoc, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating doctor: convert to PB: %w", err)
	}
	return &pbDoc, nil
}

func (r *DoctorsRepository) Update(ctx context.Context, d *pt.Doctor) (*pt.Doctor, error) {
	orm, err := d.ToORM(ctx)
	if err != nil {
		return nil, fmt.Errorf("updating doctor: convert to ORM: %w", err)
	}
	res := r.db.WithContext(ctx).Model(&orm).Where("uuid = ?", d.GetUuid()).Updates(&orm)
	if res.Error != nil {
		return nil, fmt.Errorf("updating doctor: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("updating doctor: %w", re.ErrNotFound)
	}
	pbDoc, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("updating doctor: convert to PB: %w", err)
	}
	return &pbDoc, nil
}

func (r *DoctorsRepository) Get(ctx context.Context, uuid string) (*pt.Doctor, error) {
	var orm pt.DoctorORM
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&orm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("getting doctor: %w", re.ErrNotFound)
		}
		return nil, fmt.Errorf("getting doctor: %w", err)
	}
	pbDoc, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting doctor: convert to PB: %w", err)
	}
	return &pbDoc, nil
}

func (r *DoctorsRepository) Delete(ctx context.Context, uuid string) error {
	res := r.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&pt.DoctorORM{})
	if res.Error != nil {
		return fmt.Errorf("deleting doctor: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("deleting doctor: %w", re.ErrNotFound)
	}
	return nil
}

func (r *DoctorsRepository) GetByIdentifier(ctx context.Context, identifier string) (*pt.Doctor, error) {
	term := strings.ToLower(strings.TrimSpace(identifier))
	if term == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var orm pt.DoctorORM
	if err := r.db.WithContext(ctx).
		Where("LOWER(email) = ? OR LOWER(username) = ?", term, term).
		First(&orm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("getting doctor by identifier: %w", re.ErrNotFound)
		}
		return nil, fmt.Errorf("getting doctor by identifier: %w", err)
	}
	pbDoc, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting doctor by identifier: convert to PB: %w", err)
	}
	return &pbDoc, nil
}

func (r *DoctorsRepository) List(ctx context.Context, query string, limit, offset int) ([]*pt.Doctor, error) {
	var models []pt.DoctorORM
	q := r.db.WithContext(ctx).Model(&pt.DoctorORM{})
	if strings.TrimSpace(query) != "" {
		like := "%" + strings.ToLower(strings.TrimSpace(query)) + "%"
		q = q.Where("LOWER(email) LIKE ? OR LOWER(username) LIKE ? OR LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ?", like, like, like, like)
	}
	if limit > 0 {
		q = q.Limit(limit).Offset(offset)
	}
	if err := q.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("listing doctors: %w", err)
	}
	res := make([]*pt.Doctor, 0, len(models))
	for _, orm := range models {
		pbDoc, err := orm.ToPB(ctx)
		if err != nil {
			return nil, fmt.Errorf("listing doctors: convert to PB: %w", err)
		}
		res = append(res, &pbDoc)
	}
	return res, nil
}
