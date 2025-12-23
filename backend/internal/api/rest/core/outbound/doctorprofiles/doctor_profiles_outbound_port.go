package doctorprofiles

import (
	"context"

	pt "github.com/OPetricevic/physio-tracker/backend/golang/patients"
)

type Repository interface {
	GetByDoctor(ctx context.Context, doctorUUID string) (*pt.DoctorProfile, error)
	Upsert(ctx context.Context, profile *pt.DoctorProfile) (*pt.DoctorProfile, error)
}
