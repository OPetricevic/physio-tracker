package anamneses

import (
	"context"
	"errors"
	"fmt"
	"strings"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	out "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/outbound/anamneses"
	re "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/repoerrors"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, a *pb.Anamnesis) (*pb.Anamnesis, error) {
	orm, err := a.ToORM(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating anamnesis: convert to ORM: %w", err)
	}
	if err := r.db.WithContext(ctx).Create(&orm).Error; err != nil {
		return nil, fmt.Errorf("creating anamnesis: insert: %w", err)
	}
	pbObj, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating anamnesis: convert to PB: %w", err)
	}
	return &pbObj, nil
}

func (r *Repository) Update(ctx context.Context, a *pb.Anamnesis) (*pb.Anamnesis, error) {
	orm, err := a.ToORM(ctx)
	if err != nil {
		return nil, fmt.Errorf("updating anamnesis: convert to ORM: %w", err)
	}
	res := r.db.WithContext(ctx).Model(&orm).Where("uuid = ?", a.GetUuid()).Updates(&orm)
	if res.Error != nil {
		return nil, fmt.Errorf("updating anamnesis: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("updating anamnesis: %w", re.ErrNotFound)
	}
	pbObj, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("updating anamnesis: convert to PB: %w", err)
	}
	return &pbObj, nil
}

func (r *Repository) Get(ctx context.Context, uuid string) (*pb.Anamnesis, error) {
	var orm pb.AnamnesisORM
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&orm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("getting anamnesis: %w", re.ErrNotFound)
		}
		return nil, fmt.Errorf("getting anamnesis: %w", err)
	}
	pbObj, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting anamnesis: convert to PB: %w", err)
	}
	return &pbObj, nil
}

func (r *Repository) Delete(ctx context.Context, uuid string) error {
	res := r.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&pb.AnamnesisORM{})
	if res.Error != nil {
		return fmt.Errorf("delete anamnesis: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete anamnesis: %w", re.ErrNotFound)
	}
	return nil
}

func (r *Repository) List(ctx context.Context, patientUUID string, doctorUUID string, query string, limit, offset int) ([]*pb.Anamnesis, error) {
	var orms []pb.AnamnesisORM
	q := r.db.WithContext(ctx).Model(&pb.AnamnesisORM{}).
		Joins("JOIN patients ON patients.uuid = anamneses.patient_uuid").
		Where("patients.uuid = ?", patientUUID)
	if strings.TrimSpace(doctorUUID) != "" {
		q = q.Where("patients.doctor_uuid = ?", doctorUUID)
	}
	if strings.TrimSpace(query) != "" {
		like := "%" + strings.ToLower(strings.TrimSpace(query)) + "%"
		q = q.Where("LOWER(anamneses.diagnosis) LIKE ?", like)
	}
	if limit > 0 {
		q = q.Limit(limit).Offset(offset)
	}
	if err := q.Order("anamneses.created_at DESC").Find(&orms).Error; err != nil {
		return nil, fmt.Errorf("listing anamneses: %w", err)
	}
	res := make([]*pb.Anamnesis, 0, len(orms))
	for _, orm := range orms {
		pbObj, err := orm.ToPB(ctx)
		if err != nil {
			return nil, fmt.Errorf("listing anamneses: convert to PB: %w", err)
		}
		res = append(res, &pbObj)
	}
	return res, nil
}

func (r *Repository) ListByUUIDs(ctx context.Context, uuids []string) ([]*pb.Anamnesis, error) {
	if len(uuids) == 0 {
		return []*pb.Anamnesis{}, nil
	}
	var orms []pb.AnamnesisORM
	if err := r.db.WithContext(ctx).Where("uuid IN ?", uuids).Find(&orms).Error; err != nil {
		return nil, fmt.Errorf("listing anamneses by uuids: %w", err)
	}
	res := make([]*pb.Anamnesis, 0, len(orms))
	for _, orm := range orms {
		pbObj, err := orm.ToPB(ctx)
		if err != nil {
			return nil, fmt.Errorf("listing anamneses by uuids: convert to PB: %w", err)
		}
		res = append(res, &pbObj)
	}
	return res, nil
}

var _ out.Repository = (*Repository)(nil)
