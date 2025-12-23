package anamneses

import (
	"context"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
)

// Repository defines outbound persistence for anamneses.
type Repository interface {
	Create(ctx context.Context, a *pb.Anamnesis) (*pb.Anamnesis, error)
	Update(ctx context.Context, a *pb.Anamnesis) (*pb.Anamnesis, error)
	Get(ctx context.Context, uuid string) (*pb.Anamnesis, error)
	Delete(ctx context.Context, uuid string) error
	List(ctx context.Context, patientUUID string, doctorUUID string, query string, limit, offset int) ([]*pb.Anamnesis, error)
}
