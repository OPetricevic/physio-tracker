package doctors

import (
	"context"
	"fmt"
	"strings"

	pt "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		return nil, fmt.Errorf("to ORM: %w", err)
	}
	if err := r.db.WithContext(ctx).Create(&orm).Error; err != nil {
		return nil, fmt.Errorf("insert doctor: %w", err)
	}
	pbDoc, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("to PB: %w", err)
	}
	return &pbDoc, nil
}

func (r *DoctorsRepository) Update(ctx context.Context, d *pt.Doctor) (*pt.Doctor, error) {
	orm, err := d.ToORM(ctx)
	if err != nil {
		return nil, fmt.Errorf("to ORM: %w", err)
	}
	res := r.db.WithContext(ctx).Model(&orm).Where("uuid = ?", d.GetUuid()).Updates(&orm)
	if res.Error != nil {
		return nil, fmt.Errorf("update doctor: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("update doctor: %w", gorm.ErrRecordNotFound)
	}
	pbDoc, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("to PB: %w", err)
	}
	return &pbDoc, nil
}

func (r *DoctorsRepository) List(ctx context.Context, filter *pt.ListDoctorsRequest, limit, offset int) ([]*pt.Doctor, error) {
	var models []pt.DoctorORM
	q := r.db.WithContext(ctx).Model(&pt.DoctorORM{})
	if strings.TrimSpace(filter.GetQuery()) != "" {
		like := "%" + strings.ToLower(strings.TrimSpace(filter.GetQuery())) + "%"
		q = q.Where("LOWER(email) LIKE ? OR LOWER(username) LIKE ? OR LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ?", like, like, like, like)
	}
	if limit > 0 {
		q = q.Limit(limit).Offset(offset)
	}
	if err := q.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list doctors: %w", err)
	}
	return doctorORMsToProto(ctx, models)
}

func (r *DoctorsRepository) Get(ctx context.Context, uuid string) (*pt.Doctor, error) {
	var orm pt.DoctorORM
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&orm).Error; err != nil {
		return nil, fmt.Errorf("get doctor: %w", err)
	}
	pbDoc, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("to PB: %w", err)
	}
	return &pbDoc, nil
}

func (r *DoctorsRepository) Delete(ctx context.Context, uuid string) error {
	res := r.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&pt.DoctorORM{})
	if res.Error != nil {
		return fmt.Errorf("delete doctor: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete doctor: %w", gorm.ErrRecordNotFound)
	}
	return nil
}

func doctorORMsToProto(ctx context.Context, orms []pt.DoctorORM) ([]*pt.Doctor, error) {
	res := make([]*pt.Doctor, 0, len(orms))
	for _, orm := range orms {
		pbDoc, err := orm.ToPB(ctx)
		if err != nil {
			return nil, err
		}
		// Ensure timestamps are not nil for consistency
		if pbDoc.CreatedAt == nil {
			pbDoc.CreatedAt = timestamppb.New(orm.CreatedAt.AsTime())
		}
		res = append(res, &pbDoc)
	}
	return res, nil
}
