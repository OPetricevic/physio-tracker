package anamneses

import (
	"time"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	"github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// anamnesisRecord maps the table including include_visit_uuids text[] column.
type anamnesisRecord struct {
	Uuid              string         `gorm:"column:uuid;primaryKey"`
	PatientUuid       string         `gorm:"column:patient_uuid"`
	Anamnesis         string         `gorm:"column:anamnesis"`
	Diagnosis         string         `gorm:"column:diagnosis"`
	Therapy           string         `gorm:"column:therapy"`
	OtherInfo         string         `gorm:"column:other_info"`
	IncludeVisitUuids pq.StringArray `gorm:"column:include_visit_uuids;type:text[]"`
	CreatedAt         time.Time      `gorm:"column:created_at"`
	UpdatedAt         *time.Time     `gorm:"column:updated_at"`
}

func (anamnesisRecord) TableName() string { return "anamneses" }

func recordToPB(rec anamnesisRecord) *pb.Anamnesis {
	var upd *timestamppb.Timestamp
	if rec.UpdatedAt != nil {
		upd = timestamppb.New(*rec.UpdatedAt)
	}
	return &pb.Anamnesis{
		Uuid:              rec.Uuid,
		PatientUuid:       rec.PatientUuid,
		Anamnesis:         rec.Anamnesis,
		Diagnosis:         rec.Diagnosis,
		Therapy:           rec.Therapy,
		OtherInfo:         rec.OtherInfo,
		IncludeVisitUuids: []string(rec.IncludeVisitUuids),
		CreatedAt:         timestamppb.New(rec.CreatedAt),
		UpdatedAt:         upd,
	}
}

func pbToRecord(a *pb.Anamnesis) (anamnesisRecord, error) {
	rec := anamnesisRecord{
		Uuid:              a.GetUuid(),
		PatientUuid:       a.GetPatientUuid(),
		Anamnesis:         a.GetAnamnesis(),
		Diagnosis:         a.GetDiagnosis(),
		Therapy:           a.GetTherapy(),
		OtherInfo:         a.GetOtherInfo(),
		IncludeVisitUuids: pq.StringArray(a.GetIncludeVisitUuids()),
	}
	if a.GetCreatedAt() != nil {
		rec.CreatedAt = a.GetCreatedAt().AsTime()
	}
	if a.GetUpdatedAt() != nil {
		t := a.GetUpdatedAt().AsTime()
		rec.UpdatedAt = &t
	}
	return rec, nil
}
