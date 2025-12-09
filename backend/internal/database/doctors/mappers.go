package doctors

import (
	"time"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// doctorRecord is the DB-facing shape with time.Time for timestamps.
type doctorRecord struct {
	UUID      string     `gorm:"column:uuid;primaryKey"`
	Email     string     `gorm:"column:email"`
	Username  string     `gorm:"column:username"`
	FirstName string     `gorm:"column:first_name"`
	LastName  string     `gorm:"column:last_name"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

func protoToDoctorRecord(d *pb.Doctor) *doctorRecord {
	return &doctorRecord{
		UUID:      d.GetUuid(),
		Email:     d.GetEmail(),
		Username:  d.GetUsername(),
		FirstName: d.GetFirstName(),
		LastName:  d.GetLastName(),
		CreatedAt: tsToTime(d.CreatedAt),
		UpdatedAt: tsToTimePtr(d.UpdatedAt),
	}
}

func doctorRecordToProto(m doctorRecord) *pb.Doctor {
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

func doctorRecordsToProto(list []doctorRecord) []*pb.Doctor {
	res := make([]*pb.Doctor, 0, len(list))
	for _, m := range list {
		res = append(res, doctorRecordToProto(m))
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
