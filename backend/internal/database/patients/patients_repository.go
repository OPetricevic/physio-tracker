package patients

import (
	"context"
	"errors"
	"fmt"
	"strings"

	pt "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	out "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/outbound/patients"
	re "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/repoerrors"
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
		return nil, fmt.Errorf("creating patient: convert to ORM: %w", err)
	}
	if err := r.db.WithContext(ctx).Create(&orm).Error; err != nil {
		return nil, fmt.Errorf("creating patient: insert: %w", err)
	}
	pbObj, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating patient: convert to PB: %w", err)
	}
	return &pbObj, nil
}

func (r *PatientsRepository) Update(ctx context.Context, p *pt.Patient) (*pt.Patient, error) {
	orm, err := p.ToORM(ctx)
	if err != nil {
		return nil, fmt.Errorf("updating patient: convert to ORM: %w", err)
	}
	res := r.db.WithContext(ctx).Model(&orm).Where("uuid = ?", p.GetUuid()).Updates(&orm)
	if res.Error != nil {
		return nil, fmt.Errorf("updating patient: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("updating patient: %w", re.ErrNotFound)
	}
	pbObj, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("updating patient: convert to PB: %w", err)
	}
	return &pbObj, nil
}

func (r *PatientsRepository) List(ctx context.Context, filter *pt.ListPatientsRequest, doctorUUID string, limit, offset int) ([]*pt.Patient, error) {
	var orms []pt.PatientORM
	q := r.db.WithContext(ctx).Model(&pt.PatientORM{})
	if strings.TrimSpace(doctorUUID) != "" {
		q = q.Where("doctor_uuid = ?", doctorUUID)
	}
	if terms := parseSearchTerms(filter.GetQuery()); len(terms) > 0 {
		for _, term := range terms {
			like := "%" + term + "%"
			q = q.Where("LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(phone) LIKE ?", like, like, like)
		}
	}
	if limit > 0 {
		q = q.Limit(limit).Offset(offset)
	}
	if err := q.Order("created_at DESC").Find(&orms).Error; err != nil {
		return nil, fmt.Errorf("listing patients: %w", err)
	}
	return patientORMsToProto(ctx, orms)
}

func (r *PatientsRepository) Get(ctx context.Context, uuid string) (*pt.Patient, error) {
	var orm pt.PatientORM
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&orm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("getting patient: %w", re.ErrNotFound)
		}
		return nil, fmt.Errorf("getting patient: %w", err)
	}
	pbObj, err := orm.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting patient: convert to PB: %w", err)
	}
	return &pbObj, nil
}

func (r *PatientsRepository) Delete(ctx context.Context, uuid string) error {
	res := r.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&pt.PatientORM{})
	if res.Error != nil {
		return fmt.Errorf("delete patient: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete patient: %w", re.ErrNotFound)
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

func parseSearchTerms(query string) []string {
	cleaned := strings.ToLower(strings.TrimSpace(query))
	if cleaned == "" {
		return nil
	}
	return strings.Fields(cleaned)
}

var _ out.Repository = (*PatientsRepository)(nil)
