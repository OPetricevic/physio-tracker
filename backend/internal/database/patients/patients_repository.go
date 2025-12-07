package patients

import (
	"context"
	"fmt"
	"strings"
	"time"

	patientspb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	out "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/outbound/patients"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type PatientsRepository struct {
	db *gorm.DB
}

func NewPatientsRepository(db *gorm.DB) *PatientsRepository {
	return &PatientsRepository{db: db}
}

func (r *PatientsRepository) Create(ctx context.Context, p *patientspb.Patient) (*patientspb.Patient, error) {
	if err := r.db.WithContext(ctx).Table("patients").Create(protoToModel(p)).Error; err != nil {
		return nil, fmt.Errorf("insert patient: %w", err)
	}
	return p, nil
}

func (r *PatientsRepository) Update(ctx context.Context, p *patientspb.Patient) (*patientspb.Patient, error) {
	res := r.db.WithContext(ctx).Table("patients").Where("uuid = ?", p.GetUuid()).Updates(protoToModel(p))
	if res.Error != nil {
		return nil, fmt.Errorf("update patient: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("update patient: %w", gorm.ErrRecordNotFound)
	}
	return p, nil
}

func (r *PatientsRepository) List(ctx context.Context, filter *patientspb.ListPatientsRequest, doctorUUID string, limit, offset int) ([]*patientspb.Patient, error) {
	var models []patientModel
	q := r.db.WithContext(ctx).Table("patients")
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
	if err := q.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list patients: %w", err)
	}
	return modelsToProto(models), nil
}

func (r *PatientsRepository) Get(ctx context.Context, uuid string) (*patientspb.Patient, error) {
	var m patientModel
	if err := r.db.WithContext(ctx).Table("patients").Where("uuid = ?", uuid).First(&m).Error; err != nil {
		return nil, fmt.Errorf("get patient: %w", err)
	}
	return modelToProto(m), nil
}

func (r *PatientsRepository) Delete(ctx context.Context, uuid string) error {
	res := r.db.WithContext(ctx).Table("patients").Where("uuid = ?", uuid).Delete(&patientModel{})
	if res.Error != nil {
		return fmt.Errorf("delete patient: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete patient: %w", gorm.ErrRecordNotFound)
	}
	return nil
}

// GORM model
type patientModel struct {
	UUID        string     `gorm:"column:uuid;primaryKey"`
	DoctorUUID  string     `gorm:"column:doctor_uuid"`
	FirstName   string     `gorm:"column:first_name"`
	LastName    string     `gorm:"column:last_name"`
	Phone       *string    `gorm:"column:phone"`
	Address     *string    `gorm:"column:address"`
	DateOfBirth *string    `gorm:"column:date_of_birth"`
	Sex         *string    `gorm:"column:sex"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
}

func protoToModel(p *patientspb.Patient) *patientModel {
	return &patientModel{
		UUID:        p.GetUuid(),
		DoctorUUID:  p.GetDoctorUuid(),
		FirstName:   p.GetFirstName(),
		LastName:    p.GetLastName(),
		Phone:       unwrap(p.Phone),
		Address:     unwrap(p.Address),
		DateOfBirth: unwrap(p.DateOfBirth),
		Sex:         unwrap(p.Sex),
		CreatedAt:   tsToTime(p.CreatedAt),
		UpdatedAt:   tsToTimePtr(p.UpdatedAt),
	}
}

func modelToProto(m patientModel) *patientspb.Patient {
	return &patientspb.Patient{
		Uuid:        m.UUID,
		DoctorUuid:  m.DoctorUUID,
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		Phone:       wrap(m.Phone),
		Address:     wrap(m.Address),
		DateOfBirth: wrap(m.DateOfBirth),
		Sex:         wrap(m.Sex),
		CreatedAt:   timestamppb.New(m.CreatedAt),
		UpdatedAt:   timeToTsPtr(m.UpdatedAt),
	}
}

func modelsToProto(list []patientModel) []*patientspb.Patient {
	res := make([]*patientspb.Patient, 0, len(list))
	for _, m := range list {
		res = append(res, modelToProto(m))
	}
	return res
}

func unwrap(sv *patientspb.StringValue) *string {
	if sv == nil {
		return nil
	}
	val := strings.TrimSpace(sv.GetValue())
	if val == "" {
		return nil
	}
	return &val
}

func wrap(s *string) *patientspb.StringValue {
	if s == nil {
		return nil
	}
	val := strings.TrimSpace(*s)
	if val == "" {
		return nil
	}
	return &patientspb.StringValue{Value: val}
}

func tsToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

func tsToTimePtr(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

func timeToTsPtr(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

var _ out.Repository = (*PatientsRepository)(nil)
