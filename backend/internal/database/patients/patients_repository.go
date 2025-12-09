package patients

import (
	"context"
	"fmt"
	"strings"

	pt "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	out "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/outbound/patients"
	"gorm.io/gorm"
)

type PatientsRepository struct {
	db *gorm.DB
}

func NewPatientsRepository(db *gorm.DB) *PatientsRepository {
	return &PatientsRepository{db: db}
}

func (r *PatientsRepository) Create(ctx context.Context, p *pt.Patient) (*pt.Patient, error) {
	orm, err := p.ToORM(ctx)
	if err != nil {
		return nil, fmt.Errorf("to ORM: %w", err)
	}
	if err := r.db.WithContext(ctx).Create(&orm).Error; err != nil {
		return nil, fmt.Errorf("insert patient: %w", err)
	}
	pbObj, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("to PB: %w", err)
	}
	return &pbObj, nil
}

func (r *PatientsRepository) Update(ctx context.Context, p *pt.Patient) (*pt.Patient, error) {
	orm, err := p.ToORM(ctx)
	if err != nil {
		return nil, fmt.Errorf("to ORM: %w", err)
	}
	res := r.db.WithContext(ctx).Model(&orm).Where("uuid = ?", p.GetUuid()).Updates(&orm)
	if res.Error != nil {
		return nil, fmt.Errorf("update patient: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("update patient: %w", gorm.ErrRecordNotFound)
	}
	pbObj, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("to PB: %w", err)
	}
	return &pbObj, nil
}

func (r *PatientsRepository) List(ctx context.Context, filter *pt.ListPatientsRequest, doctorUUID string, limit, offset int) ([]*pt.Patient, error) {
	var orms []pt.PatientORM
	q := r.db.WithContext(ctx).Model(&pt.PatientORM{})
	if strings.TrimSpace(doctorUUID) != "" {
		q = q.Where("doctor_uuid = ?", doctorUUID)
	}
	if strings.TrimSpace(filter.GetQuery()) != "" {
		like := "%" + strings.ToLower(strings.TrimSpace(filter.GetQuery())) + "%"
		q = q.Where("LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(phone) LIKE ?", like, like, like)
	}
	if limit > 0 {
		q = q.Limit(limit).Offset(offset)
	}
	if err := q.Order("created_at DESC").Find(&orms).Error; err != nil {
		return nil, fmt.Errorf("list patients: %w", err)
	}
	return patientORMsToProto(ctx, orms)
}

func (r *PatientsRepository) Get(ctx context.Context, uuid string) (*pt.Patient, error) {
	var orm pt.PatientORM
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&orm).Error; err != nil {
		return nil, fmt.Errorf("get patient: %w", err)
	}
	pbObj, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("to PB: %w", err)
	}
	return &pbObj, nil
}

func (r *PatientsRepository) Delete(ctx context.Context, uuid string) error {
	res := r.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&pt.PatientORM{})
	if res.Error != nil {
		return fmt.Errorf("delete patient: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete patient: %w", gorm.ErrRecordNotFound)
	}
	return nil
}

func patientORMsToProto(ctx context.Context, orms []pt.PatientORM) ([]*pt.Patient, error) {
	res := make([]*pt.Patient, 0, len(orms))
	for _, orm := range orms {
		pbObj, err := orm.ToPB(ctx)
		if err != nil {
			return nil, err
		}
		res = append(res, &pbObj)
	}
	return res, nil
}

var _ out.Repository = (*PatientsRepository)(nil)
