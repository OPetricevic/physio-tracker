package doctors

import (
	"context"
	"fmt"
	"strings"
	"time"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	out "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/outbound/doctors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type DoctorsRepository struct {
	db *gorm.DB
}

func NewDoctorsRepository(db *gorm.DB) *DoctorsRepository {
	return &DoctorsRepository{db: db}
}

func (r *DoctorsRepository) Create(ctx context.Context, d *pb.Doctor) (*pb.Doctor, error) {
	if err := r.db.WithContext(ctx).Table("doctors").Create(protoToModel(d)).Error; err != nil {
		return nil, fmt.Errorf("insert doctor: %w", err)
	}
	return d, nil
}

func (r *DoctorsRepository) Update(ctx context.Context, d *pb.Doctor) (*pb.Doctor, error) {
	res := r.db.WithContext(ctx).Table("doctors").Where("uuid = ?", d.GetUuid()).Updates(protoToModel(d))
	if res.Error != nil {
		return nil, fmt.Errorf("update doctor: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("update doctor: %w", gorm.ErrRecordNotFound)
	}
	return d, nil
}

func (r *DoctorsRepository) List(ctx context.Context, filter *pb.ListDoctorsRequest, limit, offset int) ([]*pb.Doctor, error) {
	var models []doctorModel
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
	return modelsToProto(models), nil
}

func (r *DoctorsRepository) Get(ctx context.Context, uuid string) (*pb.Doctor, error) {
	var m doctorModel
	if err := r.db.WithContext(ctx).Table("doctors").Where("uuid = ?", uuid).First(&m).Error; err != nil {
		return nil, fmt.Errorf("get doctor: %w", err)
	}
	return modelToProto(m), nil
}

func (r *DoctorsRepository) Delete(ctx context.Context, uuid string) error {
	res := r.db.WithContext(ctx).Table("doctors").Where("uuid = ?", uuid).Delete(&doctorModel{})
	if res.Error != nil {
		return fmt.Errorf("delete doctor: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete doctor: %w", gorm.ErrRecordNotFound)
	}
	return nil
}

type doctorModel struct {
	UUID      string     `gorm:"column:uuid;primaryKey"`
	Email     string     `gorm:"column:email"`
	Username  string     `gorm:"column:username"`
	FirstName string     `gorm:"column:first_name"`
	LastName  string     `gorm:"column:last_name"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

func protoToModel(d *pb.Doctor) *doctorModel {
	return &doctorModel{
		UUID:      d.GetUuid(),
		Email:     d.GetEmail(),
		Username:  d.GetUsername(),
		FirstName: d.GetFirstName(),
		LastName:  d.GetLastName(),
		CreatedAt: tsToTime(d.CreatedAt),
		UpdatedAt: tsToTimePtr(d.UpdatedAt),
	}
}

func modelToProto(m doctorModel) *pb.Doctor {
	return &pb.Doctor{
		Uuid:      m.UUID,
		Email:     m.Email,
		Username:  m.Username,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timeToTsPtr(m.UpdatedAt),
	}
}

func modelsToProto(list []doctorModel) []*pb.Doctor {
	res := make([]*pb.Doctor, 0, len(list))
	for _, m := range list {
		res = append(res, modelToProto(m))
	}
	return res
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

var _ out.Repository = (*DoctorsRepository)(nil)
