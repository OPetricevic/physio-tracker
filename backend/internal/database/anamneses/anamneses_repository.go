package anamneses

import (
	"context"
	"errors"
	"fmt"
	"strings"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	out "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/outbound/anamneses"
	re "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/repoerrors"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, a *pb.Anamnesis) (*pb.Anamnesis, error) {
	rec, err := pbToRecord(a)
	if err != nil {
		return nil, fmt.Errorf("creating anamnesis: convert: %w", err)
	}
	if err := r.db.WithContext(ctx).Create(&rec).Error; err != nil {
		return nil, fmt.Errorf("creating anamnesis: insert: %w", err)
	}
	return recordToPB(rec), nil
}

func (r *Repository) Update(ctx context.Context, a *pb.Anamnesis) (*pb.Anamnesis, error) {
	rec, err := pbToRecord(a)
	if err != nil {
		return nil, fmt.Errorf("updating anamnesis: convert: %w", err)
	}
	res := r.db.WithContext(ctx).
		Model(&anamnesisRecord{}).
		Where("uuid = ?", a.GetUuid()).
		Updates(map[string]interface{}{
			"patient_uuid":        rec.PatientUuid,
			"anamnesis":           rec.Anamnesis,
			"status":              rec.Status,
			"diagnosis":           rec.Diagnosis,
			"therapy":             rec.Therapy,
			"other_info":          rec.OtherInfo,
			"include_visit_uuids": pq.StringArray(rec.IncludeVisitUuids),
			"updated_at":          rec.UpdatedAt,
		})
	if res.Error != nil {
		return nil, fmt.Errorf("updating anamnesis: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("updating anamnesis: %w", re.ErrNotFound)
	}
	return recordToPB(rec), nil
}

func (r *Repository) Get(ctx context.Context, uuid string) (*pb.Anamnesis, error) {
	var rec anamnesisRecord
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&rec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("getting anamnesis: %w", re.ErrNotFound)
		}
		return nil, fmt.Errorf("getting anamnesis: %w", err)
	}
	return recordToPB(rec), nil
}

func (r *Repository) Delete(ctx context.Context, uuid string) error {
	res := r.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&anamnesisRecord{})
	if res.Error != nil {
		return fmt.Errorf("delete anamnesis: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete anamnesis: %w", re.ErrNotFound)
	}
	return nil
}

func (r *Repository) List(ctx context.Context, patientUUID string, doctorUUID string, query string, limit, offset int) ([]*pb.Anamnesis, error) {
	var recs []anamnesisRecord
	q := r.db.WithContext(ctx).Model(&anamnesisRecord{}).
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
	if err := q.Order("anamneses.created_at DESC").Find(&recs).Error; err != nil {
		return nil, fmt.Errorf("listing anamneses: %w", err)
	}
	res := make([]*pb.Anamnesis, 0, len(recs))
	for _, rec := range recs {
		res = append(res, recordToPB(rec))
	}
	return res, nil
}

func (r *Repository) ListByUUIDs(ctx context.Context, uuids []string) ([]*pb.Anamnesis, error) {
	if len(uuids) == 0 {
		return []*pb.Anamnesis{}, nil
	}
	var recs []anamnesisRecord
	if err := r.db.WithContext(ctx).Where("uuid IN ?", uuids).Find(&recs).Error; err != nil {
		return nil, fmt.Errorf("listing anamneses by uuids: %w", err)
	}
	res := make([]*pb.Anamnesis, 0, len(recs))
	for _, rec := range recs {
		res = append(res, recordToPB(rec))
	}
	return res, nil
}

var _ out.Repository = (*Repository)(nil)
