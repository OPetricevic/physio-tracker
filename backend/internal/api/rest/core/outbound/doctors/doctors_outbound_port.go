package doctors

import (
	"context"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
)

type Repository interface {
	Create(ctx context.Context, d *pb.Doctor) (*pb.Doctor, error)
	Update(ctx context.Context, d *pb.Doctor) (*pb.Doctor, error)
	List(ctx context.Context, filter *pb.ListDoctorsRequest, limit, offset int) ([]*pb.Doctor, error)
	Get(ctx context.Context, uuid string) (*pb.Doctor, error)
	Delete(ctx context.Context, uuid string) error
}
