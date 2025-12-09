package doctors

import (
	"context"
	"fmt"
	"strings"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	"gorm.io/gorm"
)

type DoctorsRepository struct {
	db *gorm.DB
}

func NewDoctorsRepository(db *gorm.DB) *DoctorsRepository {
	return &DoctorsRepository{db: db}
}

func (r *DoctorsRepository) Create(ctx context.Context, d *pb.Doctor) (*pb.Doctor, error) {
	if err := r.db.WithContext(ctx).Table("doctors").Create(protoToDoctorRecord(d)).Error; err != nil {
		return nil, fmt.Errorf("insert doctor: %w", err)
	}
	return d, nil
}

func (r *DoctorsRepository) Update(ctx context.Context, d *pb.Doctor) (*pb.Doctor, error) {
	res := r.db.WithContext(ctx).Table("doctors").Where("uuid = ?", d.GetUuid()).Updates(protoToDoctorRecord(d))
	if res.Error != nil {
		return nil, fmt.Errorf("update doctor: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("update doctor: %w", gorm.ErrRecordNotFound)
	}
	return d, nil
}

func (r *DoctorsRepository) List(ctx context.Context, filter *pb.ListDoctorsRequest, limit, offset int) ([]*pb.Doctor, error) {
	var models []doctorRecord
	q := r.db.WithContext(ctx).Table("doctors")
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
	return doctorRecordsToProto(models), nil
}

func (r *DoctorsRepository) Get(ctx context.Context, uuid string) (*pb.Doctor, error) {
	var m doctorRecord
	if err := r.db.WithContext(ctx).Table("doctors").Where("uuid = ?", uuid).First(&m).Error; err != nil {
		return nil, fmt.Errorf("get doctor: %w", err)
	}
	return doctorRecordToProto(m), nil
}

func (r *DoctorsRepository) Delete(ctx context.Context, uuid string) error {
	res := r.db.WithContext(ctx).Table("doctors").Where("uuid = ?", uuid).Delete(&doctorRecord{})
	if res.Error != nil {
		return fmt.Errorf("delete doctor: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete doctor: %w", gorm.ErrRecordNotFound)
	}
	return nil
}
