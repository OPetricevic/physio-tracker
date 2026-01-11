package outboundportpatients

import (
	"context"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
)

type Repository interface {
	Create(ctx context.Context, p *pb.Patient) (*pb.Patient, error)
	Update(ctx context.Context, p *pb.Patient) (*pb.Patient, error)
	List(ctx context.Context, filter *pb.ListPatientsRequest, doctorUUID string, limit, offset int) ([]*pb.Patient, error)
	Get(ctx context.Context, uuid string) (*pb.Patient, error)
	Delete(ctx context.Context, uuid string) error
}
